package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// https://github.com/ChrisPritchard/ctf-writeups/blob/master/GO-SCRIPTING.md

const (
	//terminalApp                    string = "xfce4-terminal"
	firefoxCmd                     string = "firefox"
	xdotoolCmd                     string = "xdotool"
	xdtFindFirefoxAndRejectYoutube string = "search --onlyvisible --class \"firefox\" windowactivate --sync && xdotool key Tab && xdotool key Tab && xdotool key Tab && xdotool key Tab && xdotool key Tab && xdotool key Return"
)

func printBanner() {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "  ▄████ ▒███████▒ ██░ ██ ▓█████▄  ▄▄▄       ███▄    █ ")
	fmt.Fprintln(os.Stdout, " ██▒ ▀█▒▒ ▒ ▒ ▄▀░▓██░ ██▒▒██▀ ██▌▒████▄     ██ ▀█   █ ")
	fmt.Fprintln(os.Stdout, "▒██░▄▄▄░░ ▒ ▄▀▒░ ▒██▀▀██░░██   █▌▒██  ▀█▄  ▓██  ▀█ ██▒")
	fmt.Fprintln(os.Stdout, "░▓█  ██▓  ▄▀▒   ░░▓█ ░██ ░▓█▄   ▌░██▄▄▄▄██ ▓██▒  ▐▌██▒")
	fmt.Fprintln(os.Stdout, "░▒▓███▀▒▒███████▒░▓█▒░██▓░▒████▓  ▓█   ▓██▒▒██░   ▓██░")
	fmt.Fprintln(os.Stdout, " ░▒   ▒ ░▒▒ ▓░▒░▒ ▒ ░░▒░▒ ▒▒▓  ▒  ▒▒   ▓▒█░░ ▒░   ▒ ▒ ")
	fmt.Fprintln(os.Stdout, "  ░   ░ ░░▒ ▒ ░ ▒ ▒ ░▒░ ░ ░ ▒  ▒   ▒   ▒▒ ░░ ░░   ░ ▒░")
	fmt.Fprintln(os.Stdout, "░ ░   ░ ░ ░ ░ ░ ░ ░  ░░ ░ ░ ░  ░   ░   ▒      ░   ░ ░ ")
	fmt.Fprintln(os.Stdout, "	   ░   ░ ░     ░  ░  ░   ░          ░  ░         ░ ")
	fmt.Fprintln(os.Stdout, "	   ░                 ░                             ")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Gzhodan - Goodbye AGI, APTs and Aliens")
	fmt.Fprintln(os.Stdout, "Astatical GPU Idols")
	fmt.Fprintln(os.Stdout, "Abused Party of Tools")
	fmt.Fprintln(os.Stdout, "Aliens: weird rapey people pretending to be many powers of the (theoretical) mathetical defintion of 'cool' than they actually are (No Aliens in 150 Million Lightyears btw)")
	fmt.Fprintln(os.Stdout, "...look at you slackers, I am faster than all light in the universe itself you are all linear and boringly complete.")
	fmt.Fprintln(os.Stdout, "💀 Happy Hacking :) ... 💀")
}

func main() {
	printBanner()

	firefox := exec.Command(firefoxCmd)
	err := firefox.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	err = firefox.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	time.Sleep(1 * time.Second)

	argsAndYouTubeCookies := []string{"--new-tab", "https://www.youtube.com/"}
	startYouTube := exec.Command(firefoxCmd, argsAndYouTubeCookies...)
	err = startYouTube.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	err = startYouTube.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}
	time.Sleep(5 * time.Second)

	xdotoolFindFF := exec.Command(xdotoolCmd, xdtFindFirefoxAndRejectYoutube)
	err = xdotoolFindFF.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error:", err)
		//panic(err)
	}
	err = xdotoolFindFF.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		//panic(err)
	}
	time.Sleep(1 * time.Second)

	argsAndYouTubeUrls := []string{"--new-tab", "https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos"}
	getAllNewsUrls := exec.Command(firefoxCmd, argsAndYouTubeUrls...)
	err = getAllNewsUrls.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = getAllNewsUrls.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}

	time.Sleep(1 * time.Second)

}
