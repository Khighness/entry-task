## entry-task benchmark test document



### login test

| client | QPS     |
| ------ | ------- |
| 200    | 5049.55 |
| 1000   | 7441.03 |
| 1500   | 8038.22 |
| 2000   | 7199.87 |



#### client: 200

```shell
$ wrk -t6 -c500 -d10s --latency -s login.lua "http://127.0.0.1:10001/login" 
Running 10s test @ http://127.0.0.1:10001/login
  6 threads and 500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    37.02ms   32.36ms 372.53ms   76.41%
    Req/Sec     1.37k   425.78     2.63k    69.89%
  Latency Distribution
     50%   29.25ms
     75%   53.04ms
     90%   77.05ms
     99%  144.91ms
  51018 requests in 10.10s, 25.04MB read
  Socket errors: connect 0, read 698, write 0, timeout 0
Requests/sec:   5049.55
Transfer/sec:      2.48MB
```

#### client: 1000

```shell
$ wrk -t6 -c1000 -d10s --latency -s login.lua "http://127.0.0.1:10001/login" 
Running 10s test @ http://127.0.0.1:10001/login
  6 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    58.78ms   42.35ms 360.99ms   65.25%
    Req/Sec     1.45k   420.87     3.13k    75.97%
  Latency Distribution
     50%   52.74ms
     75%   89.38ms
     90%  109.73ms
     99%  181.07ms
  75151 requests in 10.10s, 36.89MB read
  Socket errors: connect 0, read 3282, write 0, timeout 0
Requests/sec:   7441.03
Transfer/sec:      3.65MB
```

#### client: 1500

```shell
$ wrk -t6 -c1500 -d10s --latency -s login.lua "http://127.0.0.1:10001/login" 
Running 10s test @ http://127.0.0.1:10001/login
  6 threads and 1500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    97.05ms   55.44ms 534.52ms   70.83%
    Req/Sec     1.38k   433.63     2.52k    78.67%
  Latency Distribution
     50%   91.79ms
     75%  136.86ms
     90%  159.32ms
     99%  304.38ms
  81106 requests in 10.09s, 39.82MB read
  Socket errors: connect 0, read 8085, write 0, timeout 0
Requests/sec:   8038.22
Transfer/sec:      3.95MB
```

#### client: 2000

```shell
$ wrk -t6 -c2000 -d10s --latency -s login.lua "http://127.0.0.1:10001/login" 
Running 10s test @ http://127.0.0.1:10001/login
  6 threads and 2000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   133.32ms   90.27ms 870.32ms   74.75%
    Req/Sec     1.26k   495.81     2.19k    73.39%
  Latency Distribution
     50%  127.34ms
     75%  189.59ms
     90%  222.17ms
     99%  468.09ms
  72494 requests in 10.07s, 35.59MB read
  Socket errors: connect 0, read 18032, write 0, timeout 13
Requests/sec:   7199.87
Transfer/sec:      3.53MB
```





