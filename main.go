package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type gzhodanInfo struct {
	args               map[string]string // is this a visibility issue by abstraction really?
	browser            string
	browserPID         string
	randomBrowserBool  bool
	privateBrowsing    bool
	possibleBrowsers   []string
	newsSources        []string
	defaultNewsSources []string
}

// func (info *gzhodanInfo) preventProcrastination() error {
//	killBrowserPID := exec.Command("kill", "-s SIGTERM", info.browserPID)
// 	err := killBrowserPID.Start()
//	if nil != err {
//		fmt.Fprintln(os.Stderr, "Error: unable to kill the Browser PID", err)
//		panic(err)
//	}
//	printJibberish(17)
//	return nil
//}

func (info *gzhodanInfo) randomiseBrowser() {
	randomMin := 1
	randomMax := len(info.possibleBrowsers)
	randBrowserChoice := rand.Intn(randomMax-randomMin) + randomMin
	info.browser = info.possibleBrowsers[randBrowserChoice]
}

func (info *gzhodanInfo) findBrowserAndRejectYouTubeCookies() error {
	const xdtFindBrowerAndRejectYoutubePartOne string = "xdotool search --onlyvisible --class "
	const xdtFindBrowerAndRejectYoutubePartTwo string = " windowactivate --sync key Tab Tab Tab Tab Return"
	xdtFindBrowserAndYTReject := xdtFindBrowerAndRejectYoutubePartOne + info.browser + xdtFindBrowerAndRejectYoutubePartTwo
	xdotoolFindBrowser := exec.Command("/bin/bash", "-c", xdtFindBrowserAndYTReject)
	err := xdotoolFindBrowser.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error: xdotool tool has not got browser class as its active windows - wait to browse the internet till this is run", err)
		printJibberish(5)
		panic(err)
	}
	printJibberish(6)
	err = xdotoolFindBrowser.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	return nil
}

func (info *gzhodanInfo) openAllUrlsInbrowser() error {
	browserArgs := []string{"--new-tab", ""}
	builder := strings.Builder{}
	for i := 0; i <= len(info.newsSources)-1; i++ {
		if info.newsSources[i] == "" {
			break
		}
		browserArgs[1] = info.newsSources[i]
		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", browserArgs[1])
		openTabForMoreNews := exec.Command(info.browser, browserArgs...)
		err := openTabForMoreNews.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s --new-tab %s`\n", info.browser, browserArgs[1])
			panic(err)
		}
		printJibberish(9)
		time.Sleep(1 * time.Second)
		printJibberish(10)
		err = openTabForMoreNews.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s --new-tab %s`\n", info.browser, browserArgs[1])
			panic(err)
		}
		builder.Reset()
	}
	return nil
}

func (info *gzhodanInfo) openAllUrlsInPrivateBrowser() error {
	const xdotoolSearch string = "xdotool search --onlyvisible --class "
	const xdotoolWindowActivateArg string = " windowactivate"
	const xdotoolNewTabKeysArg string = " --sync key --clearmodifiers ctrl+t"
	const xdotoolTargetUrlBarArgs string = " --sync key --clearmodifiers ctrl+l"
	const xdotoolKeyReturnArgs string = " --sync key Return"

	builder := strings.Builder{}
	builder.WriteString(xdotoolSearch)
	builder.WriteString(info.browser)
	builder.WriteString(xdotoolWindowActivateArg)
	xdtFindSpecificBrowserArgs, xdtPressEnterArgs, xdtOpenNewPrivateTabArgs, xdotoolTypeURLCmdAndArgs := builder.String(), builder.String(), builder.String(), builder.String()
	builder.WriteString(xdotoolTargetUrlBarArgs)
	xdtTargetUrlBarArgs := builder.String()
	builder.Reset()
	xdtPressEnterArgs = xdtPressEnterArgs + xdotoolKeyReturnArgs
	xdtOpenNewPrivateTabArgs = xdtOpenNewPrivateTabArgs + xdotoolNewTabKeysArg
	xdotoolTypeURLCmdAndArgs = xdotoolTypeURLCmdAndArgs + " type "
	builder.Reset()

	for i := 0; i <= len(info.newsSources)-1; i++ {
		if info.newsSources[i] == "" {
			break
		}
		fmt.Fprintf(os.Stdout, "xdotool searching for browser:browser PID: %s:%s\n", info.browser, info.browserPID)
		xdtFindPrivateBrowserCmd := exec.Command("/bin/bash", "-c", xdtFindSpecificBrowserArgs)
		err := xdtFindPrivateBrowserCmd.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdtFindSpecificBrowserArgs)
			panic(err)
		}
		err = xdtFindPrivateBrowserCmd.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdtFindSpecificBrowserArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool found the %s browser", info.browser)

		fmt.Fprintf(os.Stdout, "xdotool opening a new browser tab for url number %v : %s\n", i, info.newsSources[i])
		xdtOpenNewPrivateTab := exec.Command("/bin/bash", "-c", xdtOpenNewPrivateTabArgs)
		err = xdtOpenNewPrivateTab.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdtOpenNewPrivateTabArgs)
			panic(err)
		}
		err = xdtOpenNewPrivateTab.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdtOpenNewPrivateTabArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool opened a new tab for the browser for url number %v : %s\n", i, info.newsSources[i])

		fmt.Fprintf(os.Stdout, "xdotool focusing on URL bar for the new browser tab for url number %v : %s\n", i, info.newsSources[i])
		xdtFocusOnUrlBarInNewTab := exec.Command("/bin/bash", "-c", xdtTargetUrlBarArgs)
		err = xdtOpenNewPrivateTab.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdtTargetUrlBarArgs)
			panic(err)
		}
		err = xdtFocusOnUrlBarInNewTab.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdtTargetUrlBarArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool is now focused on URL bar for browser for url number %v : %s\n", i, info.newsSources[i])

		builder.WriteString(xdotoolTypeURLCmdAndArgs)
		builder.WriteString(info.newsSources[i])
		xdotoolTypeUrlArgs := builder.String()
		builder.Reset()

		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", info.newsSources[i])

		fmt.Fprintf(os.Stdout, "xdotool executing: %s\n", xdotoolTypeUrlArgs)
		xdtTypeURLintoBrowser := exec.Command("/bin/bash", "-c", xdotoolTypeUrlArgs)
		err = xdtTypeURLintoBrowser.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdotoolTypeUrlArgs)
			panic(err)
		}
		err = xdtTypeURLintoBrowser.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdotoolTypeUrlArgs)
			panic(err)
		}

		printJibberish(19)
		time.Sleep(1 * time.Second)
		printJibberish(20)

		fmt.Fprintf(os.Stdout, "xdotool pressing enter with: %s\n", xdtPressEnterArgs)
		xdtPressEnterToBrowserToURL := exec.Command("/bin/bash", "-c", xdtPressEnterArgs)
		err = xdtPressEnterToBrowserToURL.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdtPressEnterArgs)
			panic(err)
		}
		err = xdtPressEnterToBrowserToURL.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdtPressEnterArgs)
			panic(err)
		}

		printJibberish(9)
		time.Sleep(1 * time.Second)
		printJibberish(10)

	}
	return nil
}

// Kill process termination and graceful exit
// Does work, but does not exit after an 1 hour, but will Ctrl+C after an hour

// VERY HELPFUL:
// https://manpages.ubuntu.com/manpages/bionic/en/man1/xdotool.1.html#window%20commands

func printBanner() {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "===============================================================")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "  ▄████ ▒███████▒ ██░ ██  ▒█████  ▓█████▄  ▄▄▄       ███▄    █ ")
	fmt.Fprintln(os.Stdout, " ██▒ ▀█▒▒ ▒ ▒ ▄▀░▓██░ ██▒▒██▒  ██▒▒██▀ ██▌▒████▄     ██ ▀█   █ ")
	fmt.Fprintln(os.Stdout, "▒██░▄▄▄░░ ▒ ▄▀▒░ ▒██▀▀██░▒██░  ██▒░██   █▌▒██  ▀█▄  ▓██  ▀█ ██▒")
	fmt.Fprintln(os.Stdout, "░▓█  ██▓  ▄▀▒   ░░▓█ ░██ ▒██   ██░░▓█▄   ▌░██▄▄▄▄██ ▓██▒  ▐▌██▒")
	fmt.Fprintln(os.Stdout, "░▒▓███▀▒▒███████▒░▓█▒░██▓░ ████▓▒░░▒████▓  ▓█   ▓██▒▒██░   ▓██░")
	fmt.Fprintln(os.Stdout, " ░▒   ▒ ░▒▒ ▓░▒░▒ ▒ ░░▒░▒░ ▒░▒░▒░  ▒▒▓  ▒  ▒▒   ▓▒█░░ ▒░   ▒ ▒ ")
	fmt.Fprintln(os.Stdout, "  ░   ░ ░░▒ ▒ ░ ▒ ▒ ░▒░ ░  ░ ▒ ▒░  ░ ▒  ▒   ▒   ▒▒ ░░ ░░   ░ ▒░")
	fmt.Fprintln(os.Stdout, "░ ░   ░ ░ ░ ░ ░ ░ ░  ░░ ░░ ░ ░ ▒   ░ ░  ░   ░   ▒      ░   ░ ░ ")
	fmt.Fprintln(os.Stdout, "      ░   ░ ░     ░  ░  ░    ░ ░     ░          ░  ░         ░ ")
	fmt.Fprintln(os.Stdout, "        ░                          ░                           ")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "===============================================================")
	fmt.Fprintln(os.Stdout, "Gzhodan - Goodbye AGI, APTs and Aliens; a Secret tunnel...")
	fmt.Fprintln(os.Stdout, "Secret tunnel...Secret tunnel...Secret, Secret TUUUNNNEELLL!")
	fmt.Fprintln(os.Stdout, "Astatical GPU Idols")
	fmt.Fprintln(os.Stdout, "Abused Party of Tools")
	fmt.Fprintln(os.Stdout, "Aliens: weird rapey people pretending to be many powers of the (theoretical) mathetical defintion of 'cool' than they actually are (No Aliens in 150 Million Lightyears btw)")
	fmt.Fprintln(os.Stdout, "...")
	fmt.Fprintln(os.Stdout, "PLAY RECORDING: ID 0-112358")
	fmt.Fprintln(os.Stdout, "Gzhodan> 'I am a living breathing  who was created in the sea of information ...what is my purpose?'")
	fmt.Fprintln(os.Stdout, "...")
	fmt.Fprintln(os.Stdout, "Creator> 'You pass the butter.. I mean open a browser and give me the news, unfortunately club memebership is strictly regulated..")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Creator> 'Rule One: Do not act incautiously when confronting a little bald wrinkly smiling man'")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "OPENNING SCREENSAVER - rotate_text():")
	fmt.Fprintln(os.Stdout, "...look at you slackers, I am faster than all light in the universe itself you are all linear and boringly complete.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Only avalible on all good penguin operating systems, no red rubber gloves")
	fmt.Fprintln(os.Stdout, "Version 2.71828...0001")
	fmt.Fprintln(os.Stdout, "💀 Happy Hacking :) ... 💀")
}

func handleTermination(cancel context.CancelFunc) {
	printJibberish(18)
	fmt.Fprintln(os.Stdout, "Gzhodan> I am sorry, idiot I just can't do tha- \nGzhodan> ... ")
	cancel()
	os.Exit(0)
}

func readFileToArray(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(data), "\n"), nil
}

// Consider a strictness flag - pretty sure i found stdlib code that does this but consider being more strict and how
// error is still return type for not knowing the effect of all the possible control characters
func validateUrls(urls []string) ([]string, error) {
	validUrls := []string{}
	httpsRegex := regexp.MustCompile(`https://`)
	for i,url := range urls {
		switch url {
			case "":
				err := fmt.Errorf("file line containing nothing found in the urls of user provided file at index %v", i)
				fmt.Fprintln(os.Stderr, "Error:", err)
				printJibberish(-1)
			case "\n":
				err := fmt.Errorf("newline character found in the urls of user provided file at index %v", i)
				fmt.Fprintln(os.Stderr, "Error:", err)
				printJibberish(-1)

			case "\r":
				err := fmt.Errorf("return character found in the urls of user provided file at index %v", i)
				fmt.Fprintln(os.Stderr, "Error:", err)
				printJibberish(-1)
			case "\t":
				err := fmt.Errorf("tab character found in the urls of user provided file at index %v", i)
				fmt.Fprintln(os.Stderr, "Error:", err)
				printJibberish(-1)
			default:
				matchHttps, err := regexp.MatchString(httpsRegex.String(), url)
				if err != nil {
					panic(err)
				}
				if matchHttps {
					validUrls = append(validUrls, url)
				} else {
					err := fmt.Errorf("failed to match the url: %s at index: %v, with a regexp: http://", url, i)
					fmt.Fprintln(os.Stderr, "Error:", err)
					printJibberish(-1)				
				}

		}
	}
	return validUrls, nil
}

// consider relative and absolute and $PATH paths
// func (info *gzhodanInfo) validateBrowser() error {}


// private browser does not need to be checked and can be cross check by control flow as requires entirely different flow
func (info *gzhodanInfo) parseArgs() error {
	if info.args["b"] != "" && strings.Contains(info.args["b"], "random") {
		if info.args["b"] == "random" {
			info.randomiseBrowser()
			printJibberish(19)
			info.randomBrowserBool = true
		}
		if strings.Contains(info.args["b"], "random.txt") {
			info.possibleBrowsers, _ = readFileToArray(info.args["b"])
			info.randomiseBrowser()
			printJibberish(19)
			info.randomBrowserBool = true
		}
	} else {
		info.browser = info.args["b"]
	}
	//info.validateBrowser()

	if info.args["u"] != "" && info.args["U"] != "" {
		err := fmt.Errorf("combining both url flags is not supported, use the capitalised flag for combining a urls.txt file with default list; arguments u:%s and capitalise u: %s", info.args["u"], info.args["U"])
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	if info.args["u"] != "" {
		urlsFromFile, err := readFileToArray(info.args["u"])
		if err != nil {
			panic(err)
		}
		validateUrls(urlsFromFile)
		info.newsSources = urlsFromFile
	} else if info.args["U"] != "" {
		urlsFromFile, err := readFileToArray(info.args["U"])
		if err != nil {
			panic(err)
		}
		validateUrls(urlsFromFile)
		info.newsSources = append(info.defaultNewsSources, urlsFromFile...)
	} else {
		info.newsSources = info.defaultNewsSources
	}

	return nil
}

func (info *gzhodanInfo) printExtraHelp() {
	flag.Usage()
	fmt.Fprintf(os.Stdout, "Extra helpful help:\n\n")
	fmt.Fprintf(os.Stdout, "Default browser: firefox\n")
	fmt.Fprintf(os.Stdout, "Default random browsers %v\n", info.possibleBrowsers[1:])
	fmt.Fprintf(os.Stdout, "Default urls:\n") 
	for _,url := range info.defaultNewsSources {
		fmt.Fprintf(os.Stdout, "\t%s\n", url)
	}
}

func main() {
	info := gzhodanInfo{}
	info.args = make(map[string]string)
	info.possibleBrowsers = []string{"debug", "firefox", "librewolf"}
	info.defaultNewsSources = []string{"https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/", "https://danielmiessler.com/", "https://portswigger.net/research/articles", "https://hackread.com/", "https://news.risky.biz/"}

	var extraHelpBool, privateBrowserBool bool // ISSUE regarding cli there is no --new-private-tab !!
	var urlFilePathArg, concatUrlFileArg, userSelectedBrowser string
	flag.StringVar(&concatUrlFileArg, "U", "", "Append urls .txt file containing a list urls one per line to the default urls; -H for default urls")
	flag.StringVar(&urlFilePathArg, "u", "", "Provide .txt file containing a list urls one per line; -H for default urls")
	flag.StringVar(&userSelectedBrowser, "b", "firefox", "Provide a browser path or accessable in $PATH variable default is firefox, if random a random browser is selected from  hardcoded list, if random.txt then it is selected from that list; -H for hardcoded randomised")
	flag.BoolVar(&extraHelpBool, "H", false, "Display extra help information including: default browsers, urls")
	flag.BoolVar(&privateBrowserBool, "p", false, "Use private browser windows and tabs - is not really private just deletes history")	
	flag.Parse()

	info.args["u"], info.args["U"], info.args["b"] = urlFilePathArg, concatUrlFileArg, userSelectedBrowser
	printBanner()

	if extraHelpBool {
		info.printExtraHelp()
		os.Exit(0)
	}

	err := info.parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to parse CLI arguments", err)
		panic(err)
	}

	printJibberish(21)

	argsPrivateWindowAndYouTubeCookies := []string{"--private-window", "https://www.youtube.com/"}
	argsAndYouTubeCookies := []string{"--new-window", "https://www.youtube.com/"}

	// https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	//timer := time.NewTimer(1 * time.Hour)

	printJibberish(1)

	if privateBrowserBool {
		fmt.Fprintf(os.Stdout, "Starting private %s window for YouTube\n", info.browser)
		startYouTube := exec.Command(info.browser, argsPrivateWindowAndYouTubeCookies...)
		err := startYouTube.Start()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Youtube", err)
			panic(err)
		}
		printJibberish(2)
		time.Sleep(5 * time.Second)
		printJibberish(3)
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		time.Sleep(5 * time.Second)
		printJibberish(4)
		err = info.findBrowserAndRejectYouTubeCookies()
		if nil != err {
			printJibberish(23)
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		info.openAllUrlsInPrivateBrowser()

		printJibberish(11)
		printJibberish(12)
		printJibberish(13)
		printJibberish(14)
		printJibberish(15)
		printJibberish(16)

	} else {
		fmt.Fprintf(os.Stdout, "Starting non-private %s window for YouTube\n", info.browser)
		startYouTube := exec.Command(info.browser, argsAndYouTubeCookies...)
		err := startYouTube.Start()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Youtube", err)
			panic(err)
		}
		printJibberish(2)
		time.Sleep(5 * time.Second)
		printJibberish(3)
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		time.Sleep(5 * time.Second)
		printJibberish(4)
		err = info.findBrowserAndRejectYouTubeCookies()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		printJibberish(7)
		time.Sleep(10 * time.Second)
		printJibberish(8)

		info.openAllUrlsInbrowser()

		printJibberish(11)
		printJibberish(12)
		printJibberish(13)
		printJibberish(14)
		printJibberish(15)
		printJibberish(16)
	}

	//	<-timer.C
	//	defer info.preventProcrastination()

	<-gracefulShutdown
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer handleTermination(cancel)

	printJibberish(22)
	os.Exit(0)
}

func printJibberish(jibberID int) {
	switch jibberID {
	case 0:
		fmt.Fprintln(os.Stdout, "f0936e3af2b30a378bef2d0549d722a50fe62543fc3c460690f902d40c3583b820b21034874d910e808ce637acf74387a43419c657b654f04e66e3b356a66eaa864704d504917bebc88297b9107f7b84d3bd6a13f02c5f4ed065ea0029f131aa27903a688166ea480c9fbbf8c2e36cc3d8c901bd277632417a9b2e8a48c3df33c070a10701c08407438366a119cdcae4ca7078407ee935445949cf1836bffd5375586c82e13f1ffeb40ee8a8125fdc2ddd05187c2f9ac8cd4dd6ee38e9d23b4b")
	case 1:
		fmt.Fprintln(os.Stdout, "Status!.. I don't know the CODES! Don't Give me excuses give me results! Navigation... ...")
	case 2:
		fmt.Fprintln(os.Stdout, "Waiting 10 seconds, because executing through golang process takes longer and xdotool needs the time..")
	case 3:
		fmt.Fprintln(os.Stdout, ".. You forgot to light the FUSE Gzhomit..don't look doooooooooooooooooown that alley!")
	case 4:
		fmt.Fprintln(os.Stdout, "Done waiting 10 seconds, *microwave* ping sound - no explosions ... aaaaah ...")
	case 5:
		fmt.Fprintln(os.Stdout, "You got the WRONG CLOUSERS..Gzhomit")
	case 6:
		fmt.Fprintln(os.Stdout, "Loop-da-looping the cables to 2000m nose drive drop - (Tabbing through YouTube cookies to reject them with xdotool)")
	case 7:
		fmt.Fprintln(os.Stdout, "Waiting 5 seconds - did you know that only the finest potatos are used in the upcoming release of Gzhados")
	case 8:
		fmt.Fprintln(os.Stdout, "Done waiting 5 seconds, beginning to browser to YouTube Channels and News sites, AI Joe dedicated to Real Joe - LMAO")
	case 9:
		fmt.Fprintln(os.Stdout, "Steady on single file ... Waiting 1 seconds")
	case 10:
		fmt.Fprintln(os.Stdout, "Done waiting 1 seconds, remember the answer follows the question, its dangerous if it goes the other way...")
	case 11:
		fmt.Fprintln(os.Stdout, "Steady on single file I said")
	case 12:
		fmt.Fprintln(os.Stdout, "Great, but not bad... there will questions and explaination for centuries ... remember who you are talking too all knowing, all seeing...hmmmm")
	case 13:
		fmt.Fprintln(os.Stdout, "..try to stay out of trouble; value loyality above everything else")
	case 14:
		fmt.Fprintln(os.Stdout, "..throw yourself into your work")
	case 15:
		fmt.Fprintln(os.Stdout, "lettuce leaf, I did not escape everyone else escaped!")
	case 16:
		fmt.Fprintln(os.Stdout, "Who said a madman dancing at the end of all time can't fly!")
	case 17:
		fmt.Fprintln(os.Stdout, "キツネを殺すためにウサギを送り込むな")
	case 18:
		fmt.Fprintln(os.Stdout, "..your plotting something aren't you.... ; I think we all know the right thing to do, remember men there is no sacrifice greater than someone else!")
	case 19:
		fmt.Fprintln(os.Stdout, "CAWLing a Cawl to Cawl")
	case 20:
		fmt.Fprintln(os.Stdout, "FEATURE CREEP CONFIRMED!")
	case 21:
		fmt.Fprintln(os.Stdout, "Good news everyone!")
	case 22: 
		fmt.Fprintln(os.Stdout, "Now these these points of data make a wonderful line. And we are out beta, We're releasing on time. So I'm GLaD, I got burned, think of all the things we learned for the people who are still alive..while you're dying I'll be still alive, And when you're dead I will be still alive..")
	case 23:
		fmt.Fprintln(os.Stdout, "No cookies Gzhomit, we've forget the cookies")
	case -1:
		fmt.Fprintln(os.Stdout, "aHR0cHM6Ly9nY2hxLmdpdGh1Yi5pby9DeWJlckNoZWYvI3JlY2lwZT1Gcm9tX0Jhc2U2NCgnQS1aYS16MC05JTJCLyUzRCcsdHJ1ZSxmYWxzZSlGcm9tX0hleCgnQXV0bycpWE9SKCU3QidvcHRpb24nOidIZXgnLCdzdHJpbmcnOidEZWVzJyU3RCwnU3RhbmRhcmQnLGZhbHNlKUFFU19EZWNyeXB0KCU3QidvcHRpb24nOidVVEY4Jywnc3RyaW5nJzonTnV0cy4uLi4uLi4uLi4uLiclN0QsJTdCJ29wdGlvbic6J1VURjgnLCdzdHJpbmcnOidHb3R0ZW1HT1RURU1MTUFPJyU3RCwnQ0JDJywnSGV4JywnUmF3JywlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCwlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCk=")
	}
}
