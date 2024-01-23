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
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type Statistics struct {
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

func checkError(err error) error {
	if err != nil {
		fmt.Errorf("%s", err)
		log.Fatal(err)
		panic(err)
	}
	return err
}

func createFile(filepath string) error {
	filePtr, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Creation Error:", err)
		//log.Fatal(err);
	}
	defer filePtr.Close()
	return nil
}

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	checkError(err)
	if os.IsNotExist(err) {
		log.Fatal("File path does not exist")
		return false, err
	}
	return true, nil
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

func softConfFFToSaveAlwaysHTMLOnly(testDir string, recursionCounter int) (int, error) {
	pageName := "test.html"
	xdotoolHandle := "xdotool"
	xdtOpenTerminalAndFirefox := " key ctrl+alt+t sleep 1 type firefox Enter"
	xdtFindFirefox := " search --onlyvisible --name firefox | head -n 1"
	xtdClick := " key click 1 "                  //
	xtdDown := " key Down "                      //
	xdtTab := " key Tab "                        //
	xdtGotoFileNaming := " key \"ctrl+l\" type " // needs xdtEnterKey
	xdtEnterKey := " Return"
	xdtSavePageToPath := " key \"ctrl+s\" sleep 2 type " // needs xdtEnterKey
	xdtCloseFirefox := " key --clearmodifiers \"ctrl+F4\""

	subCmdBuilder := strings.Builder{}
	initXdoTool := exec.Command(xdotoolHandle, xdtOpenTerminalAndFirefox)
	err := initXdoTool.Run()
	checkError(err)
	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtFindFirefox)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtSavePageToPath)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	// Click 1 // to escape writing input for file name - we will return with CTRL+l to then retype any accidentally click input into the Name: input bar
	subCmdBuilder.WriteString(xtdClick)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	// Tab 1 // To move gui to Dropdown on save All, HTML only, text
	subCmdBuilder.WriteString(xdtTab)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	// Down // Selecting Save only html
	subCmdBuilder.WriteString(xtdDown)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtGotoFileNaming)
	subCmdBuilder.WriteString(testDir)
	subCmdBuilder.WriteString(pageName)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	subCmdBuilder.Reset()

	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtCloseFirefox)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	subCmdBuilder.Reset()

	xdtAndFFSaveProperly := false
	// check if _files or .txt
	// delete files

	if recursionCounter != 6 && xdtAndFFSaveProperly != true {
		recursionCounter, err = softConfFFToSaveAlwaysHTMLOnly(testDir, recursionCounter)
	}

	return recursionCounter, nil
}

func xtdFFGetNewPages(saveDirectory string, urlArr []string) error {
	xdotoolHandle := "xdotool"
	xdtOpenTerminalAndFirefox := " key ctrl+alt+t sleep 1 type firefox Enter"
	xdtFindFirefox := " search --onlyvisible --name firefox | head -n 1"
	xdtGoToURLinFirefox := " key \"ctrl+l\" type " // needs xdtEnterKey

	xdtEnterKey := " Return"
	xdtSavePageToPath := " key \"ctrl+s\" sleep 2 type " // needs xdtEnterKey
	xdtCloseFirefox := " key --clearmodifiers \"ctrl+F4\""

	subCmdBuilder := strings.Builder{}
	initXdoTool := exec.Command(xdotoolHandle, xdtOpenTerminalAndFirefox)
	err := initXdoTool.Run()
	checkError(err)
	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtFindFirefox)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	checkError(err)
	subCmdBuilder.Reset()

	for _, url := range urlArr {
		subCmdBuilder.WriteString(xdotoolHandle)
		subCmdBuilder.WriteString(xdtGoToURLinFirefox)
		subCmdBuilder.WriteString(url)
		subCmdBuilder.WriteString(xdtEnterKey)
		initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
		err = initXdoTool.Run()
		checkError(err)
		subCmdBuilder.Reset()

		subCmdBuilder.WriteString(xdotoolHandle)
		subCmdBuilder.WriteString(xdtSavePageToPath)
		subCmdBuilder.WriteString(saveDirectory)
		// page name
		subCmdBuilder.WriteString(xdtEnterKey)
		initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
		err = initXdoTool.Run()
		subCmdBuilder.Reset()
	}

	subCmdBuilder.WriteString(xdotoolHandle)
	subCmdBuilder.WriteString(xdtCloseFirefox)
	subCmdBuilder.WriteString(xdtEnterKey)
	initXdoTool.Stdin = strings.NewReader(subCmdBuilder.String())
	err = initXdoTool.Run()
	subCmdBuilder.Reset()

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

func initaliseLogging() error {
	now := time.Now().UTC()
	dateFormatted := now.Format("2006-01-01")
	nameBuilder := strings.Builder{}
	nameBuilder.WriteString(dateFormatted)
	nameBuilder.WriteString(".log")
	file, err := os.OpenFile(nameBuilder.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
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

func main() {
	appDir := "/tmp" // replace with args flag for directory
	stat := Statistics{}
	now := time.Now().UTC()
	stat.date = now.Format("2006-01-01")
	stat.year = strconv.Itoa(now.Year())
	err := checkPrevRuntimes(appDir, stat.date)
	dirTree := []string{"test", "logs", "newletters", stat.year}
	err = mkAppDirTree(appDir, dirTree)
	checkError(err)
	testDirFP := filepath.Join(appDir, "test")
	softConfFFToSaveAlwaysHTMLOnly(testDirFP, 0)
	checkError(err)
	mkDirAndCD(stat.date)

	err = initaliseLogging()
	checkError(err)
	InfoLogger.Println("Logging initialised")

	//    InfoLogger.Println("Something noteworthy happened")
	//    WarningLogger.Println("There is something you should know about")
	//    ErrorLogger.Println("Something went wrong")

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

	err = xtdFFGetNewPages(saveDirectory, allBaseUrlsSeq)


	// Where the funky code really begins
	entries, err := os.ReadDir(saveDirectory)
	checkError(err)
	var todaysInitialPages []string
	allTheArtefacts := make(map[int]map[int]string)
	for _, entry := range entries {
		todaysInitialPages = append(files, entry.Name())
	}
	for _, pathToFile := range todaysInitialPages {
		file, err := os.ReadFile(pathToFile)
		checkError(err)
		defer file.Close()

		// TODO
		// map memory - for the same ~~paragraph~~ search for dates, url and tokens
		// soup go gets all the fields that have urls like gospider (CHECK HOW THAT WORK and do it locally)
		// gzlop buffer can then be adapter to search the buffer from address to offset for EVEN MORE SPEED
		buffer := bytes.NewBuffer()

		doc := soup.HTMLParse(string(file))

		// Naive Search for a token
		for _, token := range searchTokens {
			artifacts, err := gzlopBuffer(buffer, token)
			checkError(err)
			allTheArtefacts[token] = artifacts
		}
		//
		// BRAIN NEED THUNK HERE
		//

	}

}

// thunkage for soup ++Brain
// Account for site differences
// Account for previously ran urls - they have to be unique so no worries

// What are we looking for in html? - ANS: urls, titles
// find artefacts containing Tokens 
// Store:
// id index_of_next_chunk  time + domain + urlused + artefacts { urlfound  + titles; artefacts + tokensMatched } terminationStub; 


datastream := marshalArtefactstoDS()
func writeNewGzhodinToDisk(datastream []bytes) error {
	extension := ".bin.gzhobin"
	gzhobinName := ""
	file, _ := os.Create()
	defer file.Close()
	

	for _,data := range datastream { 	
		err := binary.Write(file, binary.LittleEndian, data)
		checkError(err)
	}
}

// It would be awesome to have a uniform binary file format so that hardcoded address exists at offsets to a being of a file

func loadLastGzhobinToMem() error {

}


// GzhobinTemplate needs to be a array of hardcoded offsets   
func compareCurrAndLastGzhobins() error {


	err := binary.Read(lastFile, binary.LittleEndian, XXXX)
	switch err {
	case io.EOF: 
		break
	case !nil:
		checkError(err)
	default:
		continue
	}

}


func siteSpecificsHandler(domain string) {
	switch domain {
	case "arstechnica.com":
	case "news.ycombinator.com":
	case "portswigger.net":
	case "thehackernews.com":
	case "www.sans.org":
	default:

	}
}

