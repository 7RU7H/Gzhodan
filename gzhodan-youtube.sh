# Guess how many links you are going extract
# It takes a long time
# For being non-obfuscatory  add both an echo for the parts
# And from singular js files for each


echo 'clearInterval(goToBottom);let arrayVideos=[];console.log("\n".repeat(50));const links=document.querySelectorAll("a");for(const e of links)"video-title"===e.id&&(e.href=e.href.split("&list=")[0],arrayVideos.push(e.title+";"+e.href),console.log(e.title+"\t"+e.href));let data=arrayVideos.join("\n"),blob=new Blob([data],{type:"text/csv"}),elem=window.document.createElement("a");elem.href=window.URL.createObjectURL(blob),elem.download="                                 ",document.body.appendChild(elem),elem.click(),document.body.removeChild(elem);'


# Base64 and run:
b64ConsolePartOne="Y2xlYXJJbnRlcnZhbChnb1RvQm90dG9tKTsKbGV0IGFycmF5VmlkZW9zID0gW107CmNvbnNvbGUubG9nKCdcbicucmVwZWF0KDUwKSk7CmNvbnN0IGxpbmtzID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgnYScpOwpmb3IgKGNvbnN0IGxpbmsgb2YgbGlua3MpIHsKaWYgKGxpbmsuaWQgPT09ICJ2aWRlby10aXRsZSIpIHsKCWxpbmsuaHJlZiA9IGxpbmsuaHJlZi5zcGxpdCgnJmxpc3Q9JylbMF07CglhcnJheVZpZGVvcy5wdXNoKGxpbmsudGl0bGUgKyAnOycgKyBsaW5rLmhyZWYpOwoJY29uc29sZS5sb2cobGluay50aXRsZSArICdcdCcgKyBsaW5rLmhyZWYpOwoJfQp9CmxldCBkYXRhID0gYXJyYXlWaWRlb3Muam9pbignXG4nKTsKbGV0IGJsb2IgPSBuZXcgQmxvYihbZGF0YV0sIHt0eXBlOiAndGV4dC9jc3YnfSk7CmxldCBlbGVtID0gd2luZG93LmRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoJ2EnKTsKZWxlbS5ocmVmID0gd2luZG93LlVSTC5jcmVhdGVPYmplY3RVUkwoYmxvYik7CmVsZW0uZG93bmxvYWQgPSAn"

b64ConsolePartOne="YnOwpkb2N1bWVudC5ib2R5LmFwcGVuZENoaWxkKGVsZW0pOwplbGVtLmNsaWNrKCk7CmRvY3VtZW50LmJvZHkucmVtb3ZlQ2hpbGQoZWxlbSk7"

echo "Input a filename with .csv as the extension" 

csvFilename2b64=$(echo $2 | base64 -w0)
b64JavaScriptForConsole=$($b64ConsolePartOne $csvFilename2b64 $b64ConsolePartTwo)
echo b64JavaScriptForConsole | base64 -d


xdotool search --onlyvisible --class "firefox" windowactivate --sync key --clearmodifiers 'ctrl+shift+k'

# By firefox anti scam 
# allow pasting 
# Press enter


# Type into console
'let goToBottom = setInterval(() => window.scrollBy(0, 400), 1000);'
# Press enter
# Wait till reached the bottom
# for 1000~ it takes

echo 'var decodedJavaScript = atob($b64JavaScriptForConsole);'

'eval(decodedJavaScript);'

