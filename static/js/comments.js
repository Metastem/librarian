let commentsArr = [];

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