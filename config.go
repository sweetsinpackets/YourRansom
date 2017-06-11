package main

// 使用confGen生成加密后的配置文件数据
// confGen: https://github.com/YourRansom/confGen
var configE = "YOUR_CONFIG"

var configPw = "YOUR_PW" // 使用confGen生成配置文件时使用的密码

var (
	procNum  = 32
	jumpPer  = 2
	jumpHead int64 = 2048
	encSize int64 = jumpHead + (1024 * 1024)
)
