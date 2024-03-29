package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Statistics struct {
	os               string
	appDir           string
	tmpDir           string
	testDir          string
	originalDomains  int
	originalUrls     int
	totalUrlsVisited int
	date             string
	year             string
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	// Bad Globals till I decide on flags, or modules
	TokenFilePathGlobal string
)

// Restructure
func checkError(err error) error {
	if err != nil {
		fmt.Errorf("%s", err)
		log.Fatal(err)
		panic(err)
	}
	return err
}

func keyToDomainString(keyUrl string) (string, error) {
	parsedURL, err := url.Parse(keyUrl)
	if err != nil {
		fmt.Errorf("invalid url: %s\n", keyUrl)
		checkError(err)
	}

	hostname := parsedURL.Hostname()

	switch hostname {
	case "arstechnica.com":
	case "portswigger.net":
	case "thehackernews.com":
	case "www.sans.org":
	default:

	}
	return hostname, nil
}

func createFile(filepath string) error {
	filePtr, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Creation Error:", err)
		checkError(err)
	}
	defer filePtr.Close()
	return nil
}

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if os.IsNotExist(err) {
		log.Fatal("File path does not exist")
		return false, err
	}
	return true, nil
}

// InfoLogger.Printf("Something noteworthy happened\n")
// WarningLogger.Printf("There is something you should know about\n")
// ErrorLogger.Printf("Something went wrong\n")
func initaliseLogging() error {
	now := time.Now().UTC()
	dateFormatted := now.Format("2006-01-01")
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(dateFormatted)
	nameBuilder.WriteString(".log")
	file, err := os.OpenFile(nameBuilder.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0661)
	checkError(err)

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// Map = domain + id : url
func marshalURLsToMap() (map[string]string, map[string]int, error) {
	file, err := os.Open("urls.txt")
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urlsMapped := make(map[string]string)
	urlStr := ""
	domainCounter := make(map[string]int)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		urlStr = scanner.Text()
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("Invalid URL: %s\n", urlStr)
			continue
		}

		hostname := parsedURL.Hostname()
		counter, ok := domainCounter[hostname]
		if !ok {
			counter = 0
		}
		key := hostname + ":" + strconv.Itoa(counter)
		urlsMapped[key] = urlStr
		domainCounter[hostname] = counter + 1
	}

	return urlsMapped, domainCounter, nil
}

func mkDirAndCD(date string) error {
	if err := os.Mkdir(date, os.ModePerm); err != nil {
		log.Fatal(err)
		return err
	}
	err := os.Chdir(date)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func checkPrevRuntimes(appDir, date string) error {
	dirListing, err := os.ReadDir(appDir)
	if err != nil {
		return err
	}

	for _, dir := range dirListing {
		if dir.Name() == date {
			log.Fatal(errors.New("Directory already exists"))
			return errors.New("Directory already exists")
		}
	}

	return nil
}

func mkAppDirTree(appDir string, dirTree []string) error {
	var PathAndName string
	for _, dirName := range dirTree {
		PathAndName = filepath.Join(appDir, dirName)
		if err := os.Mkdir(PathAndName, os.ModePerm); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func trimFilePath(path string) (result string, err error) {
	os := runtime.GOOS
	switch os {
	case "windows":
		pathSlice := strings.SplitAfterN(path, "\\", -1)
		result = pathSlice[len(pathSlice)-1]
	case "linux":
		pathSlice := strings.SplitAfterN(path, "/", -1)
		result = pathSlice[len(pathSlice)-1]
	default:
		err := fmt.Errorf("unsupported os for filepath trimming of delimited %s", os)
		checkError(err)
		return "", err
	}
	return result, err
}

func curlNewBasePages(urlArr []string) (map[string]string, error) {
	var args string = "-K curlrc -L "
	result := make(map[string]string)
	for _, url := range urlArr {
		runCurl := exec.Command("curl", args, url)
		outputBytes, err := runCurl.Output()
		checkError(err)
		result[url] = string(outputBytes[:])
	}
	return result, nil
}

// Does this really need to be -O
func curlNewArticles(urlArr []string) error {
	var args string = "-K curlrc -L -O"
	urlsStr := strings.Join(urlArr, " ")
	runCurl := exec.Command("curl", args, urlsStr)
	err := runCurl.Run()
	checkError(err)
	return nil
}

func gzlopBuffer(buffer *bytes.Buffer, patterns []byte) (map[int]string, error) {
	patCount := int(0)
	artifacts := make(map[int]string)
	builder := strings.Builder{}
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		for i := 0; i <= len(patterns)-1; i++ {
			if bytes.Contains(scanner.Bytes(), []byte{patterns[i]}) {
				patCount++
				builder.WriteString(string(patterns[i]))
				artifacts[patCount] = builder.String()
				builder.Reset()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return artifacts, nil
}

func getTokensFileContentsAsBytes() ([]byte, error) {
	var tokenFileAsBytes []byte
	exists, err := checkFileExists(TokenFilePathGlobal)
	checkError(err)
	if !exists {
		// tokens file path does not exist
	} else {
		file, err := os.Open("file.txt")
		checkError(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			words := strings.Fields(line)
			for _, word := range words {
				tokenFileAsBytes = append(tokenFileAsBytes, []byte(word)...)
			}
		}
	}
	return tokenFileAsBytes, nil
}

// Add to nestedStrStrMap is just a better styled answer from: https://stackoverflow.com/questions/64918219/how-to-assign-to-a-nested-map
// Iterating over nested String maps - 1st is the [DOMAIN] so prints keys and values as key:value; 2nd prints just the values of [Domain][url]VALUE
// for _, mapKey := range testmap { fmt.Println(url)	}
// for _, key := range testmap["portswigger"] {fmt.Println(url)	}
func addValueToNestedStrStrMap(parentMap map[string]map[string]string, parentKey, childKey string, nestedValue string) {
    childMap := parentMap[parentKey]
	if childMap == nil {
		childMap = make(map[string]string)
		parentMap[parentKey] = childMap
	}
	childMap[childKey] = nestedValue
}


// https://www.tutorialspoint.com/golang-program-to-convert-file-to-byte-array
func loadTokensIntoMem(path string) (error, int ,[]bytes) {
	tokensFile, err := os.Open(path)
	if err != nil {
		checkErr(err)
	}
	defer tokensFile.Close()

	stat, err := tokensFile.Stat()
	if err != nil {
		checkErr(err)
	}
	
	byteSlice := make([]byte, stat.Size())
	_, err = bufio.NewReader(tokensFile).Read(byteSlice)
	if err != nil && err != io.EOF {
		checkError(err)
		return err, 0, nil
	}
	bsSize = len(byteSlice)
	return nil, bsSize, byteSlice
}


func main() {
	appDir, err := os.Getwd()
	if err != nil {
		checkError(err)
	}
	stat := Statistics{}
	stat.os = runtime.GOOS
	stat.tmpDir = os.TempDir()
	stat.appDir = appDir
	now := time.Now().UTC()
	stat.date = now.Format("2006-01-01")
	stat.year = strconv.Itoa(now.Year())
	err = checkPrevRuntimes(stat.appDir, stat.date)
	if err != nil {
		checkError(err)
	}
	dirTree := []string{"test", "logs", "newletters", stat.year}
	err = mkAppDirTree(stat.appDir, dirTree)
	if err != nil {
		checkError(err)
	}
	stat.testDir = filepath.Join(appDir, "test")
	err = initaliseLogging()
	if err != nil {
		checkError(err)
	}
	InfoLogger.Printf("Logging initialised")

	err = mkDirAndCD(stat.date)
	if err != nil {
		checkError(err)
	}
	saveDirectory := stat.date

	urlsToVisit, baseDNSurlTotals, err := marshalURLsToMap()
	if err != nil {
		checkError(err)
	}
	allBaseUrlsSeq := make([]string, 0, len(urlsToVisit))
	for _, value := range urlsToVisit {
		allBaseUrlsSeq = append(strings.Split(value, ""))
	}

	totalUrls := 0
	stat.originalDomains = len(baseDNSurlTotals) - 1

	for _, val := range baseDNSurlTotals {
		totalUrls = +val
	}

	stat.originalUrls = totalUrls
	stat.totalUrlsVisited += totalUrls

	// stdout -> 4 base pages
	// create maps for each base pages
	basePagesStdoutMap, err := curlNewBasePages(allBaseUrlsSeq)
	if err != nil {
		checkError(err)
	}

	// Get all links and titles from Base pages
	artefactsFromBasePages, err := getAllTitlesAndLinks(basePagesStdoutMap)
	if err != nil {
		checkError(err)
	}

	// Load tokens into memory
	tokensBuffer, err := loadTokensIntoMem(tokensFilePath)
	
	failedLinksAndTitleByDomainMap := make(map[string]map[string]string)

	assignTokenBufferOffsets()
	go func() {
		titleCheckResult, err = parseTitles(tokenOffset)
		if !titleCheckResult {
			addValueToNestedStrStrMap(failedLinksAndTitleByDomainMap, url, url, title)
		}
	}()

	// TODO - awaiting consideration:
	// Compare against historic file of links and titles
	// - File for each? 
	finalTitlesAndLinks, err := compareArtefactsToHistoricData()


	// go routine to fork out and get the page from each link - fork by some -T threads ( threads requested % functions ) for equal links per thread

	// Manage being able reading tokens in memory based on circular buffer 	   ([ x -> y -> z -] )
	// So if there are three threads x,y,z they read from an offsest circularly [<- zstradle--|]
	
	// make sure everything converges in a go way

	// Manage being able reading tokens in memory based on circular buffer 	   ([ x -> y -> z -] )
	// So if there are three threads x,y,z they read from an offsest circularly [<- zstradle--|]
	// Curl pages to memory
	// search for token found limit (bear in mind the amount of tokens is not large so worry about closure is not a problem)
	// assessResult based on a config file on WHAT constitues
	// marshall results from enumerated pages
	
	// DO I need to actually worry its a read not a write?
	// DO I actually need mutexs for maps for writes?

	assignTokenBufferOffsets() // array -> tokenBufferthreadId
	go func() {
		page, err = curlNewArticle(url)
		if err != nil {
		checkError(err)
		}
		matchedTokens, goodPage, err := parsePage(page, tokenOffset)
		if err != nil {
			checkError(err)
		}
		err := marshallParserResults(goodPage, matchedTokens, url)
		if err != nil {
			checkError(err)
		}
	}()
	// make sure everything converges in a go way

	
	// TODO double check this jibberish:

	// IS J'SON actually the way and is csv just bad or is google data broking fucking with me we design decision dilemia doubling 
	// ---- only need to store and compare urls
	// if in the file remove from map
	// Storage 2 files one .csv per run and collective with Page rating, time, url, matched tokens, And just previous-urls-found-only.list
	// compare maps for domain against previous enumerated list file with gzlop
	// Print Alert - similiar to each row of .csv of urls
	backupDataStorage()
	updateDataStorage()

	// Output cli, file and (backup and then) organise historic data
	err = selectOutput(outputArgs)
	if err != nil {
		checkError(err)
	}


}


// REASONS keys for failed-map so that it makes sense

// TODO site specifics
//
// findTokensOnPage() != findLinksAndTitlesFromBasePages
//
// gzlopGetHtmlTagAndLink()
// search -> reborderise -> then back
// search -> reborderise -> then forward


func marshallParserResults(goodPage bool, matchedTokens string, url string) error {
	if !goodPage {
			// remove url from queue,
			addValueToNestedStrStrMap(failedLinksAndTitleByDomainMap, "Parsed-Page-Results-Negative" , url, titles)
		} else {
			
		}
	return nil
}


// -
// -
// Low Sleep Idiocy, but Past me is still got the right ideas - just needs refactoring to main function specification scaffolding
// - 
// -
func getAllTitlesAndLinks(basePagesStdoutMap map[string]string) (map[string]map[int]string, error) {

	arstechnicaTokensMap, portswiggerTokensMap, thehackernewsTokensMap, sansTokensMap := make(map[int]string), make(map[int]string), make(map[int]string), make(map[int]string)
	for key, value := range basePagesStdoutMap {
		domain, err := keyToDomainString(key)
		checkError(err)

		webPageBuffer := bytes.NewBuffer([]byte(value))
		tokens, err := getTokensFileContentsAsBytes(TokenFilePathGlobal)
		checkError(err)
		switch domain {
		// sans
		// arstechnica
		// thehackernews
		//
		case "arstechnica.com":
			arstechnicaTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
			checkError(err)
		case "portswigger.net": // portswigger -> links are just title strings.Join(titleNoAtags, "-")
			portswiggerTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
			checkError(err)
		case "thehackernews.com":
			thehackernewsTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
			checkError(err)
		case "www.sans.org":
			sansTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
			checkError(err)
		case "":
			err := fmt.Errorf("strange race condition occur with domain variable being an empty string")
			checkError(err)
		default:

		}
		domain = ""
	}
	resultMap := make(map[string]map[int]string)
	resultMap["arstechnica.com"], resultMap["portswigger.net"], resultMap["thehackernews.com"], resultMap["www.sans.org"] = arstechnicaTokensMap, portswiggerTokensMap, thehackernewsTokensMap, sansTokensMap

	return resultMap, nil
}

func compareTitlesAndLinksToHistoricData(historicUrlsFile, tokensFile string, urlsFound map[string]string) error {
	var urlsAsBytes []byte

	file, err := os.ReadFile(historicUrlsFile)
	checkError(err)

	for _, url := range urlsFound {
		urlAsBytes := []byte(url)
		append(urlsAsBytes, urlAsBytes[:])
	}

	searchTokens, err := os.ReadFile(tokensFile)
	checkError(err)
	// TODO

	// map memory - for the same ~~paragraph~~ search for dates, url and tokens
	// soup go gets all the fields that have urls like gospider (CHECK HOW THAT WORK and do it locally)
	// gzlop buffer can then be adapter to search the buffer from address to offset for EVEN MORE SPEED

	// Naive Search for a token

	allTheArtefacts := make(map[int]map[int]string)
	for _, token := range searchTokens {
		artifacts, err := gzlopBuffer(file, token)
		checkError(err)
		allTheArtefacts[token] = artifacts
	}
	//
	// BRAIN NEED THUNK HERE
	//
	return nil


// -
// -
// OUTPUT TODO
// -
// -
func selectOutput(args []string) error {
	argsSize := len(strings.Join(args, ""))
	var argsId int
	if argsSize != 0 {
		for _, arg := range args {
			switch arg {
			case "verbose":
				argsId += 1
			case "cli":
				argsId += 2
			case "markdown":
				argsId += 5

			default:
				err := fmt.Errorf("invalid output arguments provide: %v ; from slice of size: %v with the contents: %v", arg, argsSize, args)
				checkError(err)
				return err
			}
		}
	} else {
		argsId = argsSize
	}
	switch argsId {
	case 1: // verbose
		verboseOutput()
	case 2: // cli only
		cliOnlyOutput()
	case 3: // verbose cli only
		verboseCliOutput()
	case 5: // markdown only
		markdownOnlyOutput()
	case 6: // verbose markdown
		verboseMarkdownOutput()
	case 0:
		defaultOutput()
	default:
		err := fmt.Errorf("invalid arg idenfier counted %v", argsId)
		checkError(err)
		return err
	}
	return nil

}

// verboseOutput - markdown + verbose markdown report, cli is verbose
func verboseOutput() {

}

func cliOnlyOutput() {

}

func markdownOnlyOutput() {

}

func verboseCliOutput() {

}

func verboseMarkdownOutput() {

}

// defaultOutput - cli some statistics and markdown report

func defaultOutput() {

}
