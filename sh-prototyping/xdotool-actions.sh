bash -c "xdotool key \"ctrl+alt+t\" sleep 1 type firefox && xdotool key \"Return\"" # needs xdtEnterKey
bash -c "xdotool search --onlyvisible --name firefox | head -n 1"
bash -c "xdotool key click 1 "
bash -c "xdotool key Down "
bash -c "xdotool key Tab "
bash -c "xdotool key \"ctrl+l\" type "$INSERTKEYPRESSES" && xdotool key \"Return\"" # needs xdtEnterKey
bash -c "xdotool key \"Return\""
bash -c "xdotool key \"ctrl+s\" sleep 2 type "$INSERTKEYPRESSES" && xdotool key \"Return\"" # needs xdtEnterKey
bash -c "xdotool key --clearmodifiers \"ctrl+F4\""
bash -c "xdotool search --onlyvisible --class \"firefox\" windowactivate --sync key --clearmodifiers \"ctrl+shift+k\""
bash -c "xdotool type \"allow pasting\""
bash -c "xdotool type \"$INSERTKEYPRESSES\"" 
