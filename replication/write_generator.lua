math.randomseed(os.time())
request = function()
    local url_path = "/following/" .. math.random(1000000) .. "/follow"
    wrk.headers["Cookie"] = "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzYyOTEyNjIsImlkIjoxMDAwMDAxLCJvcmlnX2lhdCI6MTYzNjI4NzY2MiwidXNlcm5hbWUiOiJhZG1pbiJ9.2Al91XxR3CoxHadTfdk9mV1TwTScINc4Q26pmCAHfwY"
    return wrk.format("POST", url_path)
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