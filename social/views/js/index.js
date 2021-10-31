window.onload = async () => {

    let birthdateToAge = (birthdate) => {
        let ageInMillis = Date.now() - new Date(birthdate);
        return parseInt(ageInMillis / 31536000000);
    }

    let renderFollowingList = (following) => {
        let renderFriend = (profile) => {
            let age = birthdateToAge(profile.birthdate);
            let profileString = `
                <tr id="followed-${profile.id}">
                    <td>${profile.firstname}</td>
                    <td>${profile.lastname}</td>
                    <td>${age}</td>
                    <td>${profile.city}</td>
                    <td><button type="button" class="btn btn-success" id="unfollow-btn-${profile.id}" data-pk="${profile.id}">Unfollow</button></td>
                </tr>`;
            $("tbody#following").append(profileString);
            document.querySelector(`button#unfollow-btn-${profile.id}`).addEventListener("click", unfollowBtnHandler);
        }
        $("tbody#following").empty();
        following.forEach(renderFriend)
    };

    let renderProfilePage = (profile, following) => {
        let myUsername = window.localStorage.getItem("username");
        let isMe = profile.username === myUsername;
        let age = birthdateToAge(profile.birthdate);

        $("div#profile-id").text(profile.id);
        $("div#first-name").text(profile.firstname);
        $("div#last-name").text(profile.lastname);
        $("div#age").text(age);
        $("div#gender").text(profile.gender);
        $("div#city").text(profile.city);
        $("div#interests").text(profile.interests);
        $("h4#full-name").text(profile.firstname + ' ' + profile.lastname);

        if (isMe) {
            if (following.length > 0) {
                renderFollowingList(following);
            }
            $("main#following-table").show();
            $("button#follow-btn").hide();
            $("button#message-btn").hide();
        } else {
            $("main#following-table").hide();
            $("button#follow-btn").disabled = false;
            $("button#follow-btn").show();
            $("button#message-btn").show();
        }
    };

    let followBtnHandler = async (e) => {
        e.preventDefault()
        let followed = $("div#profile-id").text();
        let res = await fetch(`/following/${followed}/follow`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            showMeContainer(null);
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
        document.getElementById("profile-container").style.display = "none";
    };

    let showMeContainer = async (e) => {
        document.getElementById("search-container").style.display = "none";
        document.getElementById("profile-container").style.display = "block";

        const res = await fetch("/profiles/me", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });

        if (res.ok) {
            let data = await res.json()
            renderProfilePage(data.profile, data.following);
        }
    };

    let showUserProfile = async (e) => {
        let profileId = e.target.attributes["data-pk"].value;
        const res = await fetch(`/profiles/${profileId}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            let data = await res.json()
            renderProfilePage(data.profile, []);
            document.getElementById("search-container").style.display = "none";
            document.getElementById("profile-container").style.display = "block";
        }
    };

    let showProfile = (profile) => {
        let age = birthdateToAge(profile.birthdate);
        let profileString = `
            <tr class='clickable-row'>
                <td>${profile.firstname}</td>
                <td>${profile.lastname}</td>
                <td>${age}</td>
                <td>${profile.city}</td>
                <td><button type="button" class="btn btn-success" id="profile-btn-${profile.id}" data-pk="${profile.id}">Link</button></td>
            </tr>`;
        $("tbody#profiles").append(profileString);

        document.querySelector(`button#profile-btn-${profile.id}`).addEventListener("click", showUserProfile);
    }

    let searchForPersons = async (term, cursor) => {
        const res = await fetch("/profiles/?" + new URLSearchParams({
            firstname: term,
            lastname: term,
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
        $("button#prev").hide();
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
    document.querySelector("a#logout").addEventListener("click", logout);
    document.querySelector("a#me-profile").addEventListener("click", showMeContainer);
    document.querySelector("form#search-form").addEventListener("submit", searchBtnHandler);
    document.querySelector("#search-field").addEventListener("input", searchTypeHandler);
    document.querySelector("button#next").addEventListener("click", nextBtnClick);
    document.querySelector("button#prev").addEventListener("click", prevBtnClick);
    document.querySelector("button#follow-btn").addEventListener("click", followBtnHandler);
};
