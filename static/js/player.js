const player = new Plyr('#player');

if (localStorage.getItem("autoplay") === "true") {
  player.on('ready', player.play())
}

if (localStorage.getItem("autoplayNextVid") === "true") {
  let nextVid = document.getElementsByClassName("relVid__link")
  nextVid = nextVid[0].getAttribute("href")
  
  player.on('ended', () => {
    window.location.href = nextVid
  })
}

if (location.hash) {
  player.on('loadeddata', () => {player.currentTime = location.hash.replace("#", "") * 1})
}

const urlParams = new URLSearchParams(location.search);
if (urlParams.get("t")) {
  player.on('loadeddata', () => {player.currentTime = urlParams.get("t") * 1})
}

window.addEventListener('hashchange', () => {
  player.currentTime = location.hash.replace("#", "") * 1
})