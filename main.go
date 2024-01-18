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

	"github.com/anaskhan96/soup"
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

func xtdFFGetNewPages(saveDirectory string, urlArr []string) error {
	xdotoolHandle := "xdotool"
	xdtOpenTerminalAndFirefox := " key ctrl+alt+t sleep 1 type firefox Enter"
	xdtFindFirefox := " search --onlyvisible --name firefox | head -n 1"
	xdtGoToURLinFirefox := " key \"ctrl+l\" type " // needs xdtEnterKey
	
	xdtEnterKey := " Enter"
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
	allTheArtefacts := make(map[int]map[int]string)
  	for _, entry := range entries {
       todaysInitialPages = append(files, entry.Name())
    }
	for _,pathToFile := range todaysInitialPages {
		data, err := ioutil.ReadFile(pathToFile)
		checkError(err)
		defer file.Close()

		// TODO
		// map memory - for the same ~~paragraph~~ search for dates, url and tokens 
		// soup go gets all the fields that have urls like gospider (CHECK HOW THAT WORK and do it locally)
		// gzlop buffer can then be adapter to search the buffer from address to offset for EVEN MORE SPEED
		buffer := bytes.NewBuffer(data)
		doc := soup.HTMLParse(string(data))

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

