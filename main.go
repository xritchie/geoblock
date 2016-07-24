package main

import(
	"fmt"
	"net"
	"math"
)

type IpNetContainer struct {
	ip net.IP
	ipNet *net.IPNet
}

func newIpNetContainer(ipNet string) (*IpNetContainer, error) {
	ip, ipnet, err := net.ParseCIDR(ipNet)
	if err != nil {
		return nil, err
	}

	ipAdd := ip.To4()
	if (ipAdd == nil) {
		ipAdd = ip.To16()
	}

	return &IpNetContainer{
		ip: ipAdd,
		ipNet: ipnet,
	}, nil
}

type IPBitTree struct {
	isEnd bool
	positive *IPBitTree
	negative *IPBitTree
}

func newIpBitTree(isEnd bool) *IPBitTree {
	return &IPBitTree{
		isEnd:isEnd,
		positive:nil,
		negative:nil,
	}
}

func getBitFromByte(byte byte, headPos uint8) bool {
	return (uint8(math.Pow(2, float64(8-(headPos+1)))) & byte) > 0
}

func getByteFromIp(bytes []byte, headPos uint8) byte {
	return bytes[headPos/8]
}

func (self *IPBitTree) UpdateIPBitTree(ipNetContainer IpNetContainer) {
	maskSize, _ := ipNetContainer.ipNet.Mask.Size()
	self.updateIPBitTree(ipNetContainer.ip, 0, uint8(maskSize))
}

func (self *IPBitTree) updateIPBitTree(bits []byte, headPos, maskSize uint8) {
	if (self.isEnd) {
		return
	}

	if (maskSize == 0) || (headPos >= maskSize) || (headPos >= (uint8(len(bits))*8)) {
		self.isEnd = true
		return
	}

	head := getBitFromByte(getByteFromIp(bits, headPos), headPos)

	switch {
	case (head) && (self.positive == nil):
		self.positive = newIpBitTree(len(bits) == 1)
	case (!head) && (self.negative == nil):
		self.negative = newIpBitTree(len(bits) == 1)
	}

	switch head {
	case true:
		self.positive.updateIPBitTree(bits, headPos+1, maskSize)
	case false:
		self.negative.updateIPBitTree(bits, headPos+1, maskSize)
	}
}

func (self *IPBitTree) IsMatch(ipNetContainer IpNetContainer) bool {
	return self.isMatch(ipNetContainer.ip, 0)
}

func (self *IPBitTree) isMatch(bits []byte, headPos uint8) bool {
	if ((self.isEnd) || (headPos >= uint8(len(bits)*8))) {
		return true
	}

	head := getBitFromByte(getByteFromIp(bits, headPos), headPos)

	switch head {
	case true:
		if (self.positive != nil) {
			return self.positive.isMatch(bits, headPos+1)
		} else {
			return false;
		}
	case false:
		if (self.negative != nil) {
			return self.negative.isMatch(bits, headPos+1)
		} else {
			return false;
		}
	}

	return false;
}

func main() {
	ipList := []string{"255.255.255.255/16", "255.255.255.255/24", "127.127.127.1/32",
		"192.168.1.1/24", "fd04:3e42:4a4e:3381::/64"}
	for i := 0; i < len(ipList); i += 1 {
		ipNetContainer, _ := newIpNetContainer(ipList[i])

		ipBitTree := newIpBitTree(false)
		ipBitTree.UpdateIPBitTree(*ipNetContainer)

		for _, ip := range ipList {
			ipNetContainer, _ := newIpNetContainer(ip)
			fmt.Printf("isMatch: %v, %v\n", ipBitTree.IsMatch(*ipNetContainer), ipNetContainer.ip)
		}
		fmt.Println("")
	}
}
