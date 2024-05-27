clearInterval();
function getAllURL() {
let arrayVideos = [];
let tmpLinks = document.querySelectorAll('a').length;
for (var i,cmpCount = 0; cmpCount === 3; i++){
let scrollInterval = setInterval(function() {
    window.scrollTo(0, document.body.scrollHeight);
    
    let cmpLinks = document.querySelectorAll('a').length;
    if (tmpLinks !== cmpLinks) {
        tmpLinks = cmpLinks;
    } else {
        cmpCount++;
    }
    }, 2000);
}
const links = document.querySelectorAll('a');
for (const e of links)
	'video-title' === e.id && (e.href = e.href.split('&list=')[0], arrayVideos.push(e.title + ';' + e.href), console.log(e.title + '\t' + e.href));
let data = arrayVideos.join('\n'), blob = new Blob([data], { type: 'text/csv' }), elem = window.document.createElement('a');
 
return data;
};
let data = getAllURL();
elem.href = window.URL.createObjectURL(blob), elem.download = '00iiraokanii00-test-playlist.csv', document.body.appendChild(elem), elem.click(), document.body.removeChild(elem);