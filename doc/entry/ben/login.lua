function request()
    wrk.method = "POST"
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    wrk.body = "username=user_" .. math.random(10000010) .. "&password=123456"
    -- return wrk.format(wrk.method, wrk.path, wrk.headers, wrk.body)
    return wrk.format()
end

function response()

end
