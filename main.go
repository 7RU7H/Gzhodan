package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
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
	// optionalConfigs      map[string]string // ASP
	tokensFile string
}

type Statistics struct {
	operatingSystem  string
	originalDomains  int
	originalUrls     int
	totalUrlsVisited int
	date             time.Time
	year             string
	appStartTime     time.Time
	totalFailedUrls  int
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

// Simple idiot solution
func (b *CircularBuffer) assignReadPointerOffsets(concurrencyOffset, remainder int) {
	if remainder != 0 {
		for i := 1; i <= len(b.readPointers); i += concurrencyOffset {
			b.readPointers[i] = (concurrencyOffset * i) + remainder
			remainder--
		}
	} else {
		for i := 1; i <= len(b.readPointers); i += concurrencyOffset {
			b.readPointers[i] = concurrencyOffset * i
		}
	}
}

type MatchOnTitles struct {
	url    string
	titles string
	tokens []string
	count  uint
}

func newMatchOnTitlesBuilder() *MatchOnTitles {
	return &MatchOnTitles{}
}

func (m *MatchOnTitles) Url(url string) *MatchOnTitles {
	m.url = url
	return m
}
func (m *MatchOnTitles) Titles(titles string) *MatchOnTitles {
	m.titles = titles
	return m
}

func (m *MatchOnTitles) Build() MatchOnTitles {
	return MatchOnTitles{
		url:    m.url,
		titles: m.titles,
		tokens: make([]string, 0),
		count:  0,
	}
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func (a *Application) CreateWorkingDir() error {
	err := mkDirAndCD(a.appDir)
	if err != nil {
		checkError(err, 0, 0)
		return err
	}

	dirTree := []string{"test", "logs", "newletters", a.statistics.year}
	err = mkAppDirTree(a.appDir, dirTree)
	if err != nil {
		checkError(err, 0, 0)
		return err
	}
	a.testDir = filepath.Join(a.appDir, "test")
	return nil
}

func (app *Application) processCurrAndHistoricData(artefactsFromBasePages map[string]string) (foundBaseLinksAndTitles, failedLinksAndTitleByDomainMap map[string]map[string]string, err error) {
	if app.historicDataFilePath != "" {
		foundBaseLinksAndTitles, foundHistoricLinks, err := app.compareUrlsHistorically(artefactsFromBasePages)
		if err != nil {
			checkError(err, 0, 0)
			return nil, nil, err
		}
		for key := range foundHistoricLinks {
			for subKey, value := range foundHistoricLinks[key] {
				addValueToNestedStrStrMap(failedLinksAndTitleByDomainMap, key, subKey, value)
			}
		}
		return foundBaseLinksAndTitles, failedLinksAndTitleByDomainMap, nil
	} else {
		for urlKey, titlesValue := range artefactsFromBasePages {
			domain, err := urlKeyToDomainString(urlKey)
			if err != nil {
				checkError(err, 0, 0)
				return nil, nil, err
			}
			addValueToNestedStrStrMap(foundBaseLinksAndTitles, domain, urlKey, titlesValue)
			domain = ""
		}
		WarningLogger.Printf("No historic data file provided to compare new url with previously enumerated data, this may take a lot longer!")
		return foundBaseLinksAndTitles, failedLinksAndTitleByDomainMap, nil
	}
}

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

	bsLen := len(tokensFileAsBytes)
	return tokensFileAsBytes, bsLen, nil
}

func mvToAppDir(appDir string) (bool, error) {
	var result bool = false
	err := os.Chdir(appDir)
	if err != nil {
		checkError(err, 0, 0)
		return result, err
	}
	result = true
	return result, nil
}

func (a *Application) handleArgs(args []string, argsLength int) error {
	for i := 0; i <= argsLength-1; i++ {
		switch args[i] {
		case "-o":
			appDirExists, err := checkFileExists(args[i+1])
			if err != nil || !appDirExists {
				checkError(err, 0, 0)
				WarningLogger.Printf("No historic data file provided to compare new url with previously enumerated data, this may take a lot longer!")
			} else {
				a.historicDataFilePath = args[i+1]
			}
			a.appDir = args[i+1]
			mvStat, err := mvToAppDir(args[i+1])
			if err != nil || !mvStat {
				checkError(err, 0, 0)
				WarningLogger.Printf("unable to move to the specified application directory: %s", args[i+1])
			} else {
				wd, err := os.Getwd()
				if err != nil || wd == "" {
					checkError(err, 0, 0)
					continue
				}
				InfoLogger.Printf("current application directory changed to %s", wd)
			}
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

func (a *Application) selectOutput(dut map[string]map[string]string, mtt map[string]*MatchOnTitles, failedCount int) error {
	argsSize := len(a.outputType)
	var argsId int = 0
	switch argsSize {
	case 3:
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
			err := fmt.Errorf("invalid output arguments provide: %v ; from slice of size: %v", a.outputType, argsSize)
			checkError(err, 0, 0)
			return err
		}
	case 1:
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
	case 0:
		argsId = 0
	default:
		err := fmt.Errorf("invalid output arguments provide: %v of size %v", a.outputType, argsSize)
		checkError(err, 0, 0)
		return err
	}

	switch argsId {
	case 1: // verbose
		verboseOutput(a, dut, mtt, failedCount)
	case 2: // cli only
		cliOnlyOutput(a, dut, failedCount)
	case 3: // verbose cli only
		verboseCliOutput(a, dut, mtt, failedCount)
	case 5: // markdown only
		markdownOnlyOutput(a, dut, failedCount)
	case 6: // verbose markdown
		verboseMarkdownOutput(a, dut, mtt, failedCount)
	case 0:
		defaultOutput(a, dut, failedCount)
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
		allUrlsAsBytes = strconv.AppendQuote(allUrlsAsBytes, urlKey)
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
		err := fmt.Errorf("invalid url: %s", keyUrl)
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
			io.WriteString(os.Stdout, fmt.Sprintf("Invalid URL: %s\n", urlStr))
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
	checkError(err, 0, 0)
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
		checkError(err, 0, 0)
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

func addValueToNestedStrIntStrMap(parentMap map[string]map[int]string, parentKey string, childKey int, nestedValue string) {
	childMap := parentMap[parentKey]
	if childMap == nil {
		childMap = make(map[int]string)
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
				basePageAllHrefs = append(basePageAllHrefs, line)
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

func defangUrl(inputUrl string) string {
	tmpUrl := strings.ReplaceAll(inputUrl, "http", "hxxp")
	tmpUrl = strings.ReplaceAll(tmpUrl, "://", "[://]")
	outputUrl := strings.ReplaceAll(tmpUrl, ".", "[.]")
	return outputUrl
}

// https://gosamples.dev/write-file/
func updateDataStorage(file string, passed map[string]map[string]string, failed map[string]map[string]string) error {
	f, err := os.Open(file)
	if err != nil {
		checkError(err, 0, 0)
	}
	defer f.Close()
	for key := range passed {
		for subKey := range passed[key] {
			_, err := f.WriteString(subKey + "\n")
			if err != nil {
				checkError(err, 0, 0)
			}
		}
	}
	for key := range failed {
		for subKey := range failed[key] {
			_, err := f.WriteString(subKey + "\n")
			if err != nil {
				checkError(err, 0, 0)
			}
		}
	}
	return nil
}

// https://zetcode.com/golang/copyfile/
func backupDataStorage(src string) error {
	fin, err := os.Open(src)
	if err != nil {
		checkError(err, 0, 0)
		err = createFile(src)
		checkError(err, 0, 0)
		fin, err = os.Open(src)
		checkError(err, 0, 0)
	}
	defer fin.Close()

	dst := fin.Name() + ".bak"
	fout, err := os.Create(dst)
	if err != nil {
		checkError(err, 0, 0)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		checkError(err, 0, 0)
	}
	return nil
}

func verboseOutput(app *Application, domainUrlTitles map[string]map[string]string, matchedTitles map[string]*MatchOnTitles, failedCount int) {
	verboseCliOutput(app, domainUrlTitles, matchedTitles, failedCount)
	verboseMarkdownOutput(app, domainUrlTitles, matchedTitles, failedCount)
}

func cliOnlyOutput(app *Application, domainUrlTitles map[string]map[string]string, failedCount int) {
	successfulUrls := len(domainUrlTitles)
	for key := range domainUrlTitles {
		for subKey, value := range domainUrlTitles[key] {
			io.WriteString(os.Stdout, fmt.Sprintf("Found titled: \"%s\" at URL: %s\n", defangUrl(subKey), value))
		}
	}
	io.WriteString(os.Stdout, fmt.Sprintf("Successful URLs found: %v\n", successfulUrls))
	io.WriteString(os.Stdout, fmt.Sprintf("Failed URLs (As of previous runtime file: %v) refound: %v\n", app.historicDataFilePath, failedCount))

}

func lsCdTouchMarkdownFile(appDir string, date time.Time) (*os.File, error) {
	currDir, err := os.Getwd()
	checkError(err, 0, 0)
	if currDir != appDir {
		err := fmt.Errorf("current directory %v is not application specified directory: %v for some weird reason - must debug", currDir, appDir)
		checkError(err, 0, 0)
		return nil, err
	}
	err = os.Chdir("newsletters")
	if err != nil {
		checkError(err, 0, 0)
		return nil, err
	}

	mdFilename := date.Format(time.DateOnly) + ".md"
	exists, err := checkFileExists(mdFilename)
	checkError(err, 0, 0)
	if !exists {
		err := fmt.Errorf("file with filename %v already exists", mdFilename)
		checkError(err, 0, 0)
		return nil, err
	}
	fout, err := os.Create(mdFilename)
	if err != nil {
		checkError(err, 0, 0)
	}
	return fout, nil
}

func cdUpFromNewletters() error {
	err := os.Chdir("../")
	if err != nil {
		checkError(err, 0, 0)
	}
	return nil
}

func markdownOnlyOutput(app *Application, domainUrlTitles map[string]map[string]string, failedCount int) error {
	file, err := lsCdTouchMarkdownFile(app.appDir, app.statistics.date)
	if err != nil {
		checkError(err, 0, 0)
		return err
	}
	defer file.Close()
	successfulUrls := len(domainUrlTitles)
	for key := range domainUrlTitles {
		for subKey, value := range domainUrlTitles[key] {
			io.WriteString(file, fmt.Sprintf("Found titled: \"%s\" at URL: %s\n", defangUrl(subKey), value))
		}
	}
	io.WriteString(file, fmt.Sprintf("Successful URLs found: %v\n", successfulUrls))
	io.WriteString(file, fmt.Sprintf("Failed URLs (As of previous runtime file: %v) refound: %v\n", app.historicDataFilePath, failedCount))

	defer cdUpFromNewletters()
	return nil
}

func verboseCliOutput(app *Application, domainUrlTitles map[string]map[string]string, matchedTitles map[string]*MatchOnTitles, failedCount int) error {
	successfulUrls := len(domainUrlTitles)
	for key := range domainUrlTitles {
		for subKey, value := range domainUrlTitles[key] {
			io.WriteString(os.Stdout, fmt.Sprintf("Found titled: \"%s\" at URL: %s\n", defangUrl(subKey), value))
			mot := matchedTitles[key]
			io.WriteString(os.Stdout, fmt.Sprintf("Matched on %v Tokens: %v \n", mot.count, mot.tokens))
		}
	}
	io.WriteString(os.Stdout, fmt.Sprintf("\t---- Statistics - %v ----\n", app.statistics.date.String()))
	io.WriteString(os.Stdout, fmt.Sprintf("Started: %v\n", app.statistics.appStartTime))
	io.WriteString(os.Stdout, fmt.Sprintf("OS: %v\n", app.statistics.operatingSystem))
	io.WriteString(os.Stdout, fmt.Sprintf("Total URLs visited %v\n", app.statistics.totalUrlsVisited))
	io.WriteString(os.Stdout, fmt.Sprintf("Original Domains provided: %v\n", app.statistics.originalDomains))
	io.WriteString(os.Stdout, fmt.Sprintf("Original URLs provided: %v\n", app.statistics.originalUrls))
	io.WriteString(os.Stdout, fmt.Sprintf("Successful URLs found: %v\n", successfulUrls))
	io.WriteString(os.Stdout, fmt.Sprintf("Failed URLs (As of previous runtime file: %v) refound: %v\n", app.historicDataFilePath, failedCount))
	return nil
}

func verboseMarkdownOutput(app *Application, domainUrlTitles map[string]map[string]string, matchedTitles map[string]*MatchOnTitles, failedCount int) error {
	file, err := lsCdTouchMarkdownFile(app.appDir, app.statistics.date)
	if err != nil {
		checkError(err, 0, 0)
		return err
	}
	defer file.Close()
	successfulUrls := len(domainUrlTitles)
	for key := range domainUrlTitles {
		for subKey, value := range domainUrlTitles[key] {
			io.WriteString(file, fmt.Sprintf("Found titled: \"%s\" at URL: %s\n", defangUrl(subKey), value))
			mot := matchedTitles[key]
			io.WriteString(file, fmt.Sprintf("Matched on %v Tokens: %v \n", mot.count, mot.tokens))
		}
	}
	io.WriteString(file, fmt.Sprintf("\t---- Statistics - %v ----\n", app.statistics.date.String()))
	io.WriteString(file, fmt.Sprintf("Started: %v\n", app.statistics.appStartTime))
	io.WriteString(file, fmt.Sprintf("OS: %v\n", app.statistics.operatingSystem))
	io.WriteString(file, fmt.Sprintf("Total URLs visited %v\n", app.statistics.totalUrlsVisited))
	io.WriteString(file, fmt.Sprintf("Original Domains provided: %v\n", app.statistics.originalDomains))
	io.WriteString(file, fmt.Sprintf("Original URLs provided: %v\n", app.statistics.originalUrls))
	io.WriteString(file, fmt.Sprintf("Successful URLs found: %v\n", successfulUrls))
	io.WriteString(file, fmt.Sprintf("Failed URLs (As of previous runtime file: %v) refound: %v\n", app.historicDataFilePath, failedCount))

	defer cdUpFromNewletters()
	return nil
}

// defaultOutput - cli some app.statistics.stics and markdown report
func defaultOutput(app *Application, domainUrlTitles map[string]map[string]string, failedCount int) {
	cliOnlyOutput(app, domainUrlTitles, failedCount)
	markdownOnlyOutput(app, domainUrlTitles, failedCount)
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

	var gzlopTestBoolean bool = false

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

	err = app.checkPrevRuntimes()
	if err != nil {
		checkError(err, 0, 0)
	}

	err = initaliseLogging()
	if err != nil {
		checkError(err, 0, 0)
	}
	InfoLogger.Printf("Logging initialised")

	err = app.CreateWorkingDir()
	if err != nil {
		checkError(err, 0, 0)
	}

	urlsToVisit, baseDNSurlTotals, err := marshalURLsToMap()
	if err != nil {
		checkError(err, 0, 0)
	}
	allBaseUrlsSeq := make([]string, 0, len(urlsToVisit))
	for _, value := range urlsToVisit {
		allBaseUrlsSeq = append(allBaseUrlsSeq, value)
	}

	totalUrls := 0
	app.statistics.originalDomains = len(baseDNSurlTotals) - 1

	for _, val := range baseDNSurlTotals {
		totalUrls = +val
	}

	app.statistics.originalUrls = totalUrls
	app.statistics.totalUrlsVisited += totalUrls

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

	// map[string]map[string]string
	foundBaseLinksAndTitles, failedLinksAndTitleByDomainMap, err := app.processCurrAndHistoricData(artefactsFromBasePages)
	if err != nil {
		checkError(err, 0, 0)
	}

	tokensArray, tokensArrayLen, err := app.loadTokensIntoMem()
	if err != nil {
		checkError(err, 0, 0)
	}

	// DOUBLE CHECK!!
	workerCount := threadCount
	totalTokens := tokensArrayLen - 1
	remainder := totalTokens % workerCount
	concCurrOffset := 0
	if workerCount <= totalTokens {
		if workerCount > (totalTokens / 2) {
			concCurrOffset = totalTokens / workerCount
		}
	} else {
		workerCount = totalTokens
		concCurrOffset = 1
		remainder = 0
	}

	gzlopTestMap := make(map[string]map[int]string)
	if gzlopTestBoolean {
		for key := range foundBaseLinksAndTitles {
			for subKey, value := range foundBaseLinksAndTitles[key] {
				titlesAsBytes := make([]byte, 0)
				valueAsSlice := strings.SplitAfterN(value, " ", -1)
				for i := 0; i <= len(value)-1; i++ {
					titlesAsBytes = strconv.AppendQuote(titlesAsBytes, valueAsSlice[i])
				}
				result, err := gzlopBuffer(bytes.NewBuffer(titlesAsBytes), tokensArray)
				if err != nil {
					checkError(err, 0, 0)
				}
				for resultsKey, resultsVal := range result {
					addValueToNestedStrIntStrMap(gzlopTestMap, subKey, resultsKey, resultsVal)
				}
			}
		}
	}

	TokensBuffer := newCircularBuffer(tokensArray, concCurrOffset, workerCount)
	TokensBuffer.assignReadPointerOffsets(concCurrOffset, remainder)
	if err != nil {
		checkError(err, 0, 0)
	}

	titleTokeniserResults := make(map[string]*MatchOnTitles)
	motBuilder := newMatchOnTitlesBuilder()

	// TEST with gzlopBuffer for which is faster
	workerId := 0 // Work-around pre Paralellism
	// mutex foundBaseLinks
	// go func (workerId int, foundBaseLinks map[string]map[string]string) {
	for key := range foundBaseLinksAndTitles {
		for subKey, value := range foundBaseLinksAndTitles[key] {
			var bufferTokens, bufferTitles bytes.Buffer
			domUrlMoT := motBuilder.
				Url(subKey).
				Titles(value).
				Build()
			valueAsSlice := strings.SplitAfterN(value, " ", -1)
			titlesAsBytes := make([]byte, 0)
			for h := 0; h <= len(valueAsSlice)-1; h++ {
				titlesAsBytes = strconv.AppendQuote(titlesAsBytes, valueAsSlice[h])
			}

			var matchesFound, matchThreshold uint = 0, 5
			for i := 0; i <= tokensArrayLen-1; i++ {
				if matchThreshold != matchesFound {
					for j := 0; j <= len(titlesAsBytes)-1; j++ {
						bufferTokens.WriteByte(TokensBuffer.readCircularBufferFromOffset(workerId))
						bufferTitles.WriteByte(titlesAsBytes[j])
						tokenInBuf := bufferTokens.Bytes()
						titleWordInBuf := bufferTitles.Bytes()
						match := bytes.Compare(titleWordInBuf, tokenInBuf)
						if match != 0 {
							domUrlMoT.tokens = append(domUrlMoT.tokens, bufferTokens.String())
							matchesFound++
						}
					}
				} else {
					domUrlMoT.count = matchesFound
					break
				}
			}
			titleTokeniserResults[key] = &domUrlMoT
		}
	}
	// }()

	passedTokenisedLinksAndTitleByDomainMap := make(map[string]map[string]string)
	failedTokenisedLinksAndTitleByDomainMap := make(map[string]map[string]string)
	for key, value := range titleTokeniserResults {
		domain, err := urlKeyToDomainString(value.url)
		if err != nil {
			checkError(err, 0, 0)
		}
		switch value.count {
		case 0:
			addValueToNestedStrStrMap(failedTokenisedLinksAndTitleByDomainMap, domain, value.url, value.titles)
			delete(titleTokeniserResults, key)
		default:
			addValueToNestedStrStrMap(passedTokenisedLinksAndTitleByDomainMap, key, value.url, value.titles)
		}
	}

	for key := range passedTokenisedLinksAndTitleByDomainMap {
		for urlSubKey := range passedTokenisedLinksAndTitleByDomainMap[key] {
			InfoLogger.Printf("Attempting to curl: %s\n", urlSubKey)
			curlNewArticles(strings.SplitAfterN(urlSubKey, "", -1))
		}
	}

	err = backupDataStorage(app.historicDataFilePath)
	if err != nil {
		checkError(err, 0, 0)
	}

	err = updateDataStorage(app.historicDataFilePath, passedTokenisedLinksAndTitleByDomainMap, failedTokenisedLinksAndTitleByDomainMap)
	if err != nil {
		checkError(err, 0, 0)
	}

	app.statistics.totalFailedUrls = 0
	for _, key := range failedLinksAndTitleByDomainMap {
		app.statistics.totalFailedUrls = +len(key)
	}
	for _, key := range failedTokenisedLinksAndTitleByDomainMap {
		app.statistics.totalFailedUrls = +len(key)
	}

	err = app.selectOutput(passedTokenisedLinksAndTitleByDomainMap, titleTokeniserResults, app.statistics.totalFailedUrls)
	if err != nil {
		checkError(err, 0, 0)
	}

}
