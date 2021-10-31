math.randomseed(os.time())
request = function()
    local array = {
        "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r",
        "s", "t", "u", "v", "w", "x", "y", "z"
    }
    -- define the path that will search for q=%v 9%v being a random number between 0 and 1000)
    local firstName = array[math.random(26)] .. array[math.random(26)] .. array[math.random(26)]
    local lastName  = array[math.random(26)] .. array[math.random(26)] .. array[math.random(26)]
    local url_path = "/profiles?firstName=" .. encodeURI(firstName) .. "&lastName=" .. encodeURI(lastName) .. "&limit=100&cursor=0"
    -- if we want to print the path generated
--    print(url_path)
    -- Return the request object with the current URL path
    wrk.headers["Cookie"] = "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzU3MDYyNzksImlkIjoxMDAwMDAxLCJvcmlnX2lhdCI6MTYzNTcwMjY3OSwidXNlcm5hbWUiOiJhZG1pbiJ9._cUL3lJbpcMipchUry5wc2zonMUNnmm1VIxhSgcqcZU"
    return wrk.format("GET", url_path)
end

function encodeURI(str)
    if (str) then
        str = string.gsub (str, "\n", "\r\n")
        str = string.gsub (str, "([^%w ])",
            function (c) return string.format ("%%%02X", string.byte(c)) end)
        str = string.gsub (str, " ", "+")
    end
    return str
end