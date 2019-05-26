# Banner Manager
`Banner Manager` is an API server for managing banner.

## Features
- Written in Golang without any third-party package.
- Simple TCP client for connecting to manage.
- Unix timestamp to schedule display period.

## Quick Start

### Build the binary

```bash
make build
```

### Run the server

|Options|Default Value|
|:-:|:-:|
|-tcp|0.0.0.0:8080|
|-config|./config/local.json|

**Command:**

```bash
./banner-manager-server 
```
**For Example:**

```bash
./banner-manager-server -tcp=0.0.0.0:8080
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
./banner-manager-client
```

**For Example:**

```bash
./banner-manager-client -action=get
```
```bash
./banner-manager-client -action=update -serial=1 -start=1558849123 -expire=1561873123
```
```bash
./banner-manager-client -action=update_start -serial=1 -start=1558849123
```
```bash
./banner-manager-client -action=update_expire -serial=1 -expire=1561873123
```
```bash
./banner-manager-client -action=clear_all_timers
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

Response

```json
[
	{
		"serial": <serial(uint16)>,
		"event":  <event(string)>,
		"text":   <text(string)>,
		"image":  <image(string)>,
		"url":    <url(string)>,
	}
]
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

<a name="b3"></a>
**UpdateBannerStartedTimeRequest_CMD (0x0003)**

Request

```json
{
	"serial":      <serial(uint16)>,
	"startedTime": <start_time(uint32)>,
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

<a name="b5"></a>
**ClearAllBannerTimersRequest_CMD (0x0005)**

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
