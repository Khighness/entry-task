function request()
    wrk.method = "POST"
    wrk.headers["content-type"] = "application/json"
    username = "user_" .. (2 + math.random(10000000))
    password = "123456"
    wrk.body = string.format('{"username":"%s","password":"%s"}', username, password)
    return wrk.format()
end
