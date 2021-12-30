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
  if (window.location.hash == "#filter") {
    filterKeywords.forEach(keyword => {
      commentsArr = commentsArr.filter(comment => !comment.Comment.toLowerCase().includes(keyword))
    })
    // Slur filter regex from https://github.com/LemmyNet/lemmy/blob/main/config/config.hjson, licensed under AGPL 3.0
    commentsArr.filter(comment => !comment.Comment.match(/(fag(g|got|tard)?\b|cock\s?sucker(s|ing)?|ni((g{2,}|q)+|[gq]{2,})[e3r]+(s|z)?|mudslime?s?|kikes?|\bspi(c|k)s?\b|\bchinks?|gooks?|bitch(es|ing|y)?|whor(es?|ing)|\btr(a|@)nn?(y|ies?)|\b(b|re|r)tard(ed)?s?)/gm))
  }

  renderComments()

  let comments = document.getElementById("comments").innerHTML;
  document.getElementById("comments").innerHTML = comments + `<a id="loadMore">Load more</a>`
  loadMoreBtn(page)
}

function renderComments() {
  console.log(commentsArr)
  commentsArr = commentsArr.sort((a, b) => (b.Likes - b.Dislikes) - (a.Likes - a.Dislikes))

  let commentsHTML = "";
  for(let i = 0; i < commentsArr.length; i++) {
    let comment = commentsArr[i];
    
    let commentHTML = `
    <div class="comment">
      ${comment.Channel.Thumbnail ? 
        `<img src="${comment.Channel.Thumbnail}&w=56&h=56" class="pfp" width="56" height="56" loading="lazy">`
        : `<img src="/static/img/spaceman.png" class="pfp pfp--default" width="56" height="56" loading="lazy">`
      }
      <div>
        <a href="${comment.Channel.Url}">
          <p>
            ${comment.Channel.Title ?
              `<b>${comment.Channel.Title}</b><br>${comment.Channel.Name}`
              : `<b>${comment.Channel.Name}</b>`
            }
          </p>
        </a>
        ${comment.Comment}
        ${comment.RelTime == "a long while ago" ? 
            `<p>
              ${comment.Time} |
              <span class="material-icons-outlined">thumb_up</span> ${comment.Likes}
              <span class="material-icons-outlined">thumb_down</span> ${comment.Dislikes}
            </p>`
          : `<p title="${comment.Time}">
              ${comment.RelTime} |
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