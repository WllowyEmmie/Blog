const navbar = document.getElementById("navbar");
const postContainer = document.getElementById("post-container");
const userName = document.getElementById("user-name");
const API_BASE = "http://localhost:9090/api"
const API_BASE1 = "http://localhost:9090/api/all-posts";

const API_BASE2 = "http://localhost:9090/api/comments";
let posts = [];
const likeState = JSON.parse(localStorage.getItem("likeState")) || {};
const dislikeState = JSON.parse(localStorage.getItem("dislikeState")) || {};

const commentBodyInput = document.getElementById("comment-body");

window.addEventListener("scroll", () => {
    const maxScroll = 300;
    const scrollTop = window.scrollY;
    const opacity = Math.max(1 - scrollTop / maxScroll, 0);
    navbar.style.opacity = opacity
    if (opacity == 0) {
        userName.style.pointerEvents = "none";
    }
});

async function fetchUser() {
    try {
        const response = await fetch(`${API_BASE}/user`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });
        if (!response.ok) {
            throw new Error("Failed to fetch user")
        }
        const data = await response.json();
        userName.innerText = data.name;
    } catch (error) {
        console.log("Failed to fetch user : ", error);
        alert("Failed to fetch user");
    }

}
const signOutBtn = document.getElementById("signOut-btn");
signOutBtn.addEventListener("click", () => {
    localStorage.removeItem("token");
    localStorage.removeItem("likeState");
    localStorage.removeItem("dislikeState");
    window.location.href = "/frontend/login/login.html"
})
userName.addEventListener("click", () => {
    window.location.href = "/frontend/dashboard/dashboard.html";
});
async function fetchPost() {
    try {
        const response = await fetch(API_BASE1, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
        });
        const data = await response.json();
        if (!response.ok) {
            console.log("Server replies :", data);
            throw new Error("Failed to fetch posts");
        }
        posts = data.posts
        renderPosts();
    } catch (error) {
        console.log("Unable to render posts", error);
        alert("Failed to fetch posts")
    }
}
function renderPosts() {
    postContainer.innerHTML = "";
    posts.forEach(post => {
        const postID = post.id;

        const card = document.createElement("div");;
        card.style.minHeight = "30px";
        card.innerHTML = `
<div class="card shadow rounded mb-4" style="width: 100%; min-height: 60px; height: 100%; padding: 10px;">
    <div class="card-header d-flex justify-content-between align-items-center  text-white">
        <h5 class="mb-0">${post.title}</h5
        <small> ${post.user.name}</small>
    </div>
    <div class="card-body">
        <p class="card-text">${post.body}</p>
        
       <div class="d-flex align-items-center gap-4 justify-content-start mb-3">
  <div class="d-flex align-items-center gap-1">
    <button id="like-btn-${postID}" class="btn btn-outline-success btn-sm d-flex align-items-center p-1">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" clip-rule="evenodd"
          d="M13.75 0C13.1484 0 12.6127 0.352732 12.3657 0.888071L8.37734 9H3.98322C2.32771 9 1 10.3511 1 12V19C1 20.6489 2.32771 22 3.98322 22H18.6434C20.1225 22 21.3697 20.9128 21.5921 19.4549L22.9651 10.4549C23.2408 8.64775 21.854 7 20.0164 7H16.8825L16.8826 3.62846C16.8826 3.28115 16.8826 2.88736 16.8438 2.51718C16.8037 2.13526 16.7159 1.69889 16.4904 1.29245C15.9723 0.358596 14.9922 0 13.75 0ZM6 11H3.98322C3.44813 11 3 11.4398 3 12V19C3 19.5602 3.44813 20 3.98322 20H6V11Z"
          fill="#000000" />
      </svg>
    </button>
    <span id="post-reaction-like${postID}" class="ms-1 fw-semibold">0</span>
  </div>

  <!-- Dislike Section -->
  <div class="d-flex align-items-center gap-1">
    <button id="dislike-btn-${postID}" class="btn btn-outline-danger btn-sm d-flex align-items-center p-1">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
          d="M17.0001 2V13M22.0001 9.8V5.2C22.0001 4.07989 22.0001 3.51984 21.7821 3.09202C21.5903 2.71569 21.2844 2.40973 20.908 2.21799C20.4802 2 19.9202 2 18.8001 2H8.11806C6.65658 2 5.92584 2 5.33563 2.26743C4.81545 2.50314 4.37335 2.88242 4.06129 3.36072C3.70722 3.90339 3.59611 4.62564 3.37388 6.07012L2.8508 9.47012C2.5577 11.3753 2.41114 12.3279 2.69386 13.0691C2.94199 13.7197 3.4087 14.2637 4.01398 14.6079C4.70358 15 5.66739 15 7.59499 15H8.40005C8.96011 15 9.24013 15 9.45404 15.109C9.64221 15.2049 9.79519 15.3578 9.89106 15.546C10.0001 15.7599 10.0001 16.0399 10.0001 16.6V19.5342C10.0001 20.896 11.104 22 12.4659 22C12.7907 22 13.0851 21.8087 13.217 21.5119L16.5778 13.9502C16.7306 13.6062 16.807 13.4343 16.9278 13.3082C17.0346 13.1967 17.1658 13.1115 17.311 13.0592C17.4753 13 17.6635 13 18.0398 13H18.8001C19.9202 13 20.4802 13 20.908 12.782C21.2844 12.5903 21.5903 12.2843 21.7821 11.908C22.0001 11.4802 22.0001 10.9201 22.0001 9.8Z" />
      </svg>
    </button>
    <span id="post-reaction-dislike${postID}" class="ms-1 fw-semibold">0</span>
  </div>
</div>
        
        <span class="mt-4">Comments</>
        <div id="comments-container${postID}" class="mb-3"></div>
        
        <button class="btn btn-link p-0 mb-3" id="view-more-btn-${postID}">View more comments</button>
        
        <form id="comment-form-${postID}">
            <div class="mb-2">
                <textarea class="form-control mb-2" id="comment-body-${postID}" rows="2" placeholder="Write a comment..." required autocomplete="off"></textarea>
                <button id="comment-submit-button-${postID}" type="submit" class="btn btn-submit w-100">Submit</button>
            </div>
        </form>
    </div>
</div>
            `;
        postContainer.appendChild(card);
        fetchComments(postID);
        fetchreactions(postID);

        const likeButton = document.getElementById(`like-btn-${postID}`);
        likeButton.addEventListener("click", () => likePost(postID));
        const dislike = document.getElementById(`dislike-btn-${postID}`);
        dislike.addEventListener("click", () => dislikePost(postID));
        const commentForm = document.getElementById(`comment-form-${postID}`);
        const commentBody = document.getElementById(`comment-body-${postID}`);

        commentForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            try {
                const response = await fetch(`${API_BASE}/post/comment/${postID}`, {
                    headers: {
                        "Authorization": `Bearer ${localStorage.getItem("token")}`,
                        "Content-Type": "application/json"
                    },
                    method: "POST",
                    body: JSON.stringify({ body: commentBody.value }),
                })
                if (!response.ok) {
                    throw new Error("Could not send comment");
                }
                commentBody.value = "";
                fetchComments(postID);
            } catch (error) {
                console.log("Failed to post comments: ", error);
                alert("Failed to post comment")
            }
        });
        fetchComments(postID);
    })
}
async function fetchComments(postID) {
    try {


        const response = await fetch(`${API_BASE2}/${postID}`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
        });
        if (!response.ok) {
            throw new Error("Failed to fetch comments");
        }
        const data = await response.json();
        let comments = data.comments
        renderComments(postID, comments);
    } catch (error) {
        console.log("Failed to fetch Comments", error);
        alert("Failed to fetch comments");
    }

}
function renderComments(postID, comments) {
    const commentContainer = document.getElementById(`comments-container${postID}`);
    const viewMoreBtn = document.getElementById(`view-more-btn-${postID}`);
    commentContainer.innerHTML = "";
    let commentsToShow = 4;
    function showComments() {
        commentContainer.innerHTML = "";
        comments.slice(0, commentsToShow).forEach(comment => {
            const commentDiv = document.createElement("div")
            commentDiv.className = ""
            commentDiv.innerHTML = `

        <div class="card mb-3 shadow-sm">
            <div class="card-header d-flex justify-content-between align-items-center bg-light">
                <h6 class="mb-0 text-primary fw-semibold">${comment.user.name}</h6>
                <small class="text-muted">${comment.created_at || "Just now"}</small>
            </div>
             <div class="card-body py-2">
                <p class="card-text mb-0">${comment.body}</p>
             </div>
        </div>
        `;
            commentContainer.appendChild(commentDiv)
        });
        if (commentsToShow >= comments.length) {
            viewMoreBtn.style.display = "none";
        }
    }
    showComments();
    viewMoreBtn.onclick = () => {
        commentsToShow += 10;
        showComments();
    }

}
async function fetchreactions(postID) {
    try {
        const response = await fetch(`${API_BASE}/reactions/${postID}`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });
        if (!response.ok) {
            throw new Error("Failed to get Reactions");
        }
        const data = await response.json();
        const reactions = data.reactions;
        console.log(reactions)
        console.log(data);
        let totalLikes = 0;
        let totalDislikes = 0;
        for (const reaction of reactions) {
            totalLikes += reaction.likes;
            totalDislikes += reaction.dislikes
        }
        console.log(totalLikes);

        console.log(totalDislikes);
        const likenumber = document.getElementById(`post-reaction-like${postID}`);
        const dislikenumber = document.getElementById(`post-reaction-dislike${postID}`);
        likenumber.innerText = totalLikes;
        dislikenumber.innerText = totalDislikes;
    } catch (error) {
        console.log("Failed to get Reactions: ", error);
        alert("Failed to get reactions");
    }
}

async function likePost(postID) {
    let action;
    if (likeState[postID]) {
        action = "not-like"
        console.log("action1: ", action)
       await fetchreactions(postID);
    } else {
        action = "like"
        if (dislikeState[postID]) {
            await sendReaction(postID, "not-dislike")
            dislikeState[postID] = false;
            localStorage.setItem("dislikeState", JSON.stringify(dislikeState));
            await fetchreactions(postID)
            console.log("action2: ", action)
        }
    }
    try {
        await sendReaction(postID, action)
        likeState[postID] = !likeState[postID]
        localStorage.setItem("likeState", JSON.stringify(likeState));
    } catch (error) {
        console.log("Failed to like posts: ", error);
        alert("Failed to like posts");
    }
}
async function dislikePost(postID) {
    let action;
    if (dislikeState[postID]) {
        action = "not-dislike"
        await fetchreactions(postID);
        console.log("action3: ", action)
    } else {
        action = "dislike"
        if (likeState[postID]) {
            await sendReaction(postID, "not-like")
            likeState[postID] = false;
            localStorage.setItem("likeState", JSON.stringify(likeState));
            await fetchreactions(postID)
            console.log("action4: ", action)
        }
    }
    try {
        await sendReaction(postID, action)
        dislikeState[postID] = !dislikeState[postID]
        localStorage.setItem("dislikeState", JSON.stringify(dislikeState));
    } catch (error) {
        console.log("unable to dislike post: ", error)
        alert("Unable to dislike");
    }
}
async function sendReaction(postID, action) {
    try {
        const response = await fetch(`${API_BASE}/post/${postID}/reactions`, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`,
                "Content-Type": "application/json"
            },
            method: "PATCH",
            body: JSON.stringify({ action })
        })
        if (!response.ok) {
            throw new Error("Unable to change reaction")
        }
    } catch (error) {
        console.log("Unable to react: ", error);
        alert("Unable to react");
    }
}
window.addEventListener("DOMContentLoaded", () => { fetchPost(); fetchUser(); });
