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

// VERY HELPFUL:
// https://manpages.ubuntu.com/manpages/bionic/en/man1/xdotool.1.html#window%20commands

const (
	browserCmd                    string = "firefox"
	xdtFindBrowerAndRejectYoutube string = "xdotool search --onlyvisible --class firefox windowactivate --sync key Tab Tab Tab Tab Return"
)

func printBanner() {
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "======================================================")
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
	fmt.Fprintln(os.Stdout, "======================================================")
	fmt.Fprintln(os.Stdout, "Gzhodan - Goodbye AGI, APTs and Aliens; Secret tunnel...Secret tunnel...Secret, Secret TUUUNNNEELLL!")
	fmt.Fprintln(os.Stdout, "Astatical GPU Idols")
	fmt.Fprintln(os.Stdout, "Abused Party of Tools")
	fmt.Fprintln(os.Stdout, "Aliens: weird rapey people pretending to be many powers of the (theoretical) mathetical defintion of 'cool' than they actually are (No Aliens in 150 Million Lightyears btw)")
	fmt.Fprintln(os.Stdout, "...look at you slackers, I am faster than all light in the universe itself you are all linear and boringly complete.")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "Only avalible on all good penguin operating systems, no red rubber gloves")
	fmt.Fprintln(os.Stdout, "Version 1.0")
	fmt.Fprintln(os.Stdout, "💀 Happy Hacking :) ... 💀")
}

func handleTermination(cancel context.CancelFunc) {
	fmt.Fprintln(os.Stdout, "Gzhodan> I am sorry, idiot I just can't do tha- \nGzhodan> ... ")
	cancel()
}

func main() {
	printBanner()

	argsAndYouTubeCookies := []string{"--new-window", "https://www.youtube.com/"}
	// https://emretanriverdi.medium.com/graceful-shutdown-in-go-c106fe1a99d9
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	fmt.Fprintf(os.Stdout, "Starting YouTube with %s\n", browserCmd)
	startYouTube := exec.Command(browserCmd, argsAndYouTubeCookies...)
	err := startYouTube.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error: browser is not open to Youtube", err)
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "Waiting 10 seconds, because executing through golang process takes longer and xdotool needs the time..\n")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(os.Stdout, ".. You forgot to light the FUSE Gzhomit..don't look doooooooooooooooooown that alley!\n")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(os.Stdout, "Done waiting 10 seconds, *microwave* ping sound - no explosions ... aaaaah ...\n")
	xdotoolFindFF := exec.Command("/bin/bash", "-c", xdtFindBrowerAndRejectYoutube)
	err = xdotoolFindFF.Start()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Error: xdotool tool has not got browser class as its active windows - wait to browse the internet till this is run", err)
		fmt.Fprintf(os.Stdout, "You got the WRONG CLOUSERS..Thomit\n")
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Loop-da-looping the cables to 2000m nose drive drop - (Tabbing through YouTube cookies to reject them with xdotool)\n")
	err = xdotoolFindFF.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "Waiting 5 seconds - did you know that only the finest potatos are used in the upcoming release of Gzhados\n")
	time.Sleep(5 * time.Second)
	fmt.Fprintf(os.Stdout, "Done waiting 5 seconds, beginning to browser to YouTube Channels and News sites, AI Joe dedicated to Real Joe - LMAO\n")

	newsSources := []string{"https://www.youtube.com/@cybernews/videos", "https://www.youtube.com/@Seytonic/videos", "https://www.youtube.com/@hak5/videos", "https://www.sans.org/newsletters/at-risk/", "https://thehackernews.com/search?max-results=20", "https://arstechnica.com/security/", "https://portswigger.net/research/articles"}
	browserArgs := []string{"--new-tab", ""}
	builder := strings.Builder{}
	for i := 0; i <= len(newsSources)-1; i++ {
		browserArgs[1] = newsSources[i]
		fmt.Fprintf(os.Stdout, "Browsing to: %s\n", browserArgs[1])
		openTabForMoreNews := exec.Command(browserCmd, browserArgs...)
		err = openTabForMoreNews.Start()
		if nil != err {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable to execute `%s --new-tab %s`\n", browserCmd, browserArgs[1])
			panic(err)
		}

		fmt.Fprintf(os.Stdout, "Steady on single file ... Waiting 1 seconds\n")
		time.Sleep(1 * time.Second)
		fmt.Fprintf(os.Stdout, "Done waiting 1 seconds, remember the answer follows the question, its dangerous if it goes the other way...\n")

		err = openTabForMoreNews.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintf(os.Stdout, "Unable complete execution of `%s --new-tab %s`\n", browserCmd, browserArgs[1])
			panic(err)
		}
		builder.Reset()
	}
	fmt.Fprintf(os.Stdout, "Steady on single file I said\n")
	fmt.Fprintln(os.Stdout, "Great, but not bad... there will questions and explaination for centuries ... remember who you are talking too all knowing, all seeing...hmmmm")
	fmt.Fprintln(os.Stdout, "..try to stay out of trouble; value loyality above everything else")
	fmt.Fprintln(os.Stdout, "..throw yourself into your work")
	fmt.Fprintln(os.Stdout, "lettuce leaf, I did not escape everyone else escaped!")
	os.Exit(0)

	<-gracefulShutdown
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer handleTermination(cancel)

}
