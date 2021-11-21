# sakura_go

## 请求   

```
版本号
http://r3in.top:9091/version

排行榜
http://r3in.top:9091/rank

最近更新
http://r3in.top:9091/new

每日番剧
http://r3in.top:9091/weeks/{0-6}

关键字搜索
http://r3in.top:9091/search?page={page}&keyword={keyword}

详细
http://r3in.top:9091/detail/{id}

地址获取(结果为 base64 + URLEncode)
http://r3in.top:9091/play/{id}/{sid]/{nid]
```

## Examples   

```
http://r3in.top:9091/version
http://r3in.top:9091/rank
http://r3in.top:9091/new
http://r3in.top:9091/weeks/6
http://r3in.top:9091/search?page=1&keyword=凡人修仙传
http://r3in.top:9091/detail/fanrenxiuxianchuan
http://r3in.top:9091/play/fanrenxiuxianchuan/1/25
   
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/version
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/rank
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/new
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/weeks/6
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/search?page=1&keyword=凡人修仙传
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/detail/fanrenxiuxianchuan
https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/sakura.LATEST/sakurago/play/fanrenxiuxianchuan/1/25
```

## 阿里云函数计算   

上传 `release` 中的 `aliyunfc.zip` 即可，如果自己压缩打包则需要注意文件权限问题

## 其他   

仅供参考
