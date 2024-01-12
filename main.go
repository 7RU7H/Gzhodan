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

// Because is it used by people
func execCurl(args string) error {
	cmd := exec.Command("curl", args)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func main() {
	urlsToVist, baseDNSurlTotals, err := marshalURLsToMap()
	checkError(err)

	allBaseUrlsArr := make([]string, 0, len(urlsToVist))
	for k := range urlsToVist {
		allBaseUrlsArr = append(allBaseUrlsArr, k)
	}

	// Human spawl linear request by link domain/url, if possible from a base url
	// No-touching-The-Sides algo - a = 3,b = 4,c = 5,d = 1, e = 1 -> a1,b1,c1,a2,b2,c2,b3,c3,d1,e1,a1,b1,c4

	queue := make([]string, 0, len(allBaseUrlsArr))

	cmdArgsBuilder := strings.Builder{}

	curlArgs := "-X GET -A Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 "
	//	curlLimitRateFlag := "--limit-rate 10000B "
	curlOutputFlag := "-o "

	// Replace with queue
	for _, url := range urlsToVist {
		cmdArgsBuilder.WriteString(curlArgs)
		//	cmdArgsBuilder.WriteString(curlLimitRateFlag)
		cmdArgsBuilder.WriteString(url)
		cmdArgsBuilder.WriteString(" ")
		cmdArgsBuilder.WriteString(curlOutputFlag)
		cmdArgsBuilder.WriteString("test.txt")
		execCurl(cmdArgsBuilder.String())
		cmdArgsBuilder.Reset()
	}

}
