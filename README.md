# sub2clash
简单订阅转换

### 基本特点
- 支持 `ss`/`ssr`/`vmess` 协议的订阅
- 分组：全局 `load-balance` 分组和 HK/JP/US 各个地区的 `url-test` 分组
- 规则：GEOIP分流(HK/JP)，局域网、国内IP直连，默认 `load-balance` 分组，其他没了

### 基本功能
- 部署后可以快速获取订阅或转换订阅地址

### 基本用法
`base.yaml` 是模板文件，要放于同目录下
