## entry-task

<p align="center">
  <img src="https://img.shields.io/badge/go-backend-9cf?style=flat-square&logo=go" alt="go">
  <img src="https://img.shields.io/badge/vue-frontend-green?style=flat-square&logo=vue.js" alt="vue">
  <img src="https://img.shields.io/badge/docker-build-orange?style=flat-square&logo=docker" alt="docker">
  <img src="https://img.shields.io/badge/kubernetes-support-blueviolet?style=flat-square&logo=kubernetes" alt="kubernetes">
</p>

<p align="center">
  <a href="doc/entry/webapi.md">📑 WEBAPI</a> | <a href="doc/entry/deploy.md">🚀 DEPLOY</a> | <a href="doc/entry/bench.md">🛳 BENCH</a> 
</p>



### Feature

- Native http API based on Go lib
- Prefect connection pool for GRPC
- Beautiful front page built from VUE
- Elegant code style and exhaustive comments



### Structure

```
entry-task
    ├─bin              build script          
    ├─doc              mysql & docs
    ├─front            vue frontend
    ├─pb               grpc proto
    ├─rpc              rpc implement
    ├─tcp              tcp server
    └─web              web server
```



### Architecture

<p align="center">
<img src="./doc/images/architecture.svg" alt="Architecture" />
</p>


### Preview

<table>
  <tr>
    <td>login</td>
    <td>profile</td>
  </tr>
  <tr>
     <td width="50%" align="top"><img src="./doc/images/login.png"/></td>
     <td width="50%" align="top"><img src="./doc/images/profile.png"/></td>
  </tr>
</table>
