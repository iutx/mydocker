package network

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"net"
	"os"
	"path"
	"strings"
)

var ipamDefaultAllocatorPath = path.Join(container.BaseURL, "/networks/ipam/subnet.json")

type IPAM struct {
	SubnetAllocatorPath string
	Subnets             *map[string]string
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (i *IPAM) load() error {
	if _, err := os.Stat(i.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	configContent, err := ioutil.ReadFile(i.SubnetAllocatorPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(configContent, i.Subnets); err != nil {
		log.Errorf("unmarshal config data error: %v", err)
		return err
	}
	return nil
}

func (i *IPAM) dump() error {
	ipamPathDir, _ := path.Split(i.SubnetAllocatorPath)
	if _, err := os.Stat(ipamPathDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(ipamPathDir, 0644); err != nil {
				return fmt.Errorf("mkdir ipam path error: %v", err)
			}
		} else {
			return err
		}
	}
	content, err := json.Marshal(i.Subnets)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(i.SubnetAllocatorPath, content, 0644); err != nil {
		return err
	}

	return nil
}

func (i *IPAM) Allocate(subNet *net.IPNet) (ip net.IP, err error) {
	i.Subnets = &map[string]string{}
	if err := i.load(); err != nil {
		return nil, err
	}
	_, subNet, _ = net.ParseCIDR(subNet.String())
	// 127.0.0.0/8 -> 255.0.0.0 -> ret 8,24
	one, size := subNet.Mask.Size()
	if _, exist := (*i.Subnets)[subNet.String()]; !exist {
		(*i.Subnets)[subNet.String()] = strings.Repeat("0", 1<<uint8(size-one))
	}

	for c := range (*i.Subnets)[subNet.String()] {
		// if bit is zero. means ip can be use.
		if (*i.Subnets)[subNet.String()][c] == '0' {
			// change bit is 1, means this ip will be use.
			ipAlloc := []byte((*i.Subnets)[subNet.String()])
			ipAlloc[c] = '1'
			(*i.Subnets)[subNet.String()] = string(ipAlloc)

			// Get network ip; ex: 192.168.7.0/24 -》 192.168.7.0
			ip = subNet.IP

			// calculate ip will be return.
			for t := uint(4); t > 0; t -= 1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	fmt.Println(ip.To4(), "__________________________________")
	if err := i.dump(); err != nil {
		log.Errorf("ip :%v  dump error: %v", ip, err)
	}
	return
}

func (i *IPAM) Release(subNet *net.IPNet, ipAddr *net.IP) error {
	i.Subnets = &map[string]string{}

	_, subNet, _ = net.ParseCIDR(subNet.String())

	if err := i.load(); err != nil {
		log.Errorf("ipam load error: %v", err)
	}

	c := 0
	// translate ip to 4 bit.
	releaseIP := ipAddr.To4()
	releaseIP[3] -= 1
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subNet.IP[t-1]) << ((4 - t) * 8)
	}

	ipAlloc := []byte((*i.Subnets)[subNet.String()])
	ipAlloc[c] = '0'
	(*i.Subnets)[subNet.String()] = string(ipAlloc)

	if err := i.dump(); err != nil {
		return err
	}
	return nil
}
