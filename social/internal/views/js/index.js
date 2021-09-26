window.onload = async () => {
    // Handles logout
    await showMeContainer(null);
    document.querySelector("a#logout").addEventListener("click", logout);
    document.querySelector("a#me-profile").addEventListener("click", showMeContainer);
    document.querySelector("a#search").addEventListener("click", showSearchContainer);
};
