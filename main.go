package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type gzhodanInfo struct {
	browser     string
	browserPID  string
	newsSources []string
}

func (info *gzhodanInfo) preventProcrastination() error {
	killBrowserPID := exec.Command("kill", "-s SIGTERM", info.browserPID)
	err := killBrowserPID.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error: unable to kill the Browser PID", err)
		panic(err)

	}
	printJibberish(17)
	return nil
}

func (info *gzhodanInfo) randomiseBrowser(browsersArray []string) {
	randomMin := 1
	randomMax := len(browsersArray)
	randBrowserChoice := rand.Intn(randomMax-randomMin) + randomMin
	info.browser = browsersArray[randBrowserChoice]
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
			fmt.Fprintf(os.Stdout, "Unable complete execution of `%s --new-tab %s`\n", info.browser, browserArgs[1])
			panic(err)
		}
		builder.Reset()
	}
	return nil
}

func (info *gzhodanInfo) openAllUrlsInPrivateBrowser() error {
	const xdotoolBinary string = "xdotool"
	const xdotoolFindPrivateBrowserWindowArgs string = ""
	const xdotoolNewPrivateBrowserTabArgs string = ""
	const xdotoolTargetUrlBarArgs string = ""
	const xdotoolKeyReturn string = ""
	xdotoolTypeURLCmdAndArgs := ""
	builder := strings.Builder{}

	builder.WriteString(xdotoolFindPrivateBrowserWindowArgs)
	builder.Write()
	xdtFindSpecificBrowserArgs := builder.String()
	builder.Reset()
	for i := 0; i <= len(info.newsSources)-1; i++ {
		// grab browser window

		fmt.Fprintf(os.Stdout, "xdotool searching for browser:browser PID: %s:%s\n", info.browser, info.browserPID)
		xdtFindPrivateBrowserCmd := exec.Command(xdotoolBinary, xdtFindSpecificBrowserArgs)
		err := xdtFindPrivateBrowserCmd.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", xdotoolBinary, xdtFindSpecificBrowserArgs)
			panic(err)
		}
		err = xdtFindPrivateBrowserCmd.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable complete execution of `%s %s`\n")
			panic(err)
		}

		fmt.Fprintf(os.Stdout, "")
		// ctrl tab for new tab
		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", info.newsSources[i])
		xdtTypeURLintoBrowser := exec.Command(xdotoolBinary, xdotoolTypeUrlArgs)
		err = xdtTypeURLintoBrowser.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", xdotoolBinary, xdotoolTypeUrlArgs)
			panic(err)
		}
		err = xdtTypeURLintoBrowser.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable complete execution of `%s %s`\n", xdotoolBinary, xdotoolTypeUrlArgs)
			panic(err)
		}

		builder.WriteString(xdotoolTypeURLCmdAndArgs)
		builder.WriteString(info.newsSources[i])
		xdotoolTypeUrlArgs := builder.String()
		builder.Reset()

		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", info.newsSources[i])
		xdtTypeURLintoBrowser := exec.Command(xdotoolBinary, xdotoolTypeUrlArgs)
		err := xdtTypeURLintoBrowser.Start()
		if nil != err {
			fmt.Fprintf(os.Stdout, "Only %v urls accounted for...\n", i-len(info.newsSources))
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s %s`\n", xdotoolBinary, xdotoolTypeUrlArgs)
			panic(err)
		}
		err = xdtTypeURLintoBrowser.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable complete execution of `%s %s`\n", xdotoolBinary, xdotoolTypeUrlArgs)
			panic(err)
		}
		printJibberish(19)
		time.Sleep(1 * time.Second)
		printJibberish(20)
		// return key
		printJibberish(9)
		time.Sleep(1 * time.Second)
		printJibberish(10)

	}
	return nil
}

// Add a check for internet connectivity

// Kill process termination and graceful exit
// Does work, but does not exit after an 1 hour, but will Ctrl+C after an hour

// If private open private-window use a CTRL+T open new tab with xdotools
// Check brave browser hotkeys

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
	fmt.Fprintln(os.Stdout, "================================================================")
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
	fmt.Fprintln(os.Stdout, "Version 1.8532110112358")
	fmt.Fprintln(os.Stdout, "💀 Happy Hacking :) ... 💀")
}

func handleTermination(cancel context.CancelFunc) {
	printJibberish(18)
	fmt.Fprintln(os.Stdout, "Gzhodan> I am sorry, idiot I just can't do tha- \nGzhodan> ... ")
	cancel()
	os.Exit(0)
}

func main() {
	browsersArray := []string{"debug", "firefox", "librewolf"}

	info := gzhodanInfo{}
	privateBool := false // ISSUE regarding cli there is no --new-private-tab !!
	randomiseBrowserBool := false
	info.newsSources = []string{"https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/", "https://danielmiessler.com/", "https://portswigger.net/research/articles"}

	printBanner()

	if randomiseBrowserBool {
		printJibberish(19)
		info.randomiseBrowser(browsersArray)
	} else { //flag needs adding
		info.browser = "firefox"
	}

	argsPrivateWindowAndYouTubeCookies := []string{"--private-window", "https://www.youtube.com/"}
	argsAndYouTubeCookies := []string{"--new-window", "https://www.youtube.com/"}

	// https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	timer := time.NewTimer(1 * time.Hour)

	printJibberish(1)

	if privateBool {
		fmt.Fprintf(os.Stdout, "Starting private %s window for YouTube\n", info.browser)
		startYouTube := exec.Command(info.browser, argsPrivateWindowAndYouTubeCookies...)
		err := startYouTube.Start()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: browser could not open to Youtube", err)
			panic(err)
		}
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		err = info.findBrowserAndRejectYouTubeCookies()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}
		printJibberish(2)
		time.Sleep(5 * time.Second)
		printJibberish(3)
		time.Sleep(5 * time.Second)
		printJibberish(4)

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
		info.browserPID = strconv.Itoa(startYouTube.Process.Pid)
		err = info.findBrowserAndRejectYouTubeCookies()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error: could not reject YouTube cookies", err)
			panic(err)
		}

		printJibberish(2)
		time.Sleep(5 * time.Second)
		printJibberish(3)
		time.Sleep(5 * time.Second)
		printJibberish(4)

		printJibberish(7)
		fmt.Fprintf(os.Stdout, "Waiting 5 seconds - did you know that only the finest potatos are used in the upcoming release of Gzhados\n")
		time.Sleep(5 * time.Second)
		printJibberish(8)
		fmt.Fprintf(os.Stdout, "Done waiting 5 seconds, beginning to browser to YouTube Channels and News sites, AI Joe dedicated to Real Joe - LMAO\n")

		info.openAllUrlsInbrowser()

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

	<-timer.C
	defer info.preventProcrastination()
	os.Exit(0)
}

func printJibberish(jibberID int) {
	switch jibberID {
	case 0:
		fmt.Fprintf(os.Stdout, "f0936e3af2b30a378bef2d0549d722a50fe62543fc3c460690f902d40c3583b820b21034874d910e808ce637acf74387a43419c657b654f04e66e3b356a66eaa864704d504917bebc88297b9107f7b84d3bd6a13f02c5f4ed065ea0029f131aa27903a688166ea480c9fbbf8c2e36cc3d8c901bd277632417a9b2e8a48c3df33c070a10701c08407438366a119cdcae4ca7078407ee935445949cf1836bffd5375586c82e13f1ffeb40ee8a8125fdc2ddd05187c2f9ac8cd4dd6ee38e9d23b4b\n")
	case 1:
		fmt.Fprintf(os.Stdout, "Status!.. I don't know the CODES! Don't Give me excuses give me results! Navigation... ...\n")
	case 2:
		fmt.Fprintf(os.Stdout, "Waiting 10 seconds, because executing through golang process takes longer and xdotool needs the time..\n")
	case 3:
		fmt.Fprintf(os.Stdout, ".. You forgot to light the FUSE Gzhomit..don't look doooooooooooooooooown that alley!\n")
	case 4:
		fmt.Fprintf(os.Stdout, "Done waiting 10 seconds, *microwave* ping sound - no explosions ... aaaaah ...\n")
	case 5:
		fmt.Fprintf(os.Stdout, "You got the WRONG CLOUSERS..Gzhomit\n")
	case 6:
		fmt.Fprintf(os.Stdout, "Loop-da-looping the cables to 2000m nose drive drop - (Tabbing through YouTube cookies to reject them with xdotool)\n")
	case 7:
		fmt.Fprintf(os.Stdout, "Waiting 5 seconds - did you know that only the finest potatos are used in the upcoming release of Gzhados\n")
	case 8:
		fmt.Fprintf(os.Stdout, "Done waiting 5 seconds, beginning to browser to YouTube Channels and News sites, AI Joe dedicated to Real Joe - LMAO\n")
	case 9:
		fmt.Fprintf(os.Stdout, "Steady on single file ... Waiting 1 seconds\n")
	case 10:
		fmt.Fprintf(os.Stdout, "Done waiting 1 seconds, remember the answer follows the question, its dangerous if it goes the other way...\n")
	case 11:
		fmt.Fprintf(os.Stdout, "Steady on single file I said\n")
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
	case -1:
		fmt.Fprintln(os.Stdout, "aHR0cHM6Ly9nY2hxLmdpdGh1Yi5pby9DeWJlckNoZWYvI3JlY2lwZT1Gcm9tX0Jhc2U2NCgnQS1aYS16MC05JTJCLyUzRCcsdHJ1ZSxmYWxzZSlGcm9tX0hleCgnQXV0bycpWE9SKCU3QidvcHRpb24nOidIZXgnLCdzdHJpbmcnOidEZWVzJyU3RCwnU3RhbmRhcmQnLGZhbHNlKUFFU19EZWNyeXB0KCU3QidvcHRpb24nOidVVEY4Jywnc3RyaW5nJzonTnV0cy4uLi4uLi4uLi4uLiclN0QsJTdCJ29wdGlvbic6J1VURjgnLCdzdHJpbmcnOidHb3R0ZW1HT1RURU1MTUFPJyU3RCwnQ0JDJywnSGV4JywnUmF3JywlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCwlN0Inb3B0aW9uJzonSGV4Jywnc3RyaW5nJzonJyU3RCk=")
	}
}
