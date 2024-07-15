package main

import (
	"os/exec"
	"time"
)

// https://github.com/ChrisPritchard/ctf-writeups/blob/master/GO-SCRIPTING.md

const (
	terminalApp    string = "xfce4-terminal"
	firefoxCmd     string = "firefox"
	xdotoolCmd     string = "xdotool"
	xdtFindFirefox string = "search --onlyvisible --class \"firefox\" windowactivate --sync key --clearmodifiers \"ctrl+shift+k\""
)

func main() {
	firefox := exec.Command(firefoxCmd)
	err := firefox.Start()
	if nil != err {
		panic(err)
	}
	err = firefox.Wait()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	argsAndYouTubeCookies := []string{"--new-tab", "https://www.youtube.com/"}
	startYouTube := exec.Command(firefoxCmd, argsAndYouTubeCookies...)
	err = startYouTube.Start()
	if nil != err {
		panic(err)
	}
	err = startYouTube.Wait()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)

	xdotoolFindFF := exec.Command(xdotoolCmd, xdtFindFirefox)
	err = xdotoolFindFF.Start()
	if nil != err {
		panic(err)
	}
	err = xdotoolFindFF.Wait()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)

	// Reject Cookies

	argsAndYouTubeUrls := []string{"--new-tab", "https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos"}
	openAndCheckNewVideos := exec.Command(firefoxCmd, argsAndYouTubeUrls...)
	err = openAndCheckNewVideos.Start()
	if nil != err {
		panic(err)
	}
	err = openAndCheckNewVideos.Wait()
	if err != nil {
		panic(err)
	}

	// Prompt
	// "Have you watched the newest Videos and noted information required for objectives"
	// I want cool spinning loading icon

	argsAndWrittenNewsUrls := []string{"--new-tab", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/"}
	getWrittenNewsUrls := exec.Command(firefoxCmd, argsAndWrittenNewsUrls...)
	err = getWrittenNewsUrls.Start()
	if nil != err {
		panic(err)
	}
	err = getWrittenNewsUrls.Wait()
	if err != nil {
		panic(err)
	}
}
