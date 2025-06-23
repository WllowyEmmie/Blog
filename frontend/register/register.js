document.getElementById('register-form').addEventListener("submit", async function (e) {
    e.preventDefault();
    const name = document.getElementById("name").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const submitButton = document.getElementById("submit-button")

    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!name || !email || !password){
        alert("Please fill in all fields");
        return;
    }
    if(!emailRegex.test(email)){
        alert("Invalid email syntax")
        return;
    }

    submitButton.disabled = true;
    submitButton.textContent = "Registering..."
    try {
        const response = await fetch("http://localhost:9090/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name, email, password }),
        });
        const data = await response.json();
        if (response.ok) {
            alert("Successfully registered")
            window.location.href = "/frontend/login/login.html";
        }else{
            alert(data.error || "Couldnt register user")
        }
    } catch (error) {
        console.log("Registeration error: ", error)
        alert("Something went wrong!");
    } finally{
        submitButton.disabled = false;
        submitButton.textContent = "Register";
    }
});