#!/bin/bash

if [ "$#" -ne 3 ]; then
	echo "Usage: $0 <csvFilenameAndExtension.csv> <urlTotalGuess> <YoutubePlaylist URL>"
    echo "Input a filename with .csv as the extension" 
    echo "Guess how many links you are going extract as it will takes a long time to auto-manual scroll down to the bottom of a plyalist"
    echo ""
    echo "Because of laziness and handling authentication in bash is not an objective of this, please log in to your Account for Private Playlists"
    echo ""
    echo "Please beware this just a testing script to then develop a Golang application to do similar and the bulk parsing and processing of HTML and general decision making - this unsafe and will force Firefox to allow pasting and excution of JavaScript into the Developer Console"
    echo ""
    echo "This bash application uses a minimise JavaScript directly into FireFox to auto-manually scroll to the bottom of page and then recursively get all the links and titles of YouTube videos"
    echo "Minimised JavaScript:"
    echo 'clearInterval(goToBottom);let arrayVideos=[];console.log("\n".repeat(50));const links=document.querySelectorAll("a");for(const e of links)"video-title"===e.id&&(e.href=e.href.split("&list=")[0],arrayVideos.push(e.title+";"+e.href),console.log(e.title+"\t"+e.href));let data=arrayVideos.join("\n"),blob=new Blob([data],{type:"text/csv"}),elem=window.document.createElement("a");elem.href=window.URL.createObjectURL(blob),elem.download="YOURCSVFILEWILLGOHERE",document.body.appendChild(elem),elem.click(),document.body.removeChild(elem);'
	exit
fi
csvFile=$1
urlAmount=$2
youtubePlaylistURL=$3


# From Gzhodan.go
xdotoolHandle="xdotool"
xdtOpenTerminalAndFirefox=" key \"ctrl+alt+t\" sleep 1 type firefox" # needs xdtEnterKey
xdtFindFirefox=" search --onlyvisible --name firefox | head -n 1"
xdtClick=" key click 1 "                  
xdtDown=" key Down "                      
xdtTab=" key Tab "                        
xdtGotoURL=" key \"ctrl+l\" type " # needs xdtEnterKey
xdtEnterKey=" key \"Return\"" # key? 
xdtSavePageToPath=" key \"ctrl+s\" sleep 2 type " # needs xdtEnterKey
xdtCloseFirefox=" key --clearmodifiers \"ctrl+F4\""
# 

# Base64 and run:
b64JavaScriptCode=$(cat extractAllYTurls.min.js | sed "s/YOURCSVFILEWILLGOHERE/$csvFILE/g" | base64 -w0)
echo $b64JavaScriptCode

xdtOpenFFandConsole=' search --onlyvisible --class "firefox" windowactivate --sync key --clearmodifiers "ctrl+shift+k"'
# Bypass firefox anti scam 

# Commented lines do not work!!! 
bash -c "echo $xdotoolHandle $xdtOpenTerminalAndFireFox && $xdotoolHandle $xdtEnterKey"
bash -c 'sleep 3'
bash -c "$xdotoolHandle $xdtFindFirefox"
bash -c "$xdotoolHandle $xdtGotoURL $youtubePlaylistURL && $xdtoolHandle $xdtEnterKey"
bash -c 'sleep 3'
bash -c "$xdotoolHandle $xdtOpenFFandConsole"


exit 

xdtBypassFFConsoleScamProtection='type \"allow pasting\""'

# Press enter


# Type into console
xdtFFConsoleScrollToPageBottom=" type \"let goToBottom = setInterval(() => window.scrollBy(0, 400), 1000);\"" # Need xdtEnterKey
# Press enter
# Wait till reached the bottom
# for 1000~ it takes
sleep $urlAmount

xdtFFConsoleScrollToPageBottom=" type \"var decodedJavaScript = atob($b64JavaScriptForConsole);"

xdtFFConsoleEvalb64JavaScript="type \"eval(decodedJavaScript);\""




