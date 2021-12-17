package util

import (
	"net"
	"os"
	"runtime"
)

type HostInfo struct {
	HostName string
	HostIp   string
}

func GetHostInfo() *HostInfo {
	hostName, _ := os.Hostname()
	ip, _ := getLocalIPv4Address()
	return &HostInfo{
		HostName: hostName,
		HostIp:   ip,
	}
}

func GetHostString() string {
	return StructToString(GetHostInfo())
}

func getLocalIPv4Address() (ipv4Address string, err error) {
	//获取所有网卡
	addrs, err := net.InterfaceAddrs()
	//遍历
	for _, addr := range addrs {
		//取网络地址的网卡的信息
		ipNet, isIpNet := addr.(*net.IPNet)
		//是网卡并且不是本地环回网卡
		if isIpNet && !ipNet.IP.IsLoopback() {
			ipv4 := ipNet.IP.To4()
			//能正常转成ipv4
			if ipv4 != nil {
				return ipv4.String(), nil
			}
		}
	}
	return
}

func GetProcessMemoryMB() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapAlloc) / 1024 / 1024
}
