window.onload = async () => {

    let followBtnHandler = async (e) => {
        e.preventDefault()
        let followed = e.target.attributes["data-pk"].value;
        const res = await fetch(`/following/${followed}/follow`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            e.target.disabled = true;
        }
    };

    let unfollowBtnHandler = async (e) => {
        e.preventDefault()
        let followed = e.target.attributes["data-pk"].value;
        const res = await fetch(`/following/${followed}/unfollow`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            document.querySelector(`tr#followed-${followed}`).remove()
        }
    }

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
            $("td#first-name").text(data.profile.firstname)
            $("td#last-name").text(data.profile.lastname)
            $("td#age").text(data.profile.age)
            $("td#gender").text(data.profile.gender)
            $("td#city").text(data.profile.city)
            $("td#interests").text(data.profile.interests)

            let renderFriend = (profile) => {
                let profileString = `
                    <tr id="followed-${profile.id}">
                        <td>${profile.firstname}</td>
                        <td>${profile.lastname}</td>
                        <td>${profile.age}</td>
                        <td>${profile.city}</td>
                        <td><button type="button" class="btn btn-success" id="unfollow-btn-${profile.id}" data-pk="${profile.id}">Unfollow</button></td>
                    </tr>`;

                $("tbody#following").append(profileString);

                document.querySelector(`button#unfollow-btn-${profile.id}`).addEventListener("click", unfollowBtnHandler);
            }

            $("tbody#following").empty();
            data.following.forEach(renderFriend)
        }
    };

    let showProfile = (profile) => {
        let profileString = `
            <tr>
                <td>${profile.firstname}</td>
                <td>${profile.lastname}</td>
                <td>${profile.age}</td>
                <td>${profile.city}</td>
                <td><button type="button" class="btn btn-success" id="follow-btn-${profile.id}" data-pk="${profile.id}">Follow</button></td>
            </tr>`;
        $("tbody#profiles").append(profileString);

        document.querySelector(`button#follow-btn-${profile.id}`).addEventListener("click", followBtnHandler);
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
            if (profiles.length != 0) {
                window.localStorage.setItem("prev_cursor", prevCursor);
                window.localStorage.setItem("next_cursor", nextCursor);
                $("button#next").show();
            } else {
                $("button#next").hide();
            }

            $("tbody#profiles").empty();
            profiles.forEach(showProfile);

            if (prevCursor == 0) {
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

    await showMeContainer(null);
    $("button#prev").hide();
    document.querySelector("a#logout").addEventListener("click", logout);
    document.querySelector("a#me-profile").addEventListener("click", showMeContainer);
    document.querySelector("form#search-form").addEventListener("submit", searchBtnHandler);
    document.querySelector("#search-field").addEventListener("input", searchTypeHandler);
    document.querySelector("button#next").addEventListener("click", nextBtnClick);
    document.querySelector("button#prev").addEventListener("click", prevBtnClick);
};
