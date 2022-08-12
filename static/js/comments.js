const commentData = JSON.parse(document.getElementById("commentData").innerText)

async function comments(claimId, channelId, channelName, page) {
  document.getElementById("spinner").style.display = "flex"
  let res = await fetch(`/api/comments?claim_id=${claimId}&channel_id=${channelId}&channel_name=${channelName}&page=${page}&page_size=10`);
  let data = await res.json();

  document.getElementById('commentsHeader').innerText = `Comments (${data.Items})`

  let commentsElem = document.getElementById("comments")
  for (let key in data.Comments) {
    let comment = data.Comments[key]
    if (document.getElementById(`${comment.CommentId}-replies`)) {
      continue;
    }
    commentsElem.appendChild(generateCommentElem(comment))
  }
  document.getElementById("spinner").style.display = "none"

  let loadMoreElem = document.createElement('a');
  loadMoreElem.id = 'loadMore';
  loadMoreElem.innerText = 'Load more';
  commentsElem.appendChild(loadMoreElem)

  page = data.Comments.length >= 10 ? page + 1 : page
  loadMoreBtn(page)
}

function generateCommentElem(comment) {
  let commentElem = document.createElement('div')
  commentElem.className = "comment"

  let channelElem = document.createElement('div')
  channelElem.className = "videoDesc__channel"
  commentElem.appendChild(channelElem)

  let pfpElem = document.createElement('img')
  let pfp = comment.Channel.Thumbnail
  pfpElem.className = !pfp ? "pfp pfp--default" : "pfp"
  pfpElem.src = !pfp ? "/static/img/spaceman.webp" : pfp + "&w=48&h=48"
  pfpElem.width = 48
  pfpElem.height = 48
  pfpElem.loading = "lazy"
  if (comment.Channel.Name) {
    let pfpChannelLink = document.createElement('a')
    pfpChannelLink.href = comment.Channel.Url
    pfpChannelLink.appendChild(pfpElem)
    pfpElem = pfpChannelLink
  }
  channelElem.appendChild(pfpElem)

  let channelTitleElem = document.createElement('p')
  let primaryTitleElem = document.createElement('b')
  channelTitleElem.appendChild(primaryTitleElem)
  if (comment.Channel.Title) {
    primaryTitleElem.innerText = comment.Channel.Title
    channelTitleElem.appendChild(document.createElement('br'))
    let secondaryTitleElem = document.createElement('span')
    secondaryTitleElem.innerText = comment.Channel.Name
    channelTitleElem.appendChild(secondaryTitleElem)
  } else if (comment.Channel.Name) {
    primaryTitleElem.innerText = comment.Channel.Name
  } else {
    primaryTitleElem.innerText = "[deleted]"
  }
  if (comment.Channel.Name) {
    let titleChannelLink = document.createElement('a')
    titleChannelLink.href = comment.Channel.Url
    titleChannelLink.appendChild(channelTitleElem)
    channelTitleElem = titleChannelLink
  }
  channelElem.appendChild(channelTitleElem)

  let wrapperDiv = document.createElement('div')
  wrapperDiv.innerHTML = comment.Comment
  commentElem.appendChild(wrapperDiv)

  let commentMetaElem = document.createElement('p')
  commentMetaElem.innerHTML = `${comment.Pinned ? `<span class="material-icons-outlined">push_pin</span> Pinned |`: ''}  ${comment.RelTime == "a long while ago" ? comment.Time : `<span title="${comment.Time}">${comment.RelTime}</span>`} | <span class="material-icons-outlined">thumb_up</span> ${comment.Likes} <span class="material-icons-outlined">thumb_down</span> ${comment.Dislikes}`
  wrapperDiv.appendChild(commentMetaElem)

  let repliesElem = document.createElement('div')
  repliesElem.className = "replies"
  repliesElem.id = `${comment.CommentId}-replies`
  commentElem.appendChild(repliesElem)

  if (comment.Replies) {
    let showReplyBtn = document.createElement('a')
    showReplyBtn.className = "showReplyBtn"
    showReplyBtn.innerText = comment.Replies == 1 ? "Show reply" : `Show ${comment.Replies} replies`
    showReplyBtn.addEventListener('click', async () => {
      if (showReplyBtn.innerText == "Loading replies…") return;
      if (showReplyBtn.innerText == "Hide replies") {
        showReplyBtn.innerText = comment.Replies == 1 ? "Show reply" : `Show ${comment.Replies} replies`
        for (let elem of Array.from(repliesElem.children)) {
          if (elem.className === "comment") elem.remove()
        }
        return
      }

      showReplyBtn.innerText = "Loading replies…"
      let res = await fetch(`/api/comments?parent_id=${comment.CommentId}&claim_id=${commentData.claimId}&channel_id=${commentData.channelId}&channel_name=${commentData.channelName}&page=1&page_size=100`);
      let data = await res.json();
      for (let key in data.Comments) {
        repliesElem.appendChild(generateCommentElem(data.Comments[key]))
      }
      showReplyBtn.innerText = "Hide replies"
    })
    repliesElem.appendChild(showReplyBtn)
  }

  return commentElem
}

function loadMoreBtn(page) {
  let loadMore = document.getElementById("loadMore");
  loadMore.addEventListener('click', () => {
    loadMore.remove()
    comments(commentData.claimId, commentData.channelId, commentData.channelName, page)
  })
}

const commentWarningBtn = document.getElementById('commentWarningBtn')
commentWarningBtn.removeAttribute('href')
commentWarningBtn.addEventListener('click', () => {
  document.getElementById('commentsWarning').style.display = 'none';
  comments(commentData.claimId, commentData.channelId, commentData.channelName, 1)
})