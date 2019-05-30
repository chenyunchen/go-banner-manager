# Banner Manager
`Banner Manager` is an API server for managing banner.

## Features
- Written in Golang without any third-party package.
- Simple TCP client for connecting to manage.
- Unix timestamp to schedule display period.

## Quick Start

### Build the binary

```bash
> make build
```

### Run the server

|Options|Default Value|
|:-:|:-:|
|-tcp|0.0.0.0:8080|
|-config|./config/local.json|

**Command:**

```bash
> ./banner-manager-server 
```
**For Example:**

```bash
> ./banner-manager-server -tcp=0.0.0.0:8080
2019/05/27 21:49:33 Starting tcp server.
```

### Run the client

|Options|Default Value|More Detail|
|:-:|:-:|:-:|
|-tcp|0.0.0.0:8080||
|-action|get|[Here](#api)|
|-serial||[Here](#fake_data)|
|-start||[Unix Timestamp](https://timestampgenerator.com)|
|-expire||[Unix Timestamp](https://timestampgenerator.com)|

**Command:**

```bash
> ./banner-manager-client
```

**For Example:**

`get` (If there's no banner for update)

```bash
> ./banner-manager-client -action=get
2019/05/27 21:50:23 There is no banners available for display!
```
`update`

```bash
> ./banner-manager-client -action=update -serial=1 -start=1558849123 -expire=1561873123
2019/05/27 21:51:57 Display Banner:
2019/05/27 21:51:57 Serial: 1
2019/05/27 21:51:57 Event: Mercari Promotion
2019/05/27 21:51:57 Text: 20% off
2019/05/27 21:51:57 Image: https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png
2019/05/27 21:51:57 URL: https://www.mercari.com
2019/05/27 21:51:57 Started Time: 2019-05-26 13:38:43 +0800 CST
2019/05/27 21:51:57 Expired Time: 2019-06-30 13:38:43 +0800 CST
```
`update_start` (If the banner don't have expiredTime before)

```bash
> ./banner-manager-client -action=update_start -serial=2 -start=1558849123
2019/05/27 21:54:52 Expire time is not set yet!
```
`update_expire` (If the banner don't have startedTime before)

```bash
> ./banner-manager-client -action=update_expire -serial=3 -expire=1561873123
2019/05/27 22:32:29 Start time is not set yet!
```
`clear_all_timers`

```bash
> ./banner-manager-client -action=clear_all_timers
2019/05/27 22:02:01 There is no banners available for display!
```

## Internal IP Address Environment

### Build banner manager docker image

```bash
> make dockerfile.build
```

### Create docker subnet

```bash
> docker network create --subnet=10.0.0.0/24 mercari 
```

### Run banner manager server on `10.0.0.3`

```bash
> docker run --net mercari --ip 10.0.0.3 -it mercari/banner-manager sh 
/banner-manager # ./banner-manager-server 
2019/05/29 04:02:32 Starting tcp server.
```

### Run banner manager client on `10.0.0.1` or `10.0.0.2`

```bash
> docker run --net mercari --ip 10.0.0.2 -it mercari/banner-manager sh 
/banner-manager # ./banner-manager-client -tcp=10.0.0.3:8080
2019/05/29 04:03:46 There is no banners available for display!
```

## Test Case

### Set individually for each banner

|Serial|Event|Started Timestamp(ISO 8601)|Expired Timestamp(ISO 8601)|
|:-:|:-:|:-:|:-:|
|1|Mercari Promotion|1559072574(2019-05-29T03:42:54+08:00)|1561824000(2019-06-30T00:00:00+08:00)|
|2|Merpay Reward Point|1561938114(2019-07-01T07:41:54+08:00)|1564530114(2019-07-31T07:41:54+08:00)|
|3|Red Cross Donation|1564616514(2019-08-01T07:41:54+08:00)|1567208514(2019-08-31T07:41:54+08:00)|
|4|2020 Tokyo Olympics|1567294914(2019-09-01T07:41:54+08:00)|1569800514(2019-09-30T07:41:54+08:00)|
|5|Tokyo Golang Developers|1569886914(2019-10-01T07:41:54+08:00)|1572478914(2019-10-31T07:41:54+08:00)|

```bash
> ./banner-manager-client -action=update -serial=1 -start=1559072574 -expire=1561824000
> ./banner-manager-client -action=update -serial=2 -start=1561938114 -expire=1564530114
> ./banner-manager-client -action=update -serial=3 -start=1564616514 -expire=1567208514
> ./banner-manager-client -action=update -serial=4 -start=1567294914 -expire=1569800514
> ./banner-manager-client -action=update -serial=5 -start=1569886914 -expire=1572478914

2019/05/29 16:36:26 Display Banner:
2019/05/29 16:36:26 Serial: 1
2019/05/29 16:36:26 Event: Mercari Promotion
2019/05/29 16:36:26 Text: 20% off
2019/05/29 16:36:26 Image: https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png
2019/05/29 16:36:26 URL: https://www.mercari.com
2019/05/29 16:36:26 Started Time: 2019-05-29 03:42:54 +0800 CST
2019/05/29 16:36:26 Expired Time: 2019-06-30 00:00:00 +0800 CST
```

### Only one banner can be displayed at a time 

|Serial|Event|Started Timestamp(ISO 8601)|Expired Timestamp(ISO 8601)|
|:-:|:-:|:-:|:-:|
|1|Mercari Promotion|1559072574(2019-05-29T03:42:54+08:00)|1561824000(2019-06-30T00:00:00+08:00)|
|2|Merpay Reward Point|1559000514(2019-05-28T07:41:54+08:00)|1560555714(2019-06-15T07:41:54+08:00)|

```bash
> ./banner-manager-client -action=update -serial=1 -start=1559072574 -expire=1561824000
2019/05/29 16:49:55 Display Banner:
2019/05/29 16:49:55 Serial: 1
2019/05/29 16:49:55 Event: Mercari Promotion
2019/05/29 16:49:55 Text: 20% off
2019/05/29 16:49:55 Image: https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png
2019/05/29 16:49:55 URL: https://www.mercari.com
2019/05/29 16:49:55 Started Time: 2019-05-29 03:42:54 +0800 CST
2019/05/29 16:49:55 Expired Time: 2019-06-30 00:00:00 +0800 CST
> ./banner-manager-client -action=update -serial=2 -start=1559000514 -expire=1560555714
2019/05/29 16:50:20 Timestamp periods is overlap! Only one banner can be actived.
```

### Display the banner in period

**Display the banner for 1 minute.**

|Serial|Event|Started Timestamp(ISO 8601)|Expired Timestamp(ISO 8601)|
|:-:|:-:|:-:|:-:|
|1|Mercari Promotion|1559120400(2019-05-29T17:00:00+08:00)|1559120460(2019-05-29T17:01:00+08:00)|

```bash
> ./banner-manager-client -action=update -serial=1 -start=1559120400 -expire=1559120460 
2019/05/29 16:59:48 There is no banners available for display!
> ./banner-manager-client -action=get 
2019/05/29 17:00:18 Display Banner:
2019/05/29 17:00:18 Serial: 1
2019/05/29 17:00:18 Event: Mercari Promotion
2019/05/29 17:00:18 Text: 20% off
2019/05/29 17:00:18 Image: https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png
2019/05/29 17:00:18 URL: https://www.mercari.com
2019/05/29 17:00:18 Started Time: 2019-05-29 17:00:00 +0800 CST
2019/05/29 17:00:18 Expired Time: 2019-05-29 17:01:00 +0800 CST
> ./banner-manager-client -action=get
2019/05/29 17:01:05 There is no banners available for display!
```

### Display the future banner

**Client IP Address: `10.0.0.1` or `10.0.0.2`**

|Serial|Event|Started Timestamp(ISO 8601)|Expired Timestamp(ISO 8601)|
|:-:|:-:|:-:|:-:|
|2|Merpay Reward Point|1561938114(2019-07-01T07:41:54+08:00)|1564530114(2019-07-31T07:41:54+08:00)|

```bash
> ./banner-manager-client -tcp=10.0.0.3:8080 -action=update -serial=2 -start=1561938114 -expire=1564530114
2019/05/29 09:36:49 Display Banner:
2019/05/29 09:36:49 Serial: 2
2019/05/29 09:36:49 Event: Merpay Reward Point
2019/05/29 09:36:49 Text: 2% every transaction
2019/05/29 09:36:49 Image: https://jp.merpay.com/assets/homeServiceTeaser-summary.png
2019/05/29 09:36:49 URL: https://jp.merpay.com
2019/05/29 09:36:49 Started Time: 2019-06-30 23:41:54 +0000 UTC
2019/05/29 09:36:49 Expired Time: 2019-07-30 23:41:54 +0000 UTC
```

### Two active banner and display the expired time early one.

**Client IP Address: `10.0.0.1` or `10.0.0.2`**

|Serial|Event|Started Timestamp(ISO 8601)|Expired Timestamp(ISO 8601)|
|:-:|:-:|:-:|:-:|
|1|Mercari Promotion|1559072574(2019-05-29T03:42:54+08:00)|1561824000(2019-06-30T00:00:00+08:00)|
|2|Merpay Reward Point|1559000514(2019-05-28T07:41:54+08:00)|1560555714(2019-06-15T07:41:54+08:00)|

```bash
> ./banner-manager-client -tcp=10.0.0.3:8080 -action=update -serial=1 -start=1559072574 -expire=1561824000
2019/05/29 09:42:37 Display Banner:
2019/05/29 09:42:37 Serial: 1
2019/05/29 09:42:37 Event: Mercari Promotion
2019/05/29 09:42:37 Text: 20% off
2019/05/29 09:42:37 Image: https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png
2019/05/29 09:42:37 URL: https://www.mercari.com
2019/05/29 09:42:37 Started Time: 2019-05-28 19:42:54 +0000 UTC
2019/05/29 09:42:37 Expired Time: 2019-06-29 16:00:00 +0000 UTC
> ./banner-manager-client -tcp=10.0.0.3:8080 -action=update -serial=2 -start=1559000514 -expire=1560555714
2019/05/29 09:43:10 Display Banner:
2019/05/29 09:43:10 Serial: 2
2019/05/29 09:43:10 Event: Merpay Reward Point
2019/05/29 09:43:10 Text: 2% every transaction
2019/05/29 09:43:10 Image: https://jp.merpay.com/assets/homeServiceTeaser-summary.png
2019/05/29 09:43:10 URL: https://jp.merpay.com
2019/05/29 09:43:10 Started Time: 2019-05-27 23:41:54 +0000 UTC
2019/05/29 09:43:10 Expired Time: 2019-06-14 23:41:54 +0000 UTC
```

## Overview

* [1. Usage](#usage)
   * [1.1 Add More Handler](#u1)
   * [1.2 Add More Protocol](#u2)
   * [1.3 Instead of fake data to database](#u3)
* [2. Config](#config)
* [3. API](#api)
   * [3.1 Banner](#banner)
      * [3.1.1 GetBannersRequest_CMD](#b1)
      * [3.1.2 UpdateBannerRequest_CMD](#b2)
      * [3.1.3 UpdateBannerStartedTimeRequest_CMD](#b3)
      * [3.1.4 UpdateBannerExpiredTimeRequest_CMD](#b4)
      * [3.1.5 ClearAllBannerTimersRequest_CMD](#b5)
* [4. Fake Data](#fake_data)

<a name="usage"></a>
## Usage

<a name="u1"></a>
### Add More Handler

- Define Command: **src/entity/banner.go**
- Define Handler: **src/server/banner_handler.go**

```golang
// Init the router
router := NewRouter(sp)
router.Handle(entity.GetBannersRequest_CMD, NewHandler(getBannersHandler))
router.Handle(entity.UpdateBannerRequest_CMD, NewHandler(updateBannerHandler))
router.Handle(entity.UpdateBannerStartedTimeRequest_CMD, NewHandler(updateBannerStartedTimeHandler))
router.Handle(entity.UpdateBannerExpiredTimeRequest_CMD, NewHandler(updateBannerExpiredTimeHandler))
router.Handle(entity.ClearAllBannerTimersRequest_CMD, NewHandler(clearAllBannerTimersHandler))
```

<a name="u2"></a>
### Add More Protocol

Like: `TLS` `Websocket` `gRPC` ... etc.

- Implement the interface: **src/server/connection.go**
- For Example: **src/server/tcp.go**

```golang
// Conn is the interface of a general connection can write message
type Conn interface {
	LocalAddr() string
	RemoteAddr() string
	io.Reader
	io.Writer
	WriteMsg(msg string, byteOrder binary.ByteOrder, timeout time.Duration) error
	Close() error
}
```

<a name="u3"></a>
### Instead of fake data to database

Like: `MySQL` `MongoDB` `InfluxDB` ... etc.

- Implement the interface: **src/data/data.go**
- For Example: **src/data/fileData/fileData.go**

```golang
type BannersManager interface {
	GetBanners() ([]BannerInfo, error)
	GetBanner(uint16) (BannerInfo, error)
}
```

<a name="config"></a>
## Config

**config/local.json**

- Fake data path
- White list for internal debug

```json
{
	"data": "./data/fake.json",
	"whiteList": ["10.0.0.1", "10.0.0.2"]
}
```

<a name="api"></a>
## API

<a name="banner"></a>
### Banner

<a name="b1"></a>
**GetBannersRequest_CMD (0x0001)**

Request

```json
{}
```

Response

```json
{
	"serial":      <serial(uint16)>,
	"event":       <event(string)>,
	"text":        <text(string)>,
	"image":       <image(string)>,
	"url":         <url(string)>,
	"startedTime": <url(string)>,
	"expiredTime": <url(string)>,
}
```

<a name="b2"></a>
**UpdateBannerRequest_CMD (0x0002)**

Request

```json
{
	"serial":      <serial(uint16)>,
	"startedTime": <start_time(uint32)>,
	"expiredTime": <expire_time(uint32)>,
}
```

Response

```json
{
	"serial":      <serial(uint16)>,
	"event":       <event(string)>,
	"text":        <text(string)>,
	"image":       <image(string)>,
	"url":         <url(string)>,
	"startedTime": <url(string)>,
	"expiredTime": <url(string)>,
}
```

<a name="b3"></a>
**UpdateBannerStartedTimeRequest_CMD (0x0003)**

Request

```json
{
	"serial":      <serial(uint16)>,
	"startedTime": <start_time(uint32)>,
}
```

Response

```json
{
	"serial":      <serial(uint16)>,
	"event":       <event(string)>,
	"text":        <text(string)>,
	"image":       <image(string)>,
	"url":         <url(string)>,
	"startedTime": <url(string)>,
	"expiredTime": <url(string)>,
}
```

<a name="b4"></a>
**UpdateBannerExpiredTimeRequest_CMD (0x0004)**

Request

```json
{
	"serial":      <serial(uint16)>,
	"expiredTime": <expire_time(uint32)>,
}
```

Response

```json
{
	"serial":      <serial(uint16)>,
	"event":       <event(string)>,
	"text":        <text(string)>,
	"image":       <image(string)>,
	"url":         <url(string)>,
	"startedTime": <url(string)>,
	"expiredTime": <url(string)>,
}
```

<a name="b5"></a>
**ClearAllBannerTimersRequest_CMD (0x0005)**

Request

```json
{}
```

Response

```json
{}
```

<a name="fake_data"></a>
## Fake Data

**data/fake.json**

```json
[
  {
    "serial": 1,
    "event": "Mercari Promotion",
    "text": "20% off",
    "image": "https://www.mercari.com/jp/assets/img/common/jp/ogp_new.png",
    "url": "https://www.mercari.com"
  },
  {
    "serial": 2,
    "event": "Merpay Reward Point",
    "text": "2% every transaction",
    "image": "https://jp.merpay.com/assets/homeServiceTeaser-summary.png",
    "url": "https://jp.merpay.com"
  },
  {
    "serial": 3,
    "event": "Red Cross Donation",
    "text": "Hokkaido earthquakes 2018",
    "image": "http://www.jrc.or.jp/english/img/Hokkaidoearthquakeeng.jpg",
    "url": "http://www.jrc.or.jp/english/"
  },
  {
    "serial": 4,
    "event": "2020 Tokyo Olympics",
    "text": "Buy ticket",
    "image": "https://tokyo2020.org/en/assets/upload/20180718-01-top.jpg",
    "url": "https://tokyo2020.org/en/"
  },
  {
    "serial": 5,
    "event": "Tokyo Golang Developers",
    "text": "Golang Hackathon",
    "image": "https://secure.meetupstatic.com/photos/event/2/8/3/4/600_479470292.jpeg",
    "url": "https://www.meetup.com/Tokyo-Golang-Developers/"
  }
]
```
