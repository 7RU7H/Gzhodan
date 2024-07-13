package main

import (
	"os/exec"
	"strings"
)

// https://github.com/ChrisPritchard/ctf-writeups/blob/master/GO-SCRIPTING.md

const (
	cybernews   string = "https://www.youtube.com/@cybernews/videos"
	seytonic    string = "https://www.youtube.com/@Seytonic/videos"
	hakfive     string = "https://www.youtube.com/@hak5/videos"
	firefoxCmd  string = "firefox"
	firefoxArgs string = " --new-window"
)

func main() {
	urlsBuilder := strings.Builder{}
	urlsBuilder.WriteString(firefoxArgs + " " + cybernews + " " + seytonic + " " + hakfive)
	allUrls := urlsBuilder.String()

	openAndCheckNewVideos := exec.Command(firefoxCmd, allUrls)
	err := openAndCheckNewVideos.Run()
	if nil != err {
		panic(err)
	}
}
