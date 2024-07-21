package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
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

func handleTermination(cancel context.CancelFunc) {
	fmt.Fprintln(os.Stdout, "Gzhodan> I am sorry, idiot I just can't do tha- \nGzhodan> ... ")
	cancel()
}

func main() {
	printBanner()
	// https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	argsAndYouTubeCookies := []string{"--new-window", "https://www.youtube.com/"}
	startYouTube := exec.Command(firefoxCmd, argsAndYouTubeCookies...)
	err := startYouTube.Start()
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
	time.Sleep(5 * time.Second)

	newsSources := []string{"https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/", "https://portswigger.net/research/articles"}
	firefoxArgs := []string{"--new-tab ", ""}
	builder := strings.Builder{}
	for i := 0; i <= len(newsSources)-1; i++ {
		firefoxArgs[1] = newsSources[i]
		fmt.Fprintln(os.Stdout, "Browsing to: ", firefoxArgs)
		openTabForMoreNews := exec.Command(firefoxCmd, firefoxArgs...)
		err = openTabForMoreNews.Start()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error:", err)
			panic(err)
		}

		time.Sleep(1 * time.Second)

		err = openTabForMoreNews.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			panic(err)
		}
		builder.Reset()
	}

	<-gracefulShutdown
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer handleTermination(cancel)
}
