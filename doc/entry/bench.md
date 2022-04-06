## entry-task benchmark test document



### login test

| client | QPS     |
| ------ | ------- |
| 200    | 5623.24 |
| 1000   | 5326.15 |
| 1500   | 4152.13 |
| 2000   | 3622.45 |



#### client: 200

```shell
$ wrk -t6 -c200 -d10s --latency -s login.lua "http://127.0.0.1:10000/login" 
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    35.10ms    5.80ms  91.40ms   79.81%
    Req/Sec     0.94k   111.79     1.20k    69.33%
  Latency Distribution
     50%   34.78ms
     75%   37.46ms
     90%   41.31ms
     99%   52.17ms
  56357 requests in 10.02s, 27.94MB read
  Socket errors: connect 0, read 56, write 0, timeout 0
Requests/sec:   5623.24
Transfer/sec:      2.79MB

```

#### client: 1000

```shell
$ wrk -t6 -c1000 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   180.29ms   99.42ms   1.13s    64.97%
    Req/Sec     0.92k   526.45     2.43k    64.08%
  Latency Distribution
     50%  172.72ms
     75%  228.45ms
     90%  321.70ms
     99%  487.72ms
  53804 requests in 10.10s, 25.28MB read
  Socket errors: connect 0, read 3251, write 0, timeout 0
Requests/sec:   5326.15
Transfer/sec:      2.50MB
```

#### client: 1500

```shell
$ wrk -t6 -c1500 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 1500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   328.51ms  115.68ms   1.10s    84.69%
    Req/Sec   718.34    358.47     1.58k    64.59%
  Latency Distribution
     50%  302.30ms
     75%  342.56ms
     90%  450.89ms
     99%  770.59ms
  41790 requests in 10.06s, 20.46MB read
  Socket errors: connect 0, read 7925, write 0, timeout 0
Requests/sec:   4152.13
Transfer/sec:      2.03MB
```

#### client: 2000

```shell
$  wrk -t6 -c2000 -d10s --latency -s login.lua "http://127.0.0.1:10000/login"
Running 10s test @ http://127.0.0.1:10000/login
  6 threads and 2000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   432.74ms  108.63ms   1.05s    84.48%
    Req/Sec   666.14    393.20     1.52k    62.03%
  Latency Distribution
     50%  396.85ms
     75%  486.51ms
     90%  568.72ms
     99%  807.69ms
  36546 requests in 10.09s, 17.98MB read
  Socket errors: connect 0, read 16909, write 0, timeout 0
Requests/sec:   3622.45
Transfer/sec:      1.78MB
```



