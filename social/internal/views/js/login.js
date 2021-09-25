window.onload = async () => {

    document.querySelector("form#login").addEventListener("submit", async e => {
        e.preventDefault();
        const username = e.target.querySelector("#username").value;
        const password = e.target.querySelector("#password").value;
        try {
            const res = await fetch("/auth/sign-in", {
                method: "POST",
                body: JSON.stringify({ username, password }),
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (res.ok) {
                let body = await res.json()
                body.cookie
                // Store users username via local storage & redirect to stories wall
                window.localStorage.setItem("username", username);
                window.localStorage.setItem("token", body["token"]);
                window.location = "/";
            } else if (res.status === 401) {
                e.target.querySelector("#error").innerText = "bad username/password combination";
            } else {
                e.target.querySelector("#error").innerText = `unknown error`;
            }
        } catch (error) {
            e.target.querySelector("#error").innerText = `unknown error: ${e}`;
        }
    });
}