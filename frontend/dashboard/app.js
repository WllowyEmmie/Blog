const postsContainer = document.getElementById("posts-container");
const postForm = document.getElementById("post-form");
const postModal = new bootstrap.Modal(document.getElementById("post-modal"));
const addBtn = document.getElementById("add-post-btn")
const homeBtn = document.getElementById("home-btn")
const titleInput = document.getElementById("post-title");
const bodyInput = document.getElementById("post-body")
const postIDInput = document.getElementById("post-id")

const API_BASE = "http://localhost:9090/api/posts";
const API_BASE2 = "http://localhost:9090/api/post";
const API_BASE3 = "http://localhost:9090/api//post/:postID/delete"
let posts = [];

async function fetchPosts() {
    try {
        const response = await fetch(API_BASE, {
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });
        if (!response.ok) {
            throw new Error("Failed to fetch posts");
        }
        const data = await response.json();
        console.log("Fetched posts response:", data);
        posts = data.posts
        renderPosts();
    } catch (error) {
        console.error("Error fetching posts:", error);
        alert("Failed to load posts. Please check your connection.");
    }

}
function renderPosts() {
    postsContainer.innerHTML = "";
    posts.forEach(post => {
        const col = document.createElement("div");
        col.className = "col-lg-4";
        col.innerHTML = `
        <div class = "card shadow-sm">
            <div class = "card-body">
                <h5 class = "card-title">${post.title}</h5>
                <p class = "card-text">${post.body}</p>
                <button class="btn btn-sm btn-warning me-2 edit-btn" data-id="${post.id}">Edit</button>
                <button class="btn btn-sm btn-danger delete-btn" data-id="${post.id}">Delete</button>
            </div>
        </div>`;
        postsContainer.appendChild(col);
    });
}
postsContainer.addEventListener("click", (e) => {
    if (e.target.classList.contains("edit-btn")) {
        const id = e.target.getAttribute("data-id");
        editPost(id);
    }

    if (e.target.classList.contains("delete-btn")) {
        const id = e.target.getAttribute("data-id");
        deletePost(id);
    }
});
addBtn.addEventListener("click", () => {
    postIDInput.value = ""
    titleInput.value = ""
    bodyInput.value = ""
    postModal.show();
});
homeBtn.addEventListener("click", ()=>{
    window.location.href ="/frontend/main/main.html"
})
postForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const title = titleInput.value.trim();
    const body = bodyInput.value.trim();
    const id = postIDInput.value;

    const payload = { title, body };
    const url = id ? `${API_BASE}/${id}` : API_BASE2;
    const method = id ? "PATCH" : "POST";
    console.log("Submitting payload: ", payload)
    try {
            console.log("Submitting payload: ", payload)
        const response = await fetch(url, {
            method,
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            },
            body: JSON.stringify(payload)
        });
        const responseBOdy = await response.json();
        if (!response.ok) {
            console.log("Server replies :", responseBOdy);
            throw new Error("Save failed");
        }
        postModal.hide();
        fetchPosts();
    } catch (error) {
        console.error("Error saving post:", error);
        alert("Failed to save post. Please try again.");
    }
});
async function deletePost(id) {
    if (!confirm("Delete this post?")) return;

    try {
        const res = await fetch(`${API_BASE2}/${id}/delete`, {
            method: "DELETE",
            headers: {
                "Authorization": `Bearer ${localStorage.getItem("token")}`
            }
        });

        if (!res.ok) {
            throw new Error("Delete failed");
        }

        fetchPosts();
    } catch (error) {
        console.error("Error deleting post:", error);
        alert("Failed to delete post. Please try again.");
    }
}

function editPost(id) {
    const post = posts.find(p => p.id == id);
    if (!post)
        return;
    postIDInput.value = post.id;
    titleInput.value = post.title;
    bodyInput.value = post.body;
    postModal.show();
}

window.onload = fetchPosts;