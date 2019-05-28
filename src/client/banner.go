package main

import (
	"encoding/binary"
	"encoding/json"
	"net"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
)

func getBanners(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.GetBannersRequest_CMD))
	conn.Write(cmd)
}

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

func clearAllBannerTimers(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.ClearAllBannerTimersRequest_CMD))
	conn.Write(cmd)
}
