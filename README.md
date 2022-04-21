# Proxy-Server

從免費proxy server獲取proxy，並驗證proxy可用性。

避免free proxy的可效性很短，每隔一分鐘就會重新更新proxy清單。

## Installation

``` 
go get github.com/tsukiamaoto/proxy-server-go
```

## Run

1. 先安裝redis server

資料會暫存在redis裡面，需要先啟動redis server。

照著在[redis document](https://redis.io/docs/getting-started/installation/)的文件安裝，然後啟動redis server。

2. 啟動程式

``` 
go run main.go
```

## Usage

* Api

啟動web services後，在默認的配置下訪問[http://127.0.0.1:8080/](http://127.0.0.1:8080/) 的API：

| Api           | Method | Description     | Params |
|---------------|--------|-----------------|--------|
| /api/v1/proxy | Get    | Get all proxies | None   |


## Free Proxy List

目前使用的ProxyList網站：

| Proxy Name      | URL                                   |
|-----------------|---------------------------------------|
| Free Proxy List | [URL](https://free-proxy-list.net/)   |