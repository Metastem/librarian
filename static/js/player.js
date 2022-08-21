const player = new Plyr('#player', {
  keyboard: { focused: true, global: true }
});

// Keyboard shortcuts
document.addEventListener('keydown', (event) => {
  if (event.target.id === "searchBar") return;
  event.preventDefault()
  switch (event.key) {
    case 'j':
      player.rewind(15);
      break;
    case ' ':
      player.togglePlay();
    case 'l':
      player.forward(15);
      break;
  }
});

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