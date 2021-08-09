package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func startListen(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	localAddress, _ := net.ResolveTCPAddr("tcp4", address)
	l, err := net.ListenTCP("tcp", localAddress)
	if err != nil {
		panic(err)
	}
	output := fmt.Sprintf("start listen at %s:%d", host, port)
	Info(output)
	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	data := make([]byte, 1024)
	_, err = conn.Read(data)
	if err != nil {
		panic(err)
	}
	handleTCP(data, &conn)
}

func firstCheck(data []byte) bool {
	// check head
	if data[0] == 0x4a &&
		data[1] == 0x52 &&
		data[2] == 0x4d &&
		data[3] == 0x49 {
		// check version
		if data[4] != 0x00 &&
			data[4] != 0x01 {
			return false
		}
		// check protocol
		if data[6] != 0x4b &&
			data[6] != 0x4c &&
			data[6] != 0x4d {
			return false
		}
		// check other data
		lastData := data[7:]
		for _, v := range lastData {
			if v != 0x00 {
				return false
			}
		}
		return true
	}
	return false
}

func getFirstResp(conn *net.Conn) []byte {
	var ret []byte
	address := (*conn).RemoteAddr().String()
	ip := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]
	length := len(ip)
	// flag位
	ret = append(ret, 0x4e)
	// length位
	ret = append(ret, 0x00)
	ret = append(ret, uint8(length))
	// 写入ip
	for _, v := range ip {
		ret = append(ret, uint8(v))
	}
	// 空余
	ret = append(ret, 0x00)
	ret = append(ret, 0x00)
	intPort, _ := strconv.Atoi(port)
	temp := uint16(intPort)
	var b [2]byte
	// 写入端口
	b[1] = uint8(temp)
	b[0] = uint8(temp >> 8)
	ret = append(ret, b[0])
	ret = append(ret, b[1])
	return ret
}

func handleTCP(data []byte, conn *net.Conn) {
	// 检测第一个请求是否合法
	if !firstCheck(data) {
		return
	}
	// 发送IP信息的响应
	ret := getFirstResp(conn)
	_, err := (*conn).Write(ret)
	if err != nil {
		panic(err)
	}
	data = make([]byte, 1024)
	// 读取第二个请求
	_, _ = (*conn).Read(data)
	// 解析第二个请求
	parseSecond(data, conn)
}

func parseSecond(data []byte, conn *net.Conn) {
	if data[0] != 0x00 {
		return
	}
	length := data[1]
	var ip string
	for i := 2; i < int(length)+2; i++ {
		ip += fmt.Sprintf("%c", data[i])
	}
	// 判断响应的单播是否符合IPV4
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])$`
	match, _ := regexp.MatchString(ipReg, ip)
	if match {
		// 已经可以确认是rmi协议所以没必要进一步完成连接
		_ = (*conn).Close()
		socketChan <- true
	}
}
