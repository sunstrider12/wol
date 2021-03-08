package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

func main() {
	macs := []string{"00-11-22-33-44-55"}
	for _, mac := range macs {
		wol(mac)
	}
}

//魔包构成: 前6位为FF 然后循环16次mac地址 后6位可选,为密码
//广播地址为:255.255.255.255 端口为9
func wol(mac string) {
	mac_str := strings.Replace(strings.Replace(mac, ":", "", -1), "-", "", -1)
	if len(mac_str) != 12 {
		fmt.Printf("mac%s地址错误\n", mac)
		return
	}
	macHex, err := hex.DecodeString(mac_str)
	if err != nil {
		fmt.Printf("mac%s地址错误\n", mac)
		return
	}
	var begin_cast = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	var buff bytes.Buffer
	buff.Write(begin_cast)
	for i := 0; i < 16; i++ {
		buff.Write(macHex)
	}
	mp := buff.Bytes()
	if len(mp) != 102 {
		fmt.Printf("mac%s地址错误\n", mac)
		return
	}
	sendMagicPacket(mp)
}

func sendMagicPacket(mp []byte) {
	sender := net.UDPAddr{}
	target := net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 9,
	}
	conn, err := net.DialUDP("udp", &sender, &target)
	if err != nil {
		fmt.Printf("创建UDP广播错误：%v\n", err)
		return
	}
	defer func() {
		conn.Close()
	}()

	_, err = conn.Write(mp)
	if err != nil {
		fmt.Printf("魔包发送失败[%s]\n", err)
	} else {
		fmt.Println("发送魔包成功")
	}
}
