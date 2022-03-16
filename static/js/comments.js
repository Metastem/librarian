let commentsArr = [];
const filterKeywords = [
  "nigger",
	"racist",
	"leftist",
	"immigration",
	"immigrate",
	"liberal",
	"communist",
	"commie",
	"fuck",
	"cunt",
  "fag",
	"faggot",
	"retard",
	"woke",
	"bitch",
	"arab",
  "jew",
	"pussies",
	"pussy",
  "thug",
	"asshole",
	"entitled",
	"virtue",
	"signaling",
	"covid-19",
	"coronavirus",
	"vaccine",
	"soy",
	"suck",
	"hypocrite"
];

async function comments(claimId, channelId, channelName, page) {
  document.getElementById("spinner").style.display = "flex"
  let res = await fetch(`/api/comments?claim_id=${claimId}&channel_id=${channelId}&channel_name=${channelName}&page=${page}&page_size=15`);
  let data = await res.json();

  data.comments.forEach(comment => {
    commentsArr.push(comment)
  });

  renderComments()

  let comments = document.getElementById("comments").innerHTML;
  document.getElementById("comments").innerHTML = comments + `<a id="loadMore">Load more</a>`
  loadMoreBtn(page)
}

function renderComments() {
  commentsArr = commentsArr.sort((a, b) => (b.Likes - b.Dislikes) - (a.Likes - a.Dislikes))

  let commentsHTML = "";
  for(let i = 0; i < commentsArr.length; i++) {
    let comment = commentsArr[i];

    let pfpClass = "pfp"
    if(!comment.Channel.Thumbnail) {
      comment.Channel.Thumbnail = "/static/img/spaceman.png"
      pfpClass = "pfp pfp--default"
    } else {
      comment.Channel.Thumbnail = comment.Channel.Thumbnail + "&w=48&h=48"
    }
    
    let commentHTML = `
    <div class="comment">
      ${comment.Channel.Name !== "" ? `<a href="${comment.Channel.Url}">` : ""}
        <div class="videoDesc__channel">
          <img src="${comment.Channel.Thumbnail}" class="${pfpClass}" width="48" height="48" loading="lazy">   
          <p>
            ${
              comment.Channel.Title ?
              `<b>${comment.Channel.Title}</b><br>${comment.Channel.Name}`
              :
              comment.Channel.Name ? 
              `<b>${comment.Channel.Name}</b>`
              :
              "<b>[deleted]</b>"
            }
          </p>
        </div>
      ${comment.Channel.Name !== "" ? `</a>` : ""}
      <div>
        ${comment.Comment}
        ${comment.RelTime == "a long while ago" ? 
            `<p>
              ${comment.Time} |
              <span class="material-icons-outlined">thumb_up</span> ${comment.Likes}
              <span class="material-icons-outlined">thumb_down</span> ${comment.Dislikes}
            </p>`
          : `<p>
              <span title="${comment.Time}">${comment.RelTime}</span> |
              <span class="material-icons-outlined">thumb_up</span> ${comment.Likes}
              <span class="material-icons-outlined">thumb_down</span> ${comment.Dislikes}
            </p>`
        }
      </div>
    </div>`;

    commentsHTML = commentsHTML + commentHTML;
  }

  document.getElementById("comments").innerHTML = commentsHTML
  document.getElementById("spinner").style.display = "none"
}