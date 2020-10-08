package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

type ICMPPacket struct {
	Type     uint8  //ICMP_Type
	Code     uint8  //ICMP_Code
	Checksum uint16 //CheckSum
	ID       uint16 //ICMP_ID
	Seq      uint16 //ICMP_Sequence
}

func makeReqFormat(header ICMPPacket) ([]byte, error) {
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
	return nil, err
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

func IcmpReq(conn net.Conn, ch chan time.Time) {
	seq := 0
	pid := os.Getpid()
	for {
		message := ICMPPacket{
			Type:     0x08, //ECHO要求
			Code:     0x00,
			Checksum: 0x00,
			ID:       uint16(pid),
			Seq:      uint16(seq),
		}
		b, err := makeReqFormat(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, _ = conn.Write(b)
		ch <- time.Now()
		seq++
		time.Sleep(time.Second)
	}
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

	ch := make(chan time.Time)
	go IcmpReq(conn, ch)
	resp := make([]byte, 100)
	for {
		//ECHOReplyを受け取るまでブロック
		duration := time.Since(<-ch)
		d, _ := conn.Read(resp)
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%s\n", d, os.Args[1], resp[27], duration)
	}
}
