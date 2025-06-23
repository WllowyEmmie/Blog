document.getElementById("loginForm").addEventListener("submit", async function (e) {
    e.preventDefault();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const submitButton = document.getElementById("submit-button");
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

    
    if (!email || !password) {
        alert("Please fill both email and password fields");
        return;
    }
    if (!emailRegex.test(email)) {
        alert("Please enter a valid email")
        return
    }

    submitButton.disabled = true;
    submitButton.textContent = "Logging in...";
    try {
        const response = await fetch("http://localhost:9090/users/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password }),
        });
        const data = await response.json();
        if (response.ok) {
            alert("Login Successful!");
            localStorage.setItem("token", data.token);
            window.location.href = "/frontend/main/main.html";
        } else {
            alert(data.error || "Login failed!");
        }
    } catch (error) {
        console.log("Login error : ", error);
        alert("Something went wrong!");
    }finally{
        submitButton.disabled = false;
        submitButton.textContent = "Login";
    }
});