


- use a database to store information
- configuration file for gzhodan
- make curlrc -> cmd.Exec more jitter friendly
- file hashing for comparisons


```go
func softConfFFToSaveAlwaysHTMLOnly(testDir string, recursionCounter int) (int, error) {
	pageName := "test"

	xdtOpenNewTermAndFirefox := "xdotool key \"ctrl+alt+t\" sleep 1 type firefox && xdotool key \"Return\"" 
	xdtFindFirefox := "xdotool search --onlyvisible --name firefox | head -n 1"
	xdtClickOnce := "xdotool key click 1 "
	xdtKeyDown := "xdotool key Down"
	xdtKeyTab := "xdotool key Tab"
	xdtKeyCtrlLAndType := "xdotool key \"ctrl+l\" type \"INSERTKEYPRESSES\" && xdotool key \"Return\"" 
	xdtFirefoxSaveAPage := "xdotool key \"ctrl+s\" sleep 2"
	
	xdtFirefoxGoToURL := "xdtool key \"ctrl+l\" && xdotool type \"INSERTKEYPRESSES\" && xdotool key \"Return\""
	xdtCloseFirefox := "xdotool key --clearmodifiers \"ctrl+F4\""
	xdtFindFirefoxAndOpenDevToolsAndAllowPasting := "xdotool search --onlyvisible --class \"firefox\" windowactivate --sync key --clearmodifiers \"ctrl+shift+k\" && xdotool type \"allow pasting\""
	xdtTypeSomething := "xdotool type \"INSERTKEYPRESSES\" && xdotool key \"Return\"" 
	
	initXdoTool := exec.Command(xdoOpenNewTermAndFirefox)
	err := initXdoTool.Run()
	checkError(err)

    savePageXdotool := exec.Command(xdtFirefoxOpenSavePage)
	err := savePageXdotool.Run()
	checkError(err)

	xdtTypeFullPath := strings.Replace(xdtTypeSomething, "INSERTKEYPRESSES", testDir, -1) 
	typeFullPathXdotool := exec.Command(xdtTypeFullPath)
	err := typeFullePathXdotool.Run()
	checkError(err)
	
	// Click 1 // to escape writing input for file name - we will return with CTRL+l to then retype any accidentally click input into the Name: input bar
// The problem with click 1 is where is it clicking need a granteed way to escape to GUI
	escapePathBarXdotool	:= exec.Command()
	err := escapePathBarXdotool.Run()
	checkError(err)

  // Tab 1 // To move gui to Dropdown on save All, HTML only, text
	  tabToDropDownXdotool := exec.Command(xdtKeyTab)
	  err := tabToDropDownXdotool.Run()
	  checkError(err)

      // Down at some point selecting Save only html
      dropDownToSaveOptionXdotool := exec.Command(xdtKeyDown)
	  err := tabToEscapeToGUIXdotool.Run()
	  checkError(err)

	xdtAndFFSaveProperly := false
	// check if _files or .txt
	// delete files

	if recursionCounter != 6 && xdtAndFFSaveProperly != true {
		recursionCounter, err = softConfFFToSaveAlwaysHTMLOnly(testDir, recursionCounter)
	}

	return recursionCounter, nil
}

func xdtFFGetNewPages(saveDirectory string, urlArr []string) error {
	var xdtTypeFullPath string
	xdtOpenNewTermAndFirefox := "xdotool key \"ctrl+alt+t\" sleep 1 type firefox && xdotool key \"Return\"" 
	xdtFindFirefox := "xdotool search --onlyvisible --name firefox | head -n 1"
	xdtClickOnce := "xdotool key click 1 "
	xdtKeyDown := "xdotool key Down"
	xdtKeyTab := "xdotool key Tab"
//	xdtKeyCtrlLAndType := "xdotool key \"ctrl+l\" type \"INSERTKEYPRESSES\" && xdotool key \"Return\"" 
	xdtFirefoxOpenSavePage := "xdotool key \"ctrl+s\" sleep 2"
	
	xdtFirefoxGoToURL := "xdtool key \"ctrl+l\" && xdotool type \"INSERTKEYPRESSES\" && xdotool key \"Return\""
	xdtCloseFirefox := "xdotool key --clearmodifiers \"ctrl+F4\""
	xdtFindFirefoxAndOpenDevToolsAndAllowPasting := "xdotool search --onlyvisible --class \"firefox\" windowactivate --sync key --clearmodifiers \"ctrl+shift+k\" && xdotool type \"allow pasting\""
	xdtTypeSomething := "xdotool type \"INSERTKEYPRESSES\" && xdotool key \"Return\"" 

	initXdoTool := exec.Command(xdtOpenNewTermAndFirefox)
	err := initXdoTool.Run()
	checkError(err)

	findFFXdoTool := exec.Command(xdtFindFirefox)
	err := findFFxdoTool.Run()
	checkError(err)


	for _, url := range urlArr { 
	  goToUrlXdotool := exec.Command(strings.Replace(xdtFirefoxGoToURL, "INSERTKEYPRESSES",url,-1))
	  err := goToUrlXdotool.Run()
	  checkError(err)

	  savePageXdotool := exec.Command(xdtFirefoxOpenSavePage)
	  err := savePageXdotool.Run()
	  checkError(err)

	  xdtTypeFullPath = strings.Replace(xdtTypeSomething, "INSERTKEYPRESSES", saveDirectory, -1) 
	  typeFullPathXdotool := exec.Command(xdtTypeFullPath)
	  err := typeFullePathXdotool.Run()
	  checkError(err)
	}

	closeFirefoxXdotool := exec.Command(xdtCloseFirefox)
	err := closeFirefoxXdotool.Run()
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

// thunkage for soup ++Brain
// Account for site differences
// Account for previously ran urls - they have to be unique so no worries

// What are we looking for in html? - ANS: urls, titles
// find artefacts containing Tokens 
// Store:
// id index_of_next_chunk  time + domain + urlused + artefacts { urlfound  + titles; artefacts + tokensMatched } terminationStub; 


Datastream := marshalArtefactstoDS()
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

	for _, url := range passedLinksAndTitleByDomainMap {
		page, err := curlNewArticle(url)
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


func refangUrl(inputUrl string) string {
	tmpUrl := strings.ReplaceAll(inputUrl, "hxxp", "http")
	tmpUrl = strings.ReplaceAll(tmpUrl, "[]://]", "://")
	outputUrl := strings.ReplaceAll(tmpUrl, "[.]", ".")
	return outputUrl
}
```
