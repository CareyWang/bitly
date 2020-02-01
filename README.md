# bitly

基于 golang1.13 与 bitly api 实现的短链接生成服务

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Install

安装依赖

```
make install
```

生成可执行文件，目录位于 build/ 。默认当前平台，其他平台请参照 Makefile 或执行对应 go build 命令。

```
make
```

## Usage

自行获取 bitly api token。前往 [Release](https://github.com/CareyWang/bitly/releases) 下载对应平台可执行文件。

```
./build/bitly.service -h
Usage of ./build/bitly.service:
  -cache int
      是否使用 redis 缓存
  -port int
    	服务端口 (default 8001)
  -token string
    	Bitly api token

./build/bitly.service -token xxxxxxxxxxxxxxxxxxx
```

建议配合 [pm2](https://pm2.keymetrics.io/) 开启守护进程。

```
pm2 start bitly.service --watch --name bitly -- --token xxxxxxxxxxxxxxxxxxxx
```

## API

```
GET /
```

### Parameters
|Name|Type|Description|
|---|---|---|
|longUrl|string|长链接|

### Response

```
Content-Type: application/json
Status: 200 OK 

{
    "Code":1,
    "Message":"",
    "LongUrl":"https://www.baidu.com",
    "ShortUrl":"http://bit.ly/38iQlfH"
}
```


## Maintainers

[@CareyWang](https://github.com/CareyWang)

## Contributing

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2020 CareyWang
