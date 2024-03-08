#!/bin/bash

# -K $configFile.txt
# -s silent, -S show error, -O save as name
curl -X GET -L -A $USERAGENT -H "User-Agent: $USERAGENT" --connect-timeout 10 --retry 3 --retry-max-time 30 --retry-delay 60 -sS --remote-name-all $URL $URL 
