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
```