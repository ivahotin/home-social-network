window.onload = async () => {

    document.querySelector("form#register").addEventListener("submit", async e => {
        e.preventDefault();
        const username = e.target.querySelector("#username").value;
        const password = e.target.querySelector("#password").value;
        const firstname = e.target.querySelector("#firstname").value;
        const lastname = e.target.querySelector("#lastname").value;
        const birthdate = e.target.querySelector("#birthdate").value + "T00:00:00Z";
        const gender = e.target.querySelector("#gender").value;
        const interests = e.target.querySelector("#interests").value;
        const city = e.target.querySelector("#city").value;
        try {
            const res = await fetch("/auth/sign-up", {
                method: "POST",
                body: JSON.stringify({ username, password, firstname, lastname, birthdate, gender, interests, city }),
                headers: {
                    "Content-Type": "application/json"
                }
            });
            if (res.ok) {
                window.location = "/auth/sign-in";
            } else if (res.status === 400) {
                e.target.querySelector("#error").innerText = "Such user already exists"
            } else {
                const { error } = await res.json();
                if (error) {
                    e.target.querySelector("#error").innerText = error.message;
                } else {
                    e.target.querySelector("#error").innerText = "unknown error";
                }

            }
        } catch (error) {
            e.target.querySelector("#error").innerText = "unknown error";
        }
    });
}