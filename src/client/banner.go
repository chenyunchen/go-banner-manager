package main

import (
	"encoding/binary"
	"encoding/json"
	"net"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
)

// getBanners get banners and if the connection from white list, display inactive banner
func getBanners(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.GetBannersRequest_CMD))
	conn.Write(cmd)
}

// updateBanner update banner with started time and expired time
func updateBanner(conn net.Conn, serial uint16, startedTime, expiredTime uint32) {
	data := entity.UpdateBannerRequest{
		Serial:      serial,
		StartedTime: startedTime,
		ExpiredTime: expiredTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerRequest_CMD))
	conn.Write(append(output, b...))
}

// updateBannerStartedTime update banner started time and expired time is 2099/12/31 by default
func updateBannerStartedTime(conn net.Conn, serial uint16, startedTime uint32) {
	data := entity.UpdateBannerStartedTimeRequest{
		Serial:      serial,
		StartedTime: startedTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerStartedTimeRequest_CMD))
	conn.Write(append(output, b...))
}

// updateBannerExpiredTime update banner expired time if started time was set before
func updateBannerExpiredTime(conn net.Conn, serial uint16, expiredTime uint32) {
	data := entity.UpdateBannerExpiredTimeRequest{
		Serial:      serial,
		ExpiredTime: expiredTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerExpiredTimeRequest_CMD))
	conn.Write(append(output, b...))
}

// clearAllBannerTimers clear all the banner timer and all the display banner
func clearAllBannerTimers(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.ClearAllBannerTimersRequest_CMD))
	conn.Write(cmd)
}
