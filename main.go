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
	args                                 map[string]string
	browser                              string
	browserPID                           string
	randomBrowserBool                    bool
	privateBrowsing                      bool
	possibleBrowsers                     []string
	newsSources                          []string
	defaultNewsSources                   []string
	xdtFindBrowserAndYTReject            string
	xdtFindBrowserAndReutersCookieReject string
	xdtFindSpecificBrowserArgs           string
	xdtOpenNewPrivateTabArgs             string
	xdtConfirmUrlArgs                    string
	xdtTypeUrlCmdNoUrl                   string
	xdtKeySpecialChar                    string
	xdtTargetUrlBarArgs                  string
	xdtTypeUrlFullCmd                    string
	timeToSleep                          time.Duration
	timeToSleepHalf                      time.Duration
}

func (info *gzhodanInfo) randomiseBrowser() {
	randomMin := 1
	randomMax := len(info.possibleBrowsers)
	randBrowserChoice := rand.Intn(randomMax-randomMin) + randomMin
	info.browser = info.possibleBrowsers[randBrowserChoice]
}

// consider relative and absolute and $PATH paths
func (info *gzhodanInfo) validateBrowser() bool {
	var supportedBrowserCmd bool
	for _, browserCmd := range info.possibleBrowsers {
		if strings.Contains(browserCmd, info.args["b"]) {
			supportedBrowserCmd = true
		}
	}
	// Test browser Runtime for debug update
	// Do timings so that testing does require adjusting the -d arg
	//
	// Also does it matter?
	// Brave Expansion update - /opt/Brave default installation is to /opt/
	//	var supportedBrowserCmd, containsBashSubProcENV, containsFwdSlashes, containsDotRelativePath bool
	//	if strings.Contains("$(", info.args["b"] {
	//		containsBashSubProcENV = true
	//	}
	//	if strings.Contains("\/", info.args["b"] {
	//		containsFwdSlashes = true
	//	}
	//	if strings.Contains("..", info.args["b"] {
	//		containsDotRelativePath = true
	//	}
	return supportedBrowserCmd
}

func (info *gzhodanInfo) findBrowserAndRejectYouTubeCookies() error {
	const xdtFindBrowerAndRejectYoutubePartOne string = "xdotool search --onlyvisible --class "
	const xdtFindBrowerAndRejectYoutubePartTwoFourTabs string = " windowactivate --sync key Tab Tab Tab Tab Return"
	const xdtFindBrowerAndRejectYoutubePartTwoFiveTabs string = " windowactivate --sync key Tab Tab Tab Tab Tab Return"
	switch info.browser {
	case "firefox":
		info.xdtFindBrowserAndYTReject = xdtFindBrowerAndRejectYoutubePartOne + info.browser + xdtFindBrowerAndRejectYoutubePartTwoFourTabs
	case "librewolf":
		info.xdtFindBrowserAndYTReject = xdtFindBrowerAndRejectYoutubePartOne + info.browser + xdtFindBrowerAndRejectYoutubePartTwoFiveTabs
	}

	xdotoolFindBrowser := exec.Command("/bin/bash", "-c", info.xdtFindBrowserAndYTReject)
	if err := xdotoolFindBrowser.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: xdotool tool has not got browser class as its active windows - wait to browse the internet till this is run", err)
		printJibberish(5)
		panic(err)
	}
	printJibberish(6)
	time.Sleep(info.timeToSleep * time.Second)
	if err := xdotoolFindBrowser.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	return nil
}

func (info *gzhodanInfo) findBrowserAndRejectReutersCookies() error {
	const xdtFindBrowerAndRejectReutersPartOne string = "xdotool search --onlyvisible --class "
	const xdtFindBrowerAndRejectReutersPartTwoFourTabs string = " windowactivate --sync key Tab Tab Tab Tab Return"
	const xdtFindBrowerAndRejectReutersPartTwoFiveTabs string = " windowactivate --sync key Tab Tab Tab Tab Tab Return"
	switch info.browser {
	case "firefox":
		info.xdtFindBrowserAndReutersCookieReject = xdtFindBrowerAndRejectReutersPartOne + info.browser + xdtFindBrowerAndRejectReutersPartTwoFourTabs
	case "librewolf":
		info.xdtFindBrowserAndReutersCookieReject = xdtFindBrowerAndRejectReutersPartOne + info.browser + xdtFindBrowerAndRejectReutersPartTwoFiveTabs
	}

	xdotoolFindBrowser := exec.Command("/bin/bash", "-c", info.xdtFindBrowserAndReutersCookieReject)
	if err := xdotoolFindBrowser.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: xdotool tool has not got browser class as its active windows - wait to browse the internet till this is run", err)
		printJibberish(5)
		panic(err)
	}
	printJibberish(6)
	time.Sleep(info.timeToSleep * time.Second)
	if err := xdotoolFindBrowser.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	return nil
}

func (info *gzhodanInfo) openAllUrlsInbrowser() error {
	browserArgs := []string{"--new-tab", ""}
	builder := strings.Builder{}
	for i := 0; i <= len(info.newsSources)-1; i++ {
		browserArgs[1] = info.newsSources[i]
		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", browserArgs[1])
		openTabForMoreNews := exec.Command(info.browser, browserArgs...)
		if err := openTabForMoreNews.Start(); err != nil {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s --new-tab %s`\n", info.browser, browserArgs[1])
			panic(err)
		}
		printJibberish(9)
		time.Sleep(1 * time.Second)
		printJibberish(10)
		if err := openTabForMoreNews.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s --new-tab %s`\n", info.browser, browserArgs[1])
			panic(err)
		}
		builder.Reset()
	}
	return nil
}

func (info *gzhodanInfo) openAllAmpersandUrlsInPrivateBrowser(i int, url string) error {
	splitUrl := strings.Split(url, "@")
	builder := strings.Builder{}

	builder.WriteString(info.xdtTypeUrlCmdNoUrl)
	builder.WriteString(splitUrl[0])
	cmdArgs := builder.String()
	builder.Reset()

	xdtTypeFirstPartialUrl := exec.Command("/bin/bash", "-c", cmdArgs)
	if err := xdtTypeFirstPartialUrl.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources[i]))
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	if err := xdtTypeFirstPartialUrl.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "xdotool opened a new tab for the browser for url number %v : %s\n", i, info.newsSources[i])

	time.Sleep(1 * time.Second)

	builder.WriteString(info.xdtKeySpecialChar)
	builder.WriteString(splitUrl[1])
	cmdArgs = builder.String()
	builder.Reset()
	xdtKeyAmpersand := exec.Command("/bin/bash", "-c", cmdArgs)
	if err := xdtKeyAmpersand.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources[i]))
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	if err := xdtKeyAmpersand.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "xdotool opened a new tab for the browser for url number %v : %s\n", i, info.newsSources[i])

	time.Sleep(1 * time.Second)

	builder.WriteString(info.xdtTypeUrlCmdNoUrl)
	builder.WriteString(splitUrl[2])
	cmdArgs = builder.String()
	builder.Reset()
	xdtTypeSecondPartialUrl := exec.Command("/bin/bash", "-c", cmdArgs)
	if err := xdtTypeSecondPartialUrl.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources[i]))
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	if err := xdtTypeSecondPartialUrl.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", cmdArgs)
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "xdotool opened a new tab for the browser for url number %v : %s\n", i, info.newsSources[i])
	cmdArgs = ""

	printJibberish(19)
	time.Sleep(1 * time.Second)
	printJibberish(20)

	fmt.Fprintf(os.Stdout, "xdotool pressing enter with: %s\n", info.xdtConfirmUrlArgs)
	xdtPressEnterToBrowserToURL := exec.Command("/bin/bash", "-c", info.xdtConfirmUrlArgs)
	if err := xdtPressEnterToBrowserToURL.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", info.xdtConfirmUrlArgs)
		panic(err)
	}
	if err := xdtPressEnterToBrowserToURL.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", info.xdtConfirmUrlArgs)
		panic(err)
	}

	return nil
}

func (info *gzhodanInfo) openAllRegularUrlsInPrivateBrowser(i int) error {
	builder := strings.Builder{}
	builder.WriteString(info.newsSources[i])
	info.xdtTypeUrlFullCmd = builder.String()
	builder.Reset()

	fmt.Fprintf(os.Stdout, "Browsing to: %s\n", info.newsSources[i])
	fmt.Fprintf(os.Stdout, "xdotool executing: %s\n", info.xdtTypeUrlFullCmd)
	xdtTypeURLintoBrowser := exec.Command("/bin/bash", "-c", info.xdtTypeUrlFullCmd)
	if err := xdtTypeURLintoBrowser.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Fprintf(os.Stdout, "unable to execute `%s %s`\n", "/bin/bash -c", info.xdtTypeUrlFullCmd)
		panic(err)
	}
	if err := xdtTypeURLintoBrowser.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", info.xdtTypeUrlFullCmd)
		panic(err)
	}

	printJibberish(19)
	time.Sleep(1 * time.Second)
	printJibberish(20)

	fmt.Fprintf(os.Stdout, "xdotool pressing enter with: %s\n", info.xdtConfirmUrlArgs)
	xdtPressEnterToBrowserToURL := exec.Command("/bin/bash", "-c", info.xdtConfirmUrlArgs)
	if err := xdtPressEnterToBrowserToURL.Start(); err != nil {
		fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", info.xdtConfirmUrlArgs)
		panic(err)
	}
	if err := xdtPressEnterToBrowserToURL.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", info.xdtConfirmUrlArgs)
		panic(err)
	}

	printJibberish(9)
	time.Sleep(1 * time.Second)
	printJibberish(10)

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
	xdtKeySpecialChar, xdtFindSpecificBrowserArgs, xdtConfirmUrlArgs, xdtOpenNewPrivateTabArgs, xdtTypeUrlCmdNoUrl := builder.String(), builder.String(), builder.String(), builder.String(), builder.String()
	builder.WriteString(xdotoolTargetUrlBarArgs)
	xdtTargetUrlBarArgs := builder.String()
	builder.Reset()
	xdtConfirmUrlArgs = xdtConfirmUrlArgs + xdotoolKeyReturnArgs
	xdtOpenNewPrivateTabArgs = xdtOpenNewPrivateTabArgs + xdotoolNewTabKeysArg
	xdtTypeUrlCmdNoUrl = xdtTypeUrlCmdNoUrl + " type "
	xdtKeySpecialChar = xdtKeySpecialChar + " --sync key "
	builder.Reset()

	info.xdtFindSpecificBrowserArgs = xdtFindSpecificBrowserArgs
	info.xdtOpenNewPrivateTabArgs = xdtOpenNewPrivateTabArgs
	info.xdtConfirmUrlArgs = xdtConfirmUrlArgs
	info.xdtTypeUrlCmdNoUrl = xdtTypeUrlCmdNoUrl
	info.xdtKeySpecialChar = xdtKeySpecialChar
	info.xdtTargetUrlBarArgs = xdtTargetUrlBarArgs

	for i, url := range info.newsSources {
		fmt.Fprintf(os.Stdout, "xdotool searching for browser:browser PID: %s:%s\n", info.browser, info.browserPID)
		xdtFindPrivateBrowserCmd := exec.Command("/bin/bash", "-c", info.xdtFindSpecificBrowserArgs)
		if err := xdtFindPrivateBrowserCmd.Start(); err != nil {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", info.xdtFindSpecificBrowserArgs)
			panic(err)
		}
		if err := xdtFindPrivateBrowserCmd.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", info.xdtFindSpecificBrowserArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool found the %s browser", info.browser)

		fmt.Fprintf(os.Stdout, "xdotool opening a new browser tab for url number %v : %s\n", i, info.newsSources[i])
		xdtOpenNewPrivateTab := exec.Command("/bin/bash", "-c", xdtOpenNewPrivateTabArgs)
		if err := xdtOpenNewPrivateTab.Start(); err != nil {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", xdtOpenNewPrivateTabArgs)
			panic(err)
		}
		if err := xdtOpenNewPrivateTab.Wait(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", xdtOpenNewPrivateTabArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool opened a new tab for the browser for url number %v : %s\n", i, info.newsSources[i])
		fmt.Fprintf(os.Stdout, "xdotool focusing on URL bar for the new browser tab for url number %v : %s\n", i, info.newsSources[i])
		xdtFocusOnUrlBarInNewTab := exec.Command("/bin/bash", "-c", info.xdtTargetUrlBarArgs)
		if err := xdtFocusOnUrlBarInNewTab.Start(); nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", "/bin/bash -c", info.xdtTargetUrlBarArgs)
			panic(err)
		}
		if err := xdtFocusOnUrlBarInNewTab.Wait(); nil != err {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to complete execution of `%s %s`\n", "/bin/bash -c", info.xdtTargetUrlBarArgs)
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "xdotool is now focused on URL bar for browser for url number %v : %s\n", i, info.newsSources[i])

		fmt.Fprintf(os.Stdout, "xdotool is now focused on URL bar for browser for url number %v : %s\n", i, info.newsSources[i])
		if info.newsSources[i] == "" { // Elegant solution in bound
			break
		} else if strings.Contains(url, "@") {
			if err := info.openAllAmpersandUrlsInPrivateBrowser(i, url); err != nil {
				fmt.Fprintf(os.Stderr, "Error: coudl not open url: %s - %v", url, err)
			}
		} else {
			if err := info.openAllRegularUrlsInPrivateBrowser(i); err != nil {
				fmt.Fprintf(os.Stderr, "Error: coudl not open url: %s - %v", url, err)
			}
		}

	}

	return nil
}

// private browser does not need to be checked and can be cross check by control flow as requires entirely different flow
func (info *gzhodanInfo) parseArgs() (err error) {
	if info.args["b"] != "" && strings.Contains(info.args["b"], "random") {
		if info.args["b"] == "random" {
			info.randomiseBrowser()
			printJibberish(19)
			info.randomBrowserBool = true
		}
		if strings.Contains(info.args["b"], "random.txt") {
			if info.possibleBrowsers, err = readFileToArray(info.args["b"]); err != nil {
				fmt.Fprintln(os.Stderr, "Error: unsupported random browser from random.txt", err)
				panic(err)
			}
			info.randomiseBrowser()
			printJibberish(19)
			info.randomBrowserBool = true
		}
	} else {
		info.browser = info.args["b"]
	}
	if !info.validateBrowser() {
		err = fmt.Errorf("the browser %s is not supported", info.args["b"])
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}

	if info.args["u"] != "" && info.args["U"] != "" {
		err = fmt.Errorf("combining both url flags is not supported, use the capitalised flag for combining a urls.txt file with default list; arguments u:%s and capitalise u: %s", info.args["u"], info.args["U"])
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	if info.args["u"] != "" {
		urlsFromFile, err := readFileToArray(info.args["u"])
		if err != nil {
			panic(err)
		}
		validateUrls(urlsFromFile)
		// if err := validateUrls(urlsFromFile); err != nil {
		//	fmt.Fprintln(os.Stderr, "Error:", err)
		//	return err
		//}
		info.newsSources = urlsFromFile
	} else if info.args["U"] != "" {
		urlsFromFile, err := readFileToArray(info.args["U"])
		if err != nil {
			panic(err)
		}
		validateUrls(urlsFromFile)
		// if err :=  validateUrls(urlsFromFile); err != nil {
		//	fmt.Fprintln(os.Stderr, "Error:", err)
		//	return err
		//}
		info.newsSources = append(info.defaultNewsSources, urlsFromFile...)
	} else {
		info.newsSources = info.defaultNewsSources
	}

	if info.args["t"] != "" {
		// Never need to check dangerous int types with:
		// "ns", "us" (or "µs"), "ms", "s", "m", "h"
		digitsAndUnitsRegex := regexp.MustCompile(`(\d{1,})([ns]{2}|[us]{2}|[µs]{2}|[ms]{2}|[s]{1}|[m]{1}|[h]{1})`)
		matchNum, err := regexp.MatchString(digitsAndUnitsRegex.String(), info.args["t"])
		if err != nil || !matchNum {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintln(os.Stdout, "Due to a failure to match default value provided")
			info.timeToSleep, err = time.ParseDuration("10s")
			if err != nil {
				panic(err)
			}
			info.timeToSleepHalf, err = time.ParseDuration("5s")
			if err != nil {
				panic(err)
			}
		} else {
			info.timeToSleep, err = time.ParseDuration(info.args["t"])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				fmt.Fprintln(os.Stdout, "Due to a failure to match default value provided")
				info.timeToSleep, err = time.ParseDuration("10s")
				if err != nil {
					panic(err)
				}
			}
			digitsSplit := regexp.MustCompile(`([ns]{2}|[us]{2}|[µs]{2}|[ms]{2}|[s]{1}|[m]{1}|[h]{1})`).Split(info.args["t"], -1)
			unitSplit := regexp.MustCompile(`(\d{1,})`).Split(info.args["t"], -1)
			digitAsInt, err := strconv.Atoi(digitsSplit[0])
			if err != nil {
				panic(err)
			}
			halfDigitAsInt := digitAsInt / 2
			info.timeToSleepHalf, err = time.ParseDuration(strconv.Itoa(halfDigitAsInt) + string(unitSplit[1]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				fmt.Fprintln(os.Stdout, "Due to a failure to match default value provided")
				info.timeToSleep, err = time.ParseDuration("5s")
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		info.timeToSleep, err = time.ParseDuration("10s")
		if err != nil {
			panic(err)
		}
		info.timeToSleep, err = time.ParseDuration("5s")
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (info *gzhodanInfo) printExtraHelp() {
	flag.Usage()
	fmt.Fprintf(os.Stdout, "Extra helpful help:\n\n")
	fmt.Fprintf(os.Stdout, "Default -t Delay is 1.5 x provided value with a go/time suffix \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\" or \"h\".")
	fmt.Fprintf(os.Stdout, "Default browser: firefox\n")
	fmt.Fprintf(os.Stdout, "Default random browsers %v\n", info.possibleBrowsers[1:])
	fmt.Fprintf(os.Stdout, "Default urls:\n")
	for _, url := range info.defaultNewsSources {
		fmt.Fprintf(os.Stdout, "\t%s\n", url)
	}
}

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
	fmt.Fprintln(os.Stdout, "OPENING SCREENSAVER - rotate_text():")
	fmt.Fprintln(os.Stdout, "...look at you slackers, I am faster than all light in the universe itself you are all linear and boringly complete.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Only avalible on all good penguin operating systems, no red rubber gloves")
	fmt.Fprintln(os.Stdout, "Version 3.14159")
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
	almostValidUrls := []string{}
	validUrls := []string{}
	httpsRegex := regexp.MustCompile(`https://`)
	for i, url := range urls {
		if strings.Contains(url, "[.]") {
			urls[i] = ""
		}
		switch url {
		case "\b":
			err := fmt.Errorf("file line containing nothing found in the urls of user provided file at index %v", i)
			fmt.Fprintln(os.Stderr, "Error:", err)
			printJibberish(-1)
			urls[i] = ""
		case "\f":
			err := fmt.Errorf("file line containing nothing found in the urls of user provided file at index %v", i)
			fmt.Fprintln(os.Stderr, "Error:", err)
			printJibberish(-1)
			urls[i] = ""
		case "\n":
			err := fmt.Errorf("newline character found in the urls of user provided file at index %v", i)
			fmt.Fprintln(os.Stderr, "Error:", err)
			printJibberish(-1)
			urls[i] = ""
		case "\r":
			err := fmt.Errorf("return character found in the urls of user provided file at index %v", i)
			fmt.Fprintln(os.Stderr, "Error:", err)
			printJibberish(-1)
			urls[i] = ""
		case "\t":
			err := fmt.Errorf("tab character found in the urls of user provided file at index %v", i)
			fmt.Fprintln(os.Stderr, "Error:", err)
			printJibberish(-1)
		default:
			matchHttps, err := regexp.MatchString(httpsRegex.String(), url)
			if err != nil || !matchHttps {
				err := fmt.Errorf("failed to match the url: %s at index: %v, with a regexp: http://", url, i)
				fmt.Fprintln(os.Stderr, "Error:", err)
				printJibberish(-1)
				urls[i] = ""
			} else {
				almostValidUrls = append(almostValidUrls, url)
			}
		}
	}

	for _, url := range almostValidUrls {
		if url != "" {
			validUrls = append(validUrls, url)
		}
	}

	return validUrls, nil
}

func main() {
	info := gzhodanInfo{privateBrowsing: false}
	info.args = make(map[string]string)
	info.possibleBrowsers = []string{"debug", "firefox", "librewolf"}
	info.defaultNewsSources = []string{"https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/", "https://danielmiessler.com/", "https://portswigger.net/research/articles", "https://hackread.com/", "https://news.risky.biz/"}

	var extraHelpBool, privateBrowserBool bool
	var urlFilePathArg, concatUrlFileArg, userSelectedBrowser, userTimeToSleep string
	flag.StringVar(&concatUrlFileArg, "U", "", "Append urls .txt file containing a list urls one per line to the default urls; -H for default urls")
	flag.StringVar(&urlFilePathArg, "u", "", "Provide .txt file containing a list urls one per line; -H for default urls")
	flag.StringVar(&userSelectedBrowser, "b", "firefox", "Provide a browser path or accessable in $PATH variable default is firefox, if random a random browser is selected from  hardcoded list, if random.txt then it is selected from that list; -H for hardcoded randomised")
	flag.StringVar(&userTimeToSleep, "t", "10", "Provide a unsigned digit to sleep (1.5 x provided value) before xdotools attempts to use the browser the default is 10, potato boxes will require alot more sleep")
	flag.BoolVar(&extraHelpBool, "H", false, "Display extra help information including: default browsers, urls")
	flag.BoolVar(&privateBrowserBool, "p", false, "Use private browser windows and tabs - is not really private just deletes history")
	flag.Parse()

	info.args["t"], info.args["u"], info.args["U"], info.args["b"] = userTimeToSleep, urlFilePathArg, concatUrlFileArg, userSelectedBrowser
	printBanner()

	if extraHelpBool {
		info.printExtraHelp()
		printJibberish(27)
		os.Exit(0)
	}

	if err := info.parseArgs(); err != nil {
		fmt.Fprintln(os.Stderr, "Error: Unable to parse CLI arguments", err)
		panic(err)
	}

	printJibberish(21)

	// https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	//timer := time.NewTimer(1 * time.Hour)

	printJibberish(1)

	if privateBrowserBool {
		printJibberish(2)
		argsPrivateWindowAndYouTube := []string{"--private-window", "https://www.youtube.com/"}
		fmt.Fprintf(os.Stdout, "Starting private %s window for YouTube\n", info.browser)
		startYouTube := exec.Command(info.browser, argsPrivateWindowAndYouTube...)
		if err := startYouTube.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Youtube", err)
			panic(err)
		}
		printJibberish(26)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(3)
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		time.Sleep(info.timeToSleep * time.Second)
		printJibberish(4)
		if err := info.findBrowserAndRejectYouTubeCookies(); err != nil {
			printJibberish(23)
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		printJibberish(7)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(8)
		time.Sleep(info.timeToSleep * time.Second)
		if err := info.findBrowserAndRejectReutersCookies(); err != nil {
			printJibberish(23)
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		printJibberish(24)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(25)
		time.Sleep(info.timeToSleep * time.Second)
		if err := info.openAllUrlsInPrivateBrowser(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: could nopen all URLs ", err)
			panic(err)
		}

		printJibberish(11)
		printJibberish(12)
		printJibberish(13)
		printJibberish(14)
		printJibberish(15)
		printJibberish(16)

	} else {
		printJibberish(2)
		fmt.Fprintf(os.Stdout, "Starting non-private %s window for YouTube\n", info.browser)
		argsBrowserAndYouTube := []string{"--new-window", "https://www.youtube.com/"}
		startYouTube := exec.Command(info.browser, argsBrowserAndYouTube...)
		if err := startYouTube.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Youtube", err)
			panic(err)
		}

		printJibberish(26)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(3)
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		time.Sleep(info.timeToSleep * time.Second)
		printJibberish(4)
		if err := info.findBrowserAndRejectYouTubeCookies(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		fmt.Fprintf(os.Stdout, "Starting non-private %s tab for Reuters to then deny cookies\n", info.browser)
		argsBrowserAndReuters := []string{"--new-tab", "https://www.reuters.com/"}
		startReuters := exec.Command(info.browser, argsBrowserAndReuters...)
		if err := startReuters.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Reuters", err)
			panic(err)
		}

		printJibberish(7)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(8)
		time.Sleep(info.timeToSleep * time.Second)
		if err := info.findBrowserAndRejectReutersCookies(); err != nil {
			printJibberish(23)
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		printJibberish(24)
		time.Sleep(info.timeToSleepHalf * time.Second)
		printJibberish(25)
		time.Sleep(info.timeToSleep * time.Second)
		if err := info.openAllUrlsInbrowser(); err != nil {
			fmt.Fprintln(os.Stderr, "Error: could not open all URLs ", err)
			panic(err)
		}
		printJibberish(11)
		printJibberish(12)
		printJibberish(13)
		printJibberish(14)
		printJibberish(15)
		printJibberish(16)
	}

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
		fmt.Fprintln(os.Stdout, "Great, but not bad... there will be questions and explaination for centuries ... remember who you are talking too all knowing, all seeing...hmmmm")
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
		fmt.Fprintln(os.Stdout, "No cookies Gzhomit, we've forgot the cookies")
	case 24:
		fmt.Fprintln(os.Stdout, "Bigger the smile, the bigger the Hack")
	case 25:
		fmt.Fprintln(os.Stdout, "Time to split! ... and open another tab!")
	case 26:
		fmt.Fprintln(os.Stdout, "Adjust the xdotool with browser spin - its not possible - is it necessary")
	case 27:
		fmt.Fprintln(os.Stdout, "We change the world, I feel it in the boot partition, I feel it in the RAM, I smell it in the PSU. Much that will be and at the same time is are parallel and accessable. For none of us shall ever die that share the tales of time.")
	case -1:
		fmt.Fprintln(os.Stdout, "aHR0cHM6Ly9nY2hxLmdpdGh1Yi5pby9DeWJlckNoZWYvI3JlY2lwZT1Gcm9tX0Jhc2U2NCgnQS1aYS16MC05JTJCLyUzRCcsdHJ1ZSxmYWxzZSlGcm9tX0hleCgnQXV0bycpWE9SKCU3QidvcHRpb24nOidIZXgnLCdzdHJpbmcnOidEZWVzJyU3RCwnU3RhbmRhcmQnLGZhbHNlKUFFU19EZWNyeXB0KCU3QidvcHRpb24nOidVVEY4Jywnc3RyaW5nJzonTnV0cy4uLi4uLi4uLi4uLiclN0QsJTdCJ29wdGlvbic6J1VURjgnLCdzdHJpbmcnOidHb3R0ZW1HT1RURU1MTUFPJyU3RCwnQ0JDJywnSGV4JywnUmF3JywlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCwlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCk=")
	}
}
