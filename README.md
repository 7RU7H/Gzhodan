# README

Gzhodan is the G*-olang-*zombified Hyper-Optimized Daily Aggregating Newsreporter, because names are difficult, golang is good and Pseudo-Humaness is the objective. What is its purpose... PASS THE NEWS... remember:

> *A life without butter is no life at all*. - Marco Pierre White

#### Installation

For the purposes of which this is designed for, it is best to do as follows for a variety of implied reasons:
```bash
mkdir -p $HOME/go/src/github.com/7RU7H
cd /go/src/github.com/7RU7H
git clone git clone https://github.com/7RU7H/Gzhodan.git
cd Gzhodan/
# https://github.com/golang/go/issues/67424 - https://github.com/TLINDEN
go env -w GOTELEMETRY=off
go build  -ldflags="-s -w" -trimpath main.go
```

You may want to consider changing the browser and recompiling for the purposes of pattern reduction.
```bash
sed -i "s/$currentbrowser/$newbrower/g" main.go $HOME/go/src/github.com/7RU7H/main.go
# https://github.com/golang/go/issues/67424 - https://github.com/TLINDEN
go env -w GOTELEMETRY=off
go buildgo env -w GOTELEMETRY=off
go build  -ldflags="-s -w" -trimpath main.go
```

For the crontab at startup
```bash
# https://unix.stackexchange.com/questions/316486/run-a-script-at-startup-as-a-user
# https://serverfault.com/questions/1141641/setup-a-program-to-run-at-startup-as-a-specified-user-on-linux
crontab -e 
@reboot sleep 30 && /your/path/to/Ghzodan
```



If you want a more elaborate cronjob, above replacing `/your/path/to/Ghzodan` with a path to this script. Mkaing sure to add a list of `installed-browers-cli.txt` in this directory
```bash
#!/bin/bash 

finalgzhodanpath="/dev/shm"
# Edit browser in both browerCmd and xdtFindBrowser.. for Gzhodan to use a different Browser for OPSEC requirements if needed
currentbrowser=$(cat $HOME/go/src/github.com/7RU7H/main.go  | grep 'browserCmd                    string = ' | awk '{print $4}' | tr -d '\"')
# https://www.baeldung.com/linux/bash-draw-random-ints
min=1
max=$((wc -l $HOME/go/src/github.com/7RU7H/installed-browers-cli.txt | cut -d' ' -f1)) 
randomlinenumber=$((echo $(($RANDOM%($max-$min+1)+$min))))
newbrowser=$(awk "NR==$randomlinenumber{print}")
sed -i "s/$currentbrowser/$newbrower/g" main.go $HOME/go/src/github.com/7RU7H/main.go
# https://github.com/golang/go/issues/67424 - https://github.com/TLINDEN
go env -w GOTELEMETRY=off
go build  -ldflags="-s -w" -trimpath main.go
wait
# Set a directory you want to run anything from (replace /dev/shm)...
cp main $finalgzhodanpath/gzhodan
chmod +x $finalgzhodanpath/gzhodan
sleep 10
bash -c "$finalgzhodanpath/gzhodan"
```

