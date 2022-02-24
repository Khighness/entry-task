## 压测文档


### 登录压测

200并发

```shell
wrk -t6 -c200 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.11ms   27.10ms 369.94ms   86.32%
    Req/Sec     1.06k   279.62     1.61k    73.00%
  Latency Distribution
     50%   25.88ms
     75%   36.94ms
     90%   66.36ms
     99%  135.29ms
  63194 requests in 10.01s, 11.87MB read
  Socket errors: connect 0, read 61, write 0, timeout 0
Requests/sec:   6310.24
Transfer/sec:      1.19MB
```


2000并发

```shell
wrk -t6 -c2000 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 2000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    38.92ms   35.96ms 593.61ms   91.93%
    Req/Sec     1.06k   265.93     1.71k    69.33%
  Latency Distribution
     50%   28.89ms
     75%   41.86ms
     90%   67.41ms
     99%  198.54ms
  63246 requests in 10.02s, 11.88MB read
  Socket errors: connect 1751, read 108, write 0, timeout 0
Requests/sec:   6309.06
Transfer/sec:      1.19MB
```


wrk -t1 -c1 -d1s --latency -s login.lua "http://127.0.0.1:10000/login"
