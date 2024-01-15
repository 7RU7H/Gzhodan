package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Map = domain + id : url
func marshalURLsToMap() (map[string]string, map[string]int, error) {
	file, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
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

func xtdFFGetNewPages(saveDirectory string, urlArr[]string) error {
	xdotoolHandle := "xdotool" 
	xdtOpenTerminalAndFirefox := " key ctrl+alt+t sleep 1 type firefox Enter"
	xdtFindFirefox := " search --onlyvisible --name firefox | head -n 1"
	xdtGoToURLinFirefox := " key \"ctrl+l\" type " // needs xdtEnterKey
	xdtEnterKey := " Enter" 
  	xdtSavePageToPath := " key \"ctrl+s\" sleep 2 type " // needs xdtEnterKey
	xdtCloseFirefox := " key --clearmodifiers \"ctrl+F4\""
	//xdtFindSaveAsBox := " search --name \"Save as\""
	
	initXdoTool := exec.Command(xdotoolHandle, xdtOpenTerminalAndFirefox)
	err := initXdoTool.Run()
	checkError(err)	
	initXdoTool.Stdin = strings.NewReader(xdotoolHandle, xdtFindFirefox)
	err := initXdoTool.Run()
	checkError(err)	
	for _,url := range urlArr { 
		initXdoTool.Stdin = strings.NewReader(xdotoolHandle, xdtGoToURLinFirefox, url, xdtEnterKey)
		err := initXdoTool.Run()
		checkError(err)	
		initXdoTool.Stdin = strings.NewReader(xdotoolHandle, xdtSavePageToPath, saveDirectory, xdtEnterKey)
		err := initXdoTool.Run()
		checkError(err)	
	}
	initXdoTool.Stdin = strings.NewReader(xdotoolHandle, xdtCloseFirefox, xdtEnterKey)
	err := initXdoTool.Run()
	checkError(err)	

	return nil
}
	

func main() {
	urlsToVist, baseDNSurlTotals, err := marshalURLsToMap()
	checkError(err)

	allBaseUrlsSeq := make([]string, 0, len(urlsToVist))
	for k := range urlsToVist {
		allBaseUrlsSeq = append(allBaseUrlsArr, k)
	}

	totalUrls := 0
	totalDomains := len(baseDNSurlTotals)-1

	for _,val := range baseDNSurlTotals {
		totalUrls =+ val
	}

	err := xtdFFGetNewPages(saveDirectory, allBaseUrlsArr)

	
	
	entries, err := os.ReadDir(saveDirectory)
	checkError(err)
 	var todaysInitialPages []string
  	for _, entry := range entries {
       todaysInitialPages = append(files, entry.Name())
    }
	for _,pathtofile := range todaysInitialPages {
		file, err := os.Open()
		checkError(err)
		defer file.Close()
		doc, err := goquery.NewDocumentFromReader(file)
		//
		// BRAIN NEED THUNK HERE
		//

	}

}
