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
func marshalURLsToMap() (map[string]string, error) {
	file, err := os.Open("urls.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urlsMapped := make(map[string]string)
	urlStr := ""
	domainIDCounter := 1

	for scanner.Scan() {
		urlStr = scanner.Text()
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("Invalid URL: %s\n", urlStr)
			continue
		}

		hostname := parsedURL.Hostname()
		if _, ok := urlsMapped[hostname]; !ok {
			urlsMapped[hostname] = strconv.Itoa(domainIDCounter)
			domainIDCounter++
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

	}

	return urlsMapped, nil
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
	urlsToVist, err := marshalURLsToMap()
	checkError(err)

	cmdArgsBuilder := strings.Builder{}

	curlArgs := "-X GET -A Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 "
	curlLimitRateFlag := "--limit-rate 10000B "
	curlOutputFlag := "-o "

	for _, url := range urlsToVist {
		cmdArgsBuilder.WriteString(curlArgs)
		cmdArgsBuilder.WriteString(curlLimitRateFlag)
		cmdArgsBuilder.WriteString(url)
		cmdArgsBuilder.WriteString(" ")
		cmdArgsBuilder.WriteString(curlOutputFlag)
		cmdArgsBuilder.WriteString("test.txt")
		execCurl(cmdArgsBuilder.String())
		cmdArgsBuilder.Reset()
	}

}
