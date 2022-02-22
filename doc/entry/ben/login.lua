function request()
    wrk.method = "POST"
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    wrk.body = "username=user_" .. (2 + math.random(10000008)) .. "&password=123456"
    return wrk.format()
end

function response()

end
