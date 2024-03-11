package main

import (
	"bufio"
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

func siteSpecificsHandler(domain string) {
	switch domain {
	case "arstechnica.com":
	case "portswigger.net":
	case "thehackernews.com":
	case "www.sans.org":
	default:

	}
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

func curlNewArticles(urlArr []string) error {
	var args string = "-K curlrc -L -O"
	urlsStr := strings.Join(urlArr, " ")
	runCurl := exec.Command("curl", args, urlsStr)
	err := runCurl.Run()
	checkError(err)
	return nil
}

func findLinksAndTitlesFromBasePages(basePagesStdoutMap map[string]string) error {
	return nil
}

func main() {
	appDir, err := os.Getwd()
	checkError(err)
	stat := Statistics{}
	stat.os = runtime.GOOS
	stat.tmpDir = os.TempDir()
	stat.appDir = appDir
	now := time.Now().UTC()
	stat.date = now.Format("2006-01-01")
	stat.year = strconv.Itoa(now.Year())
	err = checkPrevRuntimes(stat.appDir, stat.date)
	checkError(err)
	dirTree := []string{"test", "logs", "newletters", stat.year}
	err = mkAppDirTree(stat.appDir, dirTree)
	checkError(err)
	stat.testDir = filepath.Join(appDir, "test")
	err = initaliseLogging()
	checkError(err)
	InfoLogger.Printf("Logging initialised")

	err = mkDirAndCD(stat.date)
	checkError(err)
	saveDirectory := stat.date

	urlsToVisit, baseDNSurlTotals, err := marshalURLsToMap()
	checkError(err)

	allBaseUrlsSeq := make([]string, 0, len(urlsToVisit))
	for _, value := range urlsToVisit {
		allBaseUrlsSeq = append(strings.Split(value, ""))
	}

	totalUrls := 0
	totalDomains := len(baseDNSurlTotals) - 1

	for _, val := range baseDNSurlTotals {
		totalUrls = +val
	}

	stat.originalDomains = totalDomains
	stat.originalUrls = totalUrls
	stat.totalUrlsVisited += totalUrls

	// stdout -> 4 base pages
	basePagesStdoutMap, err := curlNewBasePages(allBaseUrlsSeq)
	checkError(err)

	findLinksAndTitlesFromBasePages(basePagesStdoutMap)

	// regexp links and page titles
	//  ---- Domain Specifics:
	// portswigger -> links are just title strings.Join(titleNoAtags, "-")
	// sans
	// arstechnica
	// thehackernews
	//
	// compare maps for domain against previous enumerated list file with gzlop
	err = compareBasePagesToHistoricData()
	// ---- only need to store and compare urls
	// if in the file remove from map
	// Storage 2 files one .csv per run and collective with Page rating, time, url, matched tokens, And just previous-urls-found-only.list
	//
	// Get new Pages
	err = curlNewArticles(finalUrlsArr)
	// Print Alert - similiar to each row of .csv of urls

	// Where the funky code really begins
	entries, err := os.ReadDir(saveDirectory)
	checkError(err)

}

func compareBasePagesToHistoricData(historicUrlsFile, tokensFile string, urlsFound map[string]string) error {
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

}
