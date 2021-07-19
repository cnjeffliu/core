package netx

import (
	"fmt"
	"net"
	"testing"
)

func TestIP2UInt32(t *testing.T) {
	// IP地址转换为uint32
	IP1 := net.ParseIP("192.168.8.44")
	IPUint32 := IPToUInt32(IP1)
	fmt.Println(IPUint32)

	// uint32转换为IP地址
	IP2 := UInt32ToIP(IPUint32)
	fmt.Println(IP2.String())

	t.Log(net.SplitHostPort("10.0.2.113:9090"))
}
