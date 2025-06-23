const navbar = document.getElementById("navbar");
const postContainer = document.getElementById("post-container");
const userName = document.getElementById("user-name");
const API_BASE = "http://localhost:9090/api"
const API_BASE1 = "http://localhost:9090/api/all-posts";

const API_BASE2 = "http://localhost:9090/api/comments";
let posts = [];
const likeState = JSON.parse(localStorage.getItem("likeState")) || {};
const dislikeState = JSON.parse(localStorage.getItem("dislikeState")) || {}

const commentBodyInput = document.getElementById("comment-body");

window.addEventListener("scroll", () => {
    const maxScroll = 300;
    const scrollTop = window.scrollY;
    const opacity = Math.max(1 - scrollTop / maxScroll, 0);
    navbar.style.opacity = opacity
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
        if (!response.ok) {
            throw new Error("Failed to fetch posts");
        }
        const data = await response.json();
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

        const card = document.createElement("div");
        card.className = "card shadow-sm"
        card.style.width = "22rem"
        card.style.margin = "10px"
        card.innerHTML = `
        <div class = "card-body d-flex flex-column border-bottom mb-2 pb-3">
            <div class = "d-flex justify-content-between">
                <h2 class = "card-title text-right">${post.title}</h2>
                <h2 class = "card-title text-left">${post.user.name}</h2>
            </div>
                <p class = "card-text">${post.body}</p>
        
            <div class = "mb-3 d-flex justify-content-start">
                <span class = "me-3">
                    <button id = "like-btn-${postID}" class = "btn btn-outline-secondary btn-sm me-2">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-thumbs-up">
                            <path d="M14 9V5a3 3 0 0 0-6 0v4H4a2 2 0 0 0-2 2v7a2 2 0 0 0 2 2h8.28a2 2 0 0 0 1.94-1.47l1.38-5.53a2 2 0 0 0-.6-1.97L14 9z"></path>
                        </svg>
                    </button>
                    <span class ="likenumber" id = "post-reaction-like${postID}">0</span>
                </span>
                <span class = "me-3">
                    <button id = "dislike-btn-${postID}" class = "btn btn-outline-secondary btn-sm me-2">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-thumbs-down">
                            <path d="M10 15v4a3 3 0 0 0 6 0v-4h4a2 2 0 0 0 2-2v-7a2 2 0 0 0-2-2h-8.28a2 2 0 0 0-1.94 1.47L8.34 10.53a2 2 0 0 0 .6 1.97L10 15z"></path>
                        </svg>
                    </button>
                    <span class ="dislikenumber" id = "post-reaction-dislike${postID}">0</span>
                </span>
            </div>
            <div class = "mb-3 flex-grow-1">
                <h6 class = "text-center">Comments</h6>
                <div id = "comments-container${postID}" class ="mb-2"></div>
                <button class="btn btn-link p-0 mb-2"  id="view-more-btn-${postID}">View more comments</button>
                <form id = "comment-form-${postID}">
                    <div class="mb-3">
                            <label for="comment-body-${postID}" class="form-label">Body</label>
                            <textarea class="form-control" id="comment-body-${postID}" rows="2" required autocomplete="off"></textarea>
                        <button id = "comment-submit-button-${postID}" type="submit" class="btn btn-primary w-100 shadow">Submit</button>
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
            commentDiv.className = "card shadow-sm"
            commentDiv.innerHTML = `
        <div class = "card-title">
            <h3 class = "text-right">${comment.user.name}</h3> 
        </div>
        <div class = "card-body">
            <p>${comment.body}</p>
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
        const likes = data.likes;
        const dislikes = data.dislikes;
        const likenumber = document.getElementById(`post-reaction-like${postID}`);
        const dislikenumber = document.getElementById(`post-reaction-dislike${postID}`);
        likenumber.innerText = likes;
        dislikenumber.innerText = dislikes;
    } catch (error) {
        console.log("Failed to get Reactions: ", error);
        alert("Failed to get reactions");
    }
}

async function likePost(postID) {
    let action;
    if (likeState[postID]) {
        action = "not-like"
    } else {
        action = "like"
        if (dislikeState[postID]) {
            await sendReaction(postID, "not-dislike")
            dislikeState[postID] = false;
            localStorage.setItem("dislikeState", JSON.stringify(dislikeState));
            await fetchreactions(postID)
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
    } else {
        action = "dislike"
        if (likeState[postID]) {
            await sendReaction(postID, "not-like")
            likeState[postID] = false;
            localStorage.setItem("likeState", JSON.stringify(likeState));
            await fetchreactions(postID)
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
