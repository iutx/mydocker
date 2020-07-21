package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
	"strings"
)

type BridgeNetworkDriver struct {
}

func (b *BridgeNetworkDriver) initBridge(n *Network) error {
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("create bridge interface %s error: %v", bridgeName, err)
	}
	gatewayIP := *n.IPRange
	gatewayIP.IP = n.IPRange.IP
	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("error assgining address: %s on bridge %s, error: %v", gatewayIP, bridgeName, err)
	}
	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("error set bridge up %s, error: %v", bridgeName, err)
	}
	if err := setupIPTables(bridgeName, n.IPRange); err != nil {
		return fmt.Errorf("error set iptables for %s, error: %v", bridgeName, err)
	}
	return nil
}

func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (b *BridgeNetworkDriver) Create(subNet string, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subNet)
	ipRange.IP = ip

	n := &Network{
		Name:    name,
		IPRange: ipRange,
		Driver:  b.Name(),
	}
	err := b.initBridge(n)
	if err != nil {
		log.Errorf("init bridge driver error: %v", err)
	}

	return n, err
}

func (b *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	return netlink.LinkDel(br)
}

func (b *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	return nil
}

func (b *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}

func createBridgeInterface(bridgeName string) error {
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName

	br := &netlink.Bridge{}
	br.LinkAttrs = la

	if err := netlink.LinkAdd(br); err != nil {
		return fmt.Errorf("bridge creation failed for bridge %s: %v", bridgeName, err)
	}

	return nil
}

func setInterfaceIP(bridgeName string, rawIP string) error {
	iface, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("error get interface %v", err)
	}

	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		return err
	}
	addr := &netlink.Addr{}
	addr.IPNet = ipNet
	addr.Label = ""
	addr.Flags = 0
	addr.Scope = 0
	return netlink.AddrAdd(iface, addr)
}

func setInterfaceUP(interfaceName string) error {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
	}

	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("error enabling interface for %s: %v", interfaceName, err)
	}
	return nil
}

func setupIPTables(bridgeName string, subNet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subNet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	//err := cmd.Run()
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
