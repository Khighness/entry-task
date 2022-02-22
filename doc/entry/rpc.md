## RPC设计



### 1. 传输协议

| header                     | content    |
| -------------------------- | ---------- |
| 头部[len(content)]，4 byte | 内容[data] |

```
Data
  ├─Name  方法名称
  ├─Args  方法参数 ｜ 返回数据
  └─Err   错误信息
```



### 2. 序列化

采用`golang`的`gob`进行序列化与反序列化。





