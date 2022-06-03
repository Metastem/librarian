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