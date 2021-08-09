# SuperFastjsonScan

该工具仅是Demo版，并不完善，给各位提供一个思路

该工具的核心是：不搭建JNDI Server或LDAP Server，也不用Dnslog平台，即可进行无回显Java反序列化漏洞的扫描（例如Fastjson）

原理在于解析RMI协议：
```go
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
```

使用方式：
```shell
./super-fastjson-scan -u http://127.0.0.1:8080/deserialize
   _____                       _____                 
  / ____|                     / ____|                
 | (___  _   _ _ __   ___ _ _| (___   ___ __ _ _ __  
  \___ \| | | | '_ \ / _ \ '__\___ \ / __/ _` | '_ \ 
  ____) | |_| | |_) |  __/ |  ____) | (_| (_| | | | |
 |_____/ \__,_| .__/ \___|_| |_____/ \___\__,_|_| |_|
              | |                                    
              |_|                                    
demo version by 4ra1n
[+] start listen at 127.0.0.1:8888
[+] find fastjson
[+] scan finish

Process finished with the exit code -1
```

更多原理参考链接：https://www.anquanke.com/post/id/249402
