package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type ICMPEchoHeader struct {
	Type     uint8  //ICMP_Type
	Code     uint8  //ICMP_Code
	Checksum uint16 //CheckSum
	ID       uint16 //ICMP_Id
	Seq      uint16 //ICMP_Seq
}

func makeReqFormat(header ICMPEchoHeader) ([]byte, error) {
	data, err := time.Now().MarshalBinary()
	if err != nil {
		println(err)
		return nil, err
	}
	b := make([]byte, 8+binary.Size(data))
	b[0] = header.Type
	b[1] = header.Code
	b[2] = 0 //CheckSum
	b[3] = 0 //CheckSum
	binary.BigEndian.PutUint16(b[4:6], header.ID)
	binary.BigEndian.PutUint16(b[6:8], header.Seq)
	copy(b[8:], data)
	csum := checkSum(b)
	binary.BigEndian.PutUint16(b[2:4], csum)
	return b, nil
}

func getIPAddr(host string) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.To4(), err
		}
	}
	return nil, errors.New("Not found IPv4 addr")
}

func checkSum(buf []byte) uint16 {
	size := len(buf)
	sum := uint32(0) //uint32にする理由→checkSumする際に桁溢れを数えたいから
	//16bitずつICMPヘッダをTypeから16bitずつ取り出して加算する.

	for i := 0; i < size-1; i += 2 {
		sum += uint32(buf[i])<<8 | uint32(buf[i+1])
	}
	/*
		size:4
		i:0,2,4
		buf[0],buf[1],buf[2],buf[3] -> 全て足し合わせている
		4<4-1
		break
	*/
	/*
		size:5
		i:0,2,4
		buf[0],buf[1],buf[2],buf[3]
		problem: buf[4]がみられていない
		4<5-1
	*/
	//奇数バイトだった際は0でパディングして偶数にする
	if size%2 != 0 {
		sum += uint32(buf[size-1]) << 8
	}

	of := sum - (sum & 0xffff)
	sum = (sum & 0xffff) + (of >> 16)
	return ^uint16(sum)
}

func main() {
	ip, err := getIPAddr(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.Dial("ip4:1", ip.String())
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	pid := os.Getpid()
	seq := 0
	for {
		eHeader := ICMPEchoHeader{
			Type:     0x08, //ECHO要求
			Code:     0x00,
			Checksum: 0x00,
			ID:       uint16(pid),
			Seq:      uint16(seq),
		}
		b, err := makeReqFormat(eHeader)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, _ = conn.Write(b)
		time.Sleep(time.Second)
		seq++
	}
}
