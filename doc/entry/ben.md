## 压测文档


### 登录压测

```shell
wrk -t6 -c200 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   161.07ms  173.74ms   1.00s    87.25%
    Req/Sec   273.15     74.41   535.00     74.33%
  Latency Distribution
     50%   97.70ms
     75%  191.52ms
     90%  396.14ms
     99%  826.48ms
  16275 requests in 10.04s, 9.07MB read
  Socket errors: connect 0, read 59, write 0, timeout 0
Requests/sec:   1620.53
Transfer/sec:      0.90MB
```

