package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"log"
	"net/url"
	"strconv"

	"Gzhodan/storage.go"
)
func checkError(err error) {
	if err != nil {
    	log.Fatal(err)
    }
}

// Map = domain + id : url
func marshalURLsToMap(urlsFile) (map[string]string, error) { 
		file, err := os.Open(urlsFile)
		if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	urlsSlice := make([]string, 0, 0)
	urlsMapped := make(map[string]string)
	domainIDCounter := 1
	
	for scanner.Scan() {
		urlSlice = append(scanner.Text())
	
	if err := scanner.Err(); err != nil {
			panic(err)
	}
	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil {
			panic(err)
		}

		hostname := parsedURL.Hostname()
		if _, ok := urlsMapped[hostname]; !ok {
			urlsMapped[hostname] = strconv.Itoa(id)
			domainIDCounter++
		}
	}
	
	return urlsMapped, nil
}

func main() {

	
}