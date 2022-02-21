## 压测文档


### 登录压测

```shell
zikang.chen@C02GM1FFMD6M ben % wrk -t6 -c200 -d10s --latency -s login.lua "http://127.0.0.1:10000/login" 
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 200 connections
^C  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   195.99ms  224.37ms   1.01s    87.06%
    Req/Sec   225.77    112.20   646.00     67.66%
  Latency Distribution
     50%  106.23ms
     75%  256.62ms
     90%  529.80ms
     99%    1.00s 
  9842 requests in 7.34s, 2.51MB read
  Socket errors: connect 0, read 55, write 0, timeout 0
Requests/sec:   1340.86
Transfer/sec:    349.53KB
```

