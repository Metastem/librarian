async function sponsorblock() {
  let ytLinkDesc = document.querySelector(".description p:last-of-type").textContent;
  let ytId = ytLinkDesc.match(/(?:\.\.\..*v=)(.{11})/)
  if (ytId) {
    ytId = ytId[1]
    
    let hashedId = await sha256(ytId);
    let res = await fetch("/api/sponsorblock/" + hashedId.substring(0, 4) + "?categories=" + localStorage.getItem("sb_categories"))
    let data = await res.json()
    let videoData = data.find(v => v.videoID == ytId)

    let playerElement = document.getElementById("player")
    videoData.segments.forEach(segment => {
      playerElement.addEventListener('timeupdate', (event) => {
        const plyr = event.target.plyr;

        if (Math.floor(segment.segment[0]) == Math.floor(plyr.currentTime)) {
          plyr.forward(segment.segment[1] - plyr.currentTime)
        }
      });
    })

  }
}

if (localStorage.getItem("sb_categories")) {
  sponsorblock()
}

async function sha256(message) {
  const msgUint8 = new TextEncoder().encode(message);
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgUint8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map((b) => b.toString(16).padStart(2, '0')).join('');
  return hashHex;
}
