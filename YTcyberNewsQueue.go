package main

import (
	"os/exec"
)

// https://github.com/ChrisPritchard/ctf-writeups/blob/master/GO-SCRIPTING.md

const (
	firefoxCmd string = "firefox"
)

func main() {
	argsAndUrls := []string{"--new-window", "https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos"}
	openAndCheckNewVideos := exec.Command(firefoxCmd, argsAndUrls...)
	err := openAndCheckNewVideos.Start()
	if nil != err {
		panic(err)
	}
	err = openAndCheckNewVideos.Wait()
	if err != nil {
		panic(err)
	}
}
