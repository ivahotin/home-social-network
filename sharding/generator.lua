math.randomseed(os.time())
request = function()
    local chatId = math.random(1,100)
    local url_path = "/chats/" .. chatId .. "/messages"
    -- if we want to print the path generated
--     print(url_path)
    -- Return the request object with the current URL path
    wrk.method = "POST"
    wrk.headers["Content-Type"] = "application/json"
    wrk.headers["Cookie"] = "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzkyMjAzNDQsImlkIjoxMDAwMDAyLCJvcmlnX2lhdCI6MTYzOTIxNjc0NCwidXNlcm5hbWUiOiJpdmFob3RpbiJ9.2kcI-C6JI2jfI2NXl5ng1dmPyk8AlKh7uYFOj4n0Sks"
    wrk.body = '{"message": "'.. RandomVariable(20) ..'"}'
--     print(body)
    return wrk.format("POST", url_path, headers, body)
end

function RandomVariable(length)
	local res = ""
	for i = 1, length do
		res = res .. string.char(math.random(97, 122))
	end
	return res
end
