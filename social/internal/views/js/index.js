window.onload = async () => {

    let logout = async (e) => {
        e.preventDefault();
        window.localStorage.removeItem("username");
        try {
            await fetch("/auth/sign-out", { method: "POST" });
        } catch (e) { }
        window.location = "/auth/sign-in";
    };

    let showSearchContainer = async (e) => {
        document.getElementById("search-container").style.display = "block";
        document.getElementById("me-container").style.display = "none";
    };

    let showMeContainer = async (e) => {
        document.getElementById("search-container").style.display = "none";
        document.getElementById("me-container").style.display = "block";

        const res = await fetch("/profiles/me", {
            method: "GET",
            "headers": {
                "Content-Type": "application/json",
            }
        });

        if (res.ok) {
            let data = await res.json()
            $("td#first-name").text(data.firstname)
            $("td#last-name").text(data.lastname)
            $("td#age").text(data.age)
            $("td#gender").text(data.gender)
            $("td#city").text(data.city)
            $("td#interests").text(data.interests)
        }
    };

    let showProfile = (profile) => {
        let profileString = `
            <tr>
                <td>${profile.firstname}</td>
                <td>${profile.lastname}</td>
                <td>${profile.age}</td>
                <td>${profile.city}</td>
            </tr>
        `;
        $("tbody#profiles").append(profileString);
    }

    let searchForPersons = async (term, cursor) => {
        const res = await fetch("/profiles/?" + new URLSearchParams({
            term: term,
            cursor: cursor,
            limit: 10
        }));

        if (res.ok) {
            let body = await res.json();
            let profiles = body.profiles;
            let nextCursor = body.next_cursor;
            let prevCursor = body.prev_cursor;
            window.localStorage.setItem("prev_cursor", prevCursor - 1);
            window.localStorage.setItem("next_cursor", nextCursor);

            $("tbody#profiles").empty();
            profiles.forEach(showProfile);

            if (profiles.empty || nextCursor == 0) {
                $("button#next").hide();
            } else {
                $("button#next").show();
            }

            if (prevCursor <= 0) {
                $("button#prev").hide();
            } else {
                $("button#prev").show();
            }
        }

        await showSearchContainer(null);
    }

    let searchBtnHandler = async (e) => {
        e.preventDefault();
        let searchTerm = e.target.querySelector("#search-field").value;
        window.localStorage.setItem("currentTerm", searchTerm);
        await searchForPersons(searchTerm, 0);
    };

    let nextBtnClick = async (e) => {
        e.preventDefault();
        let term = window.localStorage.getItem("currentTerm");
        let cursor = window.localStorage.getItem("next_cursor");
        await searchForPersons(term, cursor);
    };

    let prevBtnClick = async (e) => {
        e.preventDefault();
        let term = window.localStorage.getItem("currentTerm");
        let cursor = window.localStorage.getItem("prev_cursor");
        await searchForPersons(term, cursor);
    };

    let searchTypeHandler = async (e) => {
        const maxWords = 2;
        let words = e.target.value.split(/\b[\s,\.-:;]*/);
        if (words.length > maxWords) {
            words.splice(maxWords);
            e.target.value = words.join("");
            alert("Max 2 words please");
        }
    };

    await showMeContainer(null)
    window.localStorage.setItem("prev_cursor", 0);
    window.localStorage.setItem("next_cursor", 0);
    $("button#prev").hide();
    document.querySelector("a#logout").addEventListener("click", logout);
    document.querySelector("a#me-profile").addEventListener("click", showMeContainer);
    document.querySelector("form#search-form").addEventListener("submit", searchBtnHandler);
    document.querySelector("#search-field").addEventListener("input", searchTypeHandler);
    document.querySelector("button#next").addEventListener("click", nextBtnClick);
    document.querySelector("button#prev").addEventListener("click", prevBtnClick);
};
