package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Application struct {
	appDir               string
	tmpDir               string
	testDir              string
	previousRuntime      string
	historicDataFilePath string
	statistics           Statistics
	noGzhodanConfig      bool
	multiDaily           bool
	outputType           string
	gzhodanConfig        string
	optionalConfigs      map[string]string
	tokensFile           string
}

type Statistics struct {
	operatingSystem  string
	originalDomains  int
	originalUrls     int
	totalUrlsVisited int
	date             time.Time
	year             string
	appStartTime     time.Time
}

type CircularBuffer struct {
	buffer       []byte
	size         int
	readPointers []int
}

func newCircularBuffer(data []byte, concurrencyOffset, workers int) *CircularBuffer {
	return &CircularBuffer{
		buffer:       data,
		size:         len(data) - 1,
		readPointers: make([]int, concurrencyOffset*workers),
	}
}

func (c *CircularBuffer) readCircularBufferFromOffset(worker int) (data byte) {
	data = c.buffer[c.readPointers[worker]]
	c.readPointers[worker] = (c.readPointers[worker] + 1) % c.size
	return data
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// https://gobyexample.com/reading-files
func (a *Application) loadTokensIntoMem() ([]byte, int, error) {
	exists, err := checkFileExists(a.tokensFile)
	if err != nil || !exists {
		checkError(err, 0, 0)
	}
	tokensFileAsBytes, err := os.ReadFile(a.tokensFile)
	if err != nil {
		checkError(err, 0, 0)
	}

	bsSize := len(tokensFileAsBytes) - 1
	return tokensFileAsBytes, bsSize, nil
}

func (a *Application) handleArgs(args []string, argsLength int) error {
	for i := 0; i <= argsLength-1; i++ {
		switch args[i] {
		case "-o":
			// Check Directory exist function required
			a.appDir = args[i+1]
		case "-m":
			a.multiDaily = true
		case "-H":
			historicDataExists, err := checkFileExists(args[i+1])
			if err != nil || !historicDataExists {
				checkError(err, 0, 0)
				WarningLogger.Printf("No historic data file provided to compare new url with previously enumerated data, this may take a lot longer!")
			} else {
				a.historicDataFilePath = args[i+1]
			}
		case "-O":
			configListSlice := strings.SplitAfterN(args[i+1], ",", -1)
			tmpmap := make(map[string]string)
			for _, config := range configListSlice {
				var key string
				if strings.Contains(config, "/") {
					keySlice := strings.SplitAfterN(config, "/", -1)
					key = keySlice[len(keySlice)-1]
				}
				if strings.Contains(config, "\\") {
					keySlice := strings.SplitAfterN(config, "\\", -1)
					key = keySlice[len(keySlice)-1]
				}
				keySlice := strings.Split(key, ".")
				key = keySlice[0]
				exists, err := checkFileExists(config)
				if err != nil || !exists {
					checkError(err, 0, 0)
				}
				tmpmap[key] = config
			}
		case "-t":
			exists, err := checkFileExists(args[i+1])
			if err != nil || !exists {
				checkError(err, 0, 0)
			}
			a.tokensFile = args[i+1]
		case "-G":
			exists, err := checkFileExists(args[i+1])
			if err != nil || !exists {
				checkError(err, 0, 0)
			}
			a.gzhodanConfig = args[i+1]
		case "-g":
			a.noGzhodanConfig = true
		case "-C":
			if a.outputType != "" {
				a.outputType = a.outputType + "-C"
			} else {
				a.outputType = "C"
			}
		case "-M":
			if a.outputType != "" {
				a.outputType = a.outputType + "-M"
			} else {
				a.outputType = "M"
			}
		case "-V":
			if a.outputType != "" {
				a.outputType = a.outputType + "-V"
			} else {
				a.outputType = "V"
			}
		default:
			err := fmt.Errorf("invalid arguments provided: %v", args)
			checkError(err, 0, 0)
		}
	}

	return nil
}

func (a *Application) selectOutput() error {
	argsSize := len(a.outputType)
	var argsId int = 0
	if argsSize != 1 {
		switch a.outputType {
		case "C-V":
			argsId = 3
		case "V-C":
			argsId = 3
		case "M-V":
			argsId = 7
		case "V-M":
			argsId = 7
		default:
			argsId = 0
			err := fmt.Errorf("invalid output arguments provide: %v ; from slice of size: %v", a.outputType, argsSize)
			checkError(err, 0, 0)
			return err
		}
	} else {
		switch a.outputType {
		case "V":
			argsId = 1
		case "C":
			argsId = 2
		case "M":
			argsId = 5
		default:
			err := fmt.Errorf("invalid output arguments provide: %v of size %v", a.outputType, argsSize)
			checkError(err, 0, 0)
			return err
		}
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
		checkError(err, 0, 0)
		return err
	}
	return nil
}

// Restructure err, 0, 0
func checkError(err error, errorLevel, errorCode int) {
	switch errorLevel {
	case 0:
		InfoLogger.Printf("test passed - error code:%v", errorCode)
		return
	case 1:
		WarningLogger.Printf("error code %v:%s", errorCode, err)
		return
	case 2:
		ErrorLogger.Printf("error code %v:%s", errorCode, err)
		log.Fatal(err)
	default:
		err := fmt.Errorf("incorrect errorlevel integer: %v by errorcode: %v", errorLevel, errorCode)
		log.Fatal(err)
	}
}

func (a *Application) compareUrlsHistorically(urlsFound map[string]string) (map[string]map[string]string, map[string]map[string]string, error) {
	var allUrlsAsBytes []byte
	var domain, url string = "", ""
	exists, err := checkFileExists(a.historicDataFilePath)
	if err != nil || !exists {
		checkError(err, 0, 0)
	}
	historicDataAsBytes, err := os.ReadFile(a.historicDataFilePath)
	if err != nil {
		checkError(err, 0, 0)
	}

	allTitles := make(map[int]string)
	i := 0
	for urlKey, titleValue := range urlsFound {
		allUrlsAsBytes = append([]byte(urlKey))
		allTitles[i] = titleValue
		i++
	}

	goodUrls := make(map[string]map[string]string)
	badUrls := make(map[string]map[string]string)

	dateAsBytesSize := len(historicDataAsBytes)
	for _, urlAsBytes := range allUrlsAsBytes {
		for i := 0; i <= dateAsBytesSize-1; i++ {
			if historicDataAsBytes[i] == urlAsBytes {
				url = string(urlAsBytes)
				domain, err = urlKeyToDomainString(url)
				if err != nil {
					checkError(err, 0, 0)
				}
				addValueToNestedStrStrMap(goodUrls, domain, url, allTitles[i])
				domain, url = "", ""
			} else {
				url = string(urlAsBytes)
				domain, err = urlKeyToDomainString(url)
				if err != nil {
					checkError(err, 0, 0)
				}
				addValueToNestedStrStrMap(badUrls, domain, url, allTitles[i])
			}
		}
	}

	return goodUrls, badUrls, nil
}

func urlKeyToDomainString(keyUrl string) (string, error) {
	parsedURL, err := url.Parse(keyUrl)
	if err != nil {
		fmt.Errorf("invalid url: %s\n", keyUrl)
		checkError(err, 0, 0)
	}
	return parsedURL.Hostname(), nil
}

func createFile(filepath string) error {
	filePtr, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Creation Error:", err)
		checkError(err, 0, 0)
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
	checkError(err, 0, 0)

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// Map = domain + id : url
func marshalURLsToMap() (map[string]string, map[string]int, error) {
	file, err := os.Open("urls.txt")
	checkError(err, 0, 0)
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

// Later functionality when there is alot data at some point we need condensing or checking
// to or with a historicData file, these kind of programs need data regression to best dataset (size, quality, parsability,etc)
// Do not remove
func (a *Application) checkPrevRuntimes() error {
	dirListing, err := os.ReadDir(a.appDir)
	if err != nil {
		checkError(err, 0, 0)
		return err
	}

	if len(dirListing) < 1 {
		InfoLogger.Printf("No previous data found\n")
		a.previousRuntime = ""
		return nil
	}

	currDateStr := time.Time.String(a.statistics.date)
	compare, err := time.Parse(time.DateOnly, "2006-03-03")
	for _, dir := range dirListing {
		tmp, err := time.Parse(time.DateOnly, dir.Name())
		if err != nil {
			checkError(err, 0, 0)
			continue
		}
		if dir.Name() == currDateStr {
			if !a.multiDaily {
				err := fmt.Errorf("directory already exists %s", dir.Name())
				checkError(err, 0, 0)
				continue
			} else {
				a.previousRuntime = currDateStr
				return nil
			}
		}
		if tmp.After(compare) {
			compare = tmp
		}

		if compare.After(a.statistics.date) {
			err := fmt.Errorf("current date %v, is before a directory already made of date: %v", a.statistics.date, compare)
			checkError(err, 0, 0)
			return err
		}
	}
	a.previousRuntime = time.Time.String(compare)
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
		checkError(err, 0, 0)
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
		checkError(err, 0, 0)
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
	checkError(err, 0, 0)
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

func parseAllBasePagesForLinksAndTitles(basePagesStdoutMap map[string]string) (map[string]string, error) {
	hrefPathAndTitlesRegexp := regexp.MustCompile(`<a href="\/.{1,}<\/a>`)
	basePageAllHrefs := make([]string, 0)
	allUrlsAndTitles := make(map[string]string)
	for siteUrl, page := range basePagesStdoutMap {
		pageLines := strings.SplitAfterN(page, "\n", -1)
		for _, line := range pageLines {
			match, err := regexp.MatchString(hrefPathAndTitlesRegexp.String(), line)
			if err != nil {
				checkError(err, 0, 0)
			}
			if match {
				basePageAllHrefs = append([]string{line})
			}
		}
		for _, hrefPathAndTitle := range basePageAllHrefs {
			doubleQuoteSplitHref := strings.SplitAfterN(hrefPathAndTitle, "\"", -1)
			titleTmp := strings.Replace(doubleQuoteSplitHref[2], ">", "", -1)
			titleFinal := strings.Replace(titleTmp, "</a", "", -1)
			linkUrl := siteUrl + doubleQuoteSplitHref[1]
			allUrlsAndTitles[linkUrl] = titleFinal
		}
	}
	return allUrlsAndTitles, nil
}

func main() {
	var dataDirectory, gzhodanConfig, multiDaily, noGzhodanConfig, optionalConfigs, tokensFile, markdownOnly, cliOnly, verboseOutput string
	flag.StringVar(&noGzhodanConfig, "g", "", "Use internally hardcoded configurations")
	flag.StringVar(&gzhodanConfig, "G", "gzhodan.conf", "Provide a Gzhodan configuration file!")
	flag.StringVar(&optionalConfigs, "O", "", "Optional configuration files seperated with a comma")
	flag.StringVar(&dataDirectory, "o", "", "Directory for which previous and new data is read and written to")
	flag.StringVar(&multiDaily, "m", "", "If application is running multiple times per day this is REQUIRED flag!")
	flag.StringVar(&tokensFile, "t", "", "If Gzhodan requires custom tokens -- not compatible with -g or -G !!!")
	flag.StringVar(&markdownOnly, "M", "", "Verbose output is combinable with -V for verbose")
	flag.StringVar(&cliOnly, "C", "", "CLI only output is combinable with -V for verbose")
	flag.StringVar(&verboseOutput, "V", "", "Verbose output is combinable with -C or -M")
	flag.Parse()

	args, argsLen := os.Args, len(os.Args)
	if argsLen > 2 {
		flag.PrintDefaults()
		err := fmt.Errorf("lack of arguments provided")
		checkError(err, 0, 0)
		os.Exit(1)
	}

	// Everything below need refactored into a method
	appStartTime := time.Now()
	dateFormatted := appStartTime.Format("2006-01-01")
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(dateFormatted)
	nameBuilder.WriteString(".log")

	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	app := Application{}
	app.multiDaily = false
	app.noGzhodanConfig = false
	app.statistics = Statistics{}
	app.statistics.appStartTime = appStartTime.UTC()
	app.tmpDir = os.TempDir()
	app.statistics.operatingSystem = runtime.GOOS
	app.statistics.date = appStartTime
	app.statistics.year = ""

	err := app.handleArgs(args, argsLen)
	if err != nil {
		checkError(err, 0, 0)
	}

	// Everything below could be in its own method
	//
	// Later functionality when there is alot data at some point we need condensing or checking
	// to or with a historicData file, these kind of programs need data regression to best dataset (size, quality, parsability,etc)
	// Do not remove
	err = app.checkPrevRuntimes()
	if err != nil {
		checkError(err, 0, 0)
	}

	err = initaliseLogging()
	if err != nil {
		checkError(err, 0, 0)
	}
	InfoLogger.Printf("Logging initialised")

	err = mkDirAndCD(app.appDir)
	if err != nil {
		checkError(err, 0, 0)
	}

	dirTree := []string{"test", "logs", "newletters", app.statistics.year}
	err = mkAppDirTree(app.appDir, dirTree)
	if err != nil {
		checkError(err, 0, 0)
	}
	app.testDir = filepath.Join(app.appDir, "test")

	urlsToVisit, baseDNSurlTotals, err := marshalURLsToMap()
	if err != nil {
		checkError(err, 0, 0)
	}
	allBaseUrlsSeq := make([]string, 0, len(urlsToVisit))
	for _, value := range urlsToVisit {
		allBaseUrlsSeq = append(strings.Split(value, ""))
	}

	totalUrls := 0
	app.statistics.originalDomains = len(baseDNSurlTotals) - 1

	for _, val := range baseDNSurlTotals {
		totalUrls = +val
	}

	app.statistics.originalUrls = totalUrls
	app.statistics.totalUrlsVisited += totalUrls
	// Everything above could be in its own method

	// Consider implecations of full implecation at some point...
	var defaultThreadCount int = 10
	threadCount := defaultThreadCount

	basePagesStdoutMap, err := curlNewBasePages(allBaseUrlsSeq)
	if err != nil {
		checkError(err, 0, 0)
	}

	artefactsFromBasePages, err := parseAllBasePagesForLinksAndTitles(basePagesStdoutMap)
	if err != nil {
		checkError(err, 0, 0)
	}

	failedLinksAndTitleByDomainMap := make(map[string]map[string]string)
	foundHistoricLinks := make(map[string]map[string]string)

	if app.historicDataFilePath != "" {
		foundBaseLinks, foundHistoricLinks, err := app.compareUrlsHistorically(artefactsFromBasePages)
		if err != nil {
			checkError(err, 0, 0)
		}
	} else {
		foundBaseLinks := basePagesStdoutMap
		WarningLogger.Printf("No historic data file provided to compare new url with previously enumerated data, this may take a lot longer!")
	}

	// Collect Duplicate URLs
	if foundHistoricLinks != nil {
		err := appendFalsePositiveDataToNextLTS(foundHistoricLinks)
		if err != nil {
			checkError(err, 0, 0)
		}
	}

	// Load tokens into memory
	tokensBuffer, tokenBufferSize, err := app.loadTokensIntoMem()
	//
	//
	// Need a thread flag
	// Need a max threads for maths
	offsets, err := assignTokenBufferOffsets(tokenBufferSize, threadCount)
	if err != nil {
		checkError(err, 0, 0)
	}
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
	//
	//  	artefactsFromBasePages[key]
	//  	artefactsFromBasePages[]value
	//
	go func() {
		titleCheckResult, err = parseTitles(tokenOffset)
		if !titleCheckResult {
			err := addValueToNestedStrStrMap(failedLinksAndTitleByDomainMap, url, url, title)
			if err != nil {
				checkError(err, 0, 0)
			}
		}
	}()

	assignTokenBufferOffsets() // array -> tokenBufferthreadId
	go func() {
		page, err := curlNewArticle(url)
		if err != nil {
			checkError(err, 0, 0)
		}
		matchedTokens, goodPage, err := parsePage(page, tokenOffset)
		if err != nil {
			checkError(err, 0, 0)
		}
		err := marshallParserResults(goodPage, matchedTokens, url)
		if err != nil {
			checkError(err, 0, 0)
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

	// Another kick in the really would like gzhobin data files
	backupDataStorage()

	// Another kick the really would like gzhobin data files
	updateDataStorage()

	// Output cli, file and (backup and then) organise historic data
	err = selectOutput(outputArgs)
	if err != nil {
		checkError(err, 0, 0)
	}

}

// REASONS keys for failed-map so that it makes sense

func marshallParserResults(goodPage bool, matchedTokens string, url string) error {
	if !goodPage {
		// remove url from queue,
		addValueToNestedStrStrMap(failedLinksAndTitleByDomainMap, "Parsed-Page-Results-Negative", url, titles)
	} else {

	}
	return nil
}

func ingestPagesTokensMap(tokens []byte) (map[string]map[int]string, error) {
	arstechnicaTokensMap, thehackernewsTokensMap, sansTokensMap := make(map[int]string), make(map[int]string), make(map[int]string), make(map[int]string)
	// do stuff

	resultMap := make(map[string]map[int]string)
	resultMap["arstechnica.com"], resultMap["thehackernews.com"], resultMap["www.sans.org"] = arstechnicaTokensMap, thehackernewsTokensMap, sansTokensMap

	return resultMap, nil
}

func parsePageForTokens(domain, page string) error {
	domain, err := urlKeyToDomainString(key)
	if err != nil {
		checkError(err, 0, 0)
	}
	webPageBuffer := bytes.NewBuffer([]byte(page))
	switch domain {
	case "arstechnica.com":
		arstechnicaTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
		if err != nil {
			checkError(err, 0, 0)
		}
	case "thehackernews.com":
		thehackernewsTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
		if err != nil {
			checkError(err, 0, 0)
		}
	case "www.sans.org":
		sansTokensMap, err = gzlopBuffer(webPageBuffer, tokens)
		if err != nil {
			checkError(err, 0, 0)
		}
	case "":
		err := fmt.Errorf("strange race condition occur with domain variable being an empty string")
		if err != nil {
			checkError(err, 0, 0)
		}
	default:

	}
	//  domain = "" // for loop from previous
	return nil
}

func wtfisthislackofsleepin2024() {

	searchTokens, err := os.ReadFile(tokensFile)
	if err != nil {
		checkError(err, 0, 0)
	}
	// TODO

	// map memory - for the same ~~paragraph~~ search for dates, url and tokens
	// soup go gets all the fields that have urls like gospider (CHECK HOW THAT WORK and do it locally)
	// gzlop buffer can then be adapter to search the buffer from address to offset for EVEN MORE SPEED

	// Naive Search for a token

	allTheArtefacts := make(map[int]map[int]string)
	for _, token := range searchTokens {
		artifacts, err := gzlopBuffer(file, token)
		if err != nil {
			checkError(err, 0, 0)
		}
		// WTF
		// allTheArtefacts[] = artifacts
	}
	//
	// BRAIN NEED THUNK HERE
	//
	return nil

}

// -
// -
// OUTPUT TODO
// -
// -

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

// defaultOutput - cli some app.statistics.stics and markdown report

func defaultOutput() {

}
