# Dev

https://golongwithgolang.com/thread-safety-in-golang
https://gobyexample.com/goroutines
https://medium.com/@sairavitejachintakrindi/goroutines-and-threads-exploring-concurrency-in-go-370d609038c
https://medium.com/@brenomachadodomonte/multithreading-with-go-routines-in-golang-7e8fcd33be81
https://golangdocs.com/channels-in-golang

#### Objectives

- Globally deployable human-like news consumer-aggregator for real human consumption regard CyberSecurity news.
    - Go for global deployment
    - xtdotool and firefox for news consumption
    - Use GoSoup ... unforntunately the one non-std library cannot operate without
    - Gzlop for tokenised search through Gzhobins binary data map for later cross reference 
    - Develop human-like computer usage behvaiour library to extend bypassing bot and rate limiting checks 

Therefore the collection is not time sensitive the backend data aggregation, collection and retention is!

Meta/Sub Objectives
- go std as much as possible - few libraries where possible  
- Operate as much like a normal user (from the point of a remote site) - develop the Gzombie User library
- Bypass all future the rate limiting by smart path traversal of complete url list and being slow using a custom cli-browser & curl blend that has jitter and user-agent 
    - MAC and IP address randomisation  
- Use and develop Gzlop
- Prototype Gzhobins
    - Gzhobins Golang Zombie Hyper Optimise Binaries
        - Storage format to save artefacts to check the last visited with currect to speedly check difference    
- Gzombie user library

## Gzhombie Style Guide!

This I guess, is probably going to be required thinking and development as I go to prevent the de-modularisation by complex single chaining of commands that do more than one ACTION

ISSUE 1 - There is a mathematical law of 3s that states infinite combinations so maybe 2 subactions ending in a (key return) maximum is too strict.

ACTION = x + .. x^nth -> Y : where x is a suba[text](https://firefox-source-docs.mozilla.org/index.html)ctions (use xdotool to open a new terminal and do to a maxium of two things that result in a Return) and x^nth is a Return key (something terminates the set of actions as Action) that Y is simple Result; Whereas open terminal to open firefox to go to a url would be probably too much.

- ISSUE 2 - binary block memory execution that simulate uses is not test by me so the both size and execution and the librarification of blobs that imitate respresentations of users doing something on a computer is just creative theortical dirty protesting my hopes for getting good.

- ISSUE 3 - Even if EDR or AV alerts on use the bulk of the concept here could still just be injected into trust processes as shellcode that looks like user input but the assembly version of what a human is like on a computer is the ultimate 

- Bash 
- PowerShell
- C/C++/Rust 

## Issues

- Modularisation hopes are then shifted-up and confined to a dictionary of commands that do X - increasing binary size :( - and making everything less atomic
- Codium and sh is bad for development M'Kay lessons learnt.
- xdotool saveas click or GUI variation issues  


- Curl out a backup plan with curl

## Todo

- Variablisation is wrong - full commands required 


## The Plan

0. Map flowchart to dev chart of development - DONE
1. Get head and teeth into application steps after creative outbursts
    - Test and Implement as bash first for - visitSites - Testing found variablisation is the wrong way especial with Return key requirements and typing
        - Nice bash -> Golang os.cmd() for future use
2. Restructure App
    - Track app state with application struct
    - Config File 
        - Default Config
    - Segment by phase / purpose - monolith not possible
    - Config parsing
    - Application Logging
    - Storage Logging
    - Local Article parsing
    - Zhombrowsing
    - Application Directory and Data management
3. GoSoup - Article Parsing
4. Gzlop - Gzhobins
    - Backup plan if Gzhobins are too time consuming, or better alternative exists that will be maintained 
5. Intelligence Correlation ideas...

Pieces then Asynchoronous HARMONY!
- Logging
    - application.log - added - implement after completed all features
    - storage.log
- App-External Data retention
    - markdown newletter generation - wait
    - Directory stucture for: - Done
        - YEAR-MONTH-DAY.txt    
- Gzhobins
- Ingestion
    - Gzhobins
- Aggregation-Ingestion
    - Tokens as a const 
    - Memory arena hash map regex 
    - Check *correct article* storage.log file with glzop library
- Aggregation
    - Gzlop
    - Local Article Parsing
        - https://github.com/anaskhan96/soup - DO NOT USE TO GET ANYTHING - just copy functionality to parse 
- Intelligence Correlation
    - `firefox -screenshot $URL --headless --window-size=1920,1080`


- Gzombie
    - Action Sequencer
    - Jitter that match    

Post-Release
- Aethetics
- Site changes alert system
- X/Twitter issues
 


####  Dealing with the web application outliner
    - ycombinator Hacker News - non hackery tidbits by a set number of points 300+ because I only want the suggestion to go to another site, I am not going to comment, but other might and if they want lower the threshold for using this site they with a flag ASP or fork.
        - https://news.ycombinator.com/?p=1
        - https://news.ycombinator.com/?p=2
        - https://news.ycombinator.com/?p=3    
        - by X points

    - SAN @Risk
        - https://www.sans.org/newsletters/at-risk/
        - Roman numerals url /at-risk/
            - YEAR-ISSUE
            - 2024 volume 2 would be xxiv-02   
            - Sometimes there are 49-50?
            - Do some weeks in the year check on 48 - 52

- Daily Swig 
    - https://portswigger.net/daily-swig/vulnerabilities  - this web security news
    - https://portswigger.net/daily-swig/cloud-security
    - https://portswigger.net/daily-swig/supply-chain-attacks
    - https://portswigger.net/daily-swig/network-security
    - https://portswigger.net/daily-swig/zero-day

https://thehackernews.com/search?max-results=20

https://arstechnica.com/security/


titles links

- Ingestion
- Aggregate from modifiable configuration
    - For each site by url and date from previous log by yesterdays date - if no previous then grab, prevent previous urls continuously being featured. 
    - Token , case, US and UK, 
    - Tokens in the memory arenas fo 1.2X for the glzoping 
        - CVE, Exploit, Zero-Day, 0day,etc, Attack, Vuln ,Vulnerablilities, Malware, Ransomware, PoC, RCE, Bug , exploit,cve,rce,sqli,xss,lfi,ssti,ssrf,csrf,xxe,traversal,crlf,csv,injection,deserialization, JSON web Token, cookie (paramatre, parameter, prototype) pollution,smuggling, leaks, breach, errors, misconfiguration, API Key, crlf, cvs, clickjacking, DNS Rebinding, supply chain, Dom clubbering, file inclusion, insecure, Race condition, Tabnabbing, Type Juggling, (upload file - yikes), XSLT, Open Redirect, Mass Assignment, Crytography, Account takeover, Zero Click, Privilge Escalation, Decrypt, Encryption
        
        - Store all urls and summarised statements, Tokens hit! Date, - for self-reference in ingestion and for general collection
- Printing Terminal colourful friendly Newletter and Markdown friendly 
 

## Inspirations, Ideas, etc

https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go

https://github.com/abiyani/automate-save-page-as/blob/master/save_page_as - ./spa url flags firefox

https://www.educative.io/answers/how-to-implement-a-queue-in-golang
https://dev.to/vearutop/memory-arenas-in-go-j1f

https://github.com/anaskhan96/soup



- Use bug bounty screenshotter to get artitical screenshots for markdown report or just 
- GL with 3 year old web driver for golang that proabbly wont work https://github.com/tebeka/selenium 

## Post Release 

A dump of all future concerns to prevent the distraction of completing a working prototype in Go

#### Issues

- Dealing with Site changes - Site changed warning and alert system
- X/Twitter is the infosec hivemind, but mastodon and other exist..


#### Aethetics

- Gzhodan  || Gzhombr0z3
- Good picture like Scarecrow.go picture
## Historic

#### Timeline of TODO

Pieces then Asynchoronous HARMONY!

urls.txt

0. Develop the Gzombie
    - How to develop the HUMANESS of firefox/curl cli browsing the internet from the CLI  - DONE
        - Browser extensions - FIREFOX HAS EXTENSIONS
            - No JS and configure, etc
        - Jitter that match    
    - How to test the HUMANESS of firefox cli browsing the internet from the CLI?       
    - Queueing of tasks - DONE - ANSWER simple and dumb human open browser and links some links from a bookmark - over enigneering is bad
        - What is the safest traversal over all urls to avoid rate limiting (a(1-20),b(1,1,2,3,4,5)(1-5)...e(1)) - e(1) being SANS risk b being portswigger - CHILL
        - just a array sequencer rather than a graph - CHILL
1. Develop the ingestion
    - curl to buffer - FIREFOX AND XDOTOOL TO DOWNLOAD PAGES
    - MODIFY https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/net/http/client.go;l=483 - OVERENGINEERING NOT DOING
2. Develop the logging and data retention
    - application.log - added - implement after completed all features
    - storage.log
    - markdown newletter generation - wait
    - Diretory stucture for: - Done
        - YEAR-MONTH-DAY.txt    
3. Aggregation
    - glzop Library
    - Tokens as a const 
    - Memory arena hash map regex 
    - Check *correct article* storage.log file with glzop library
4. Article parsing 
    - https://github.com/anaskhan96/soup - DO NOT USE TO GET ANYTHING - just copy functionality to parse 
5. Intelligence correlation 
    - `firefox -screenshot $URL --headless --window-size=1920,1080`
        - Also very useful for handing to a AI to look at


#### Design Decisions and Ideas
- Find how scarecrow got there picture, needs a good repo picture. 
- The name is bad. The name needs to be on the Daily Swig level of naming - Gzhodan / Gzombr0ws3
- I want the data to be easy to just hand to another Go lang backend to an AI to do all AI related Software advances the Messlier is hyping me up for not have read emails and news for. A glorious future I will live to enjoy.

- package with a gzlop.go so no additional downloads and ensure development of gzlop
- package with gurl so that is developed.

- Python3 -> GO GO GO - gzlop for the win! for recursive grep-like in go routines and memory arena
- Just avoid the X/Twitter issues of some of infosec being X/Twitter and x/Twitter is still X/Twitter.

- Gzhobins Golang Zombie Hyper Optimise Binaries
    - Storage format to save artefacts to check the last visited with currect to speedly check difference    

