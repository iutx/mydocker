package network

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"mydocker/container"
	"net"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"
)

var (
	defaultNetworkPath = path.Join(container.BaseURL, "/networks/networks")
	drivers            = map[string]NetworkDriver{}
	networks           = map[string]*Network{}
)

type Network struct {
	Name    string     // 网络名
	IPRange *net.IPNet // 地址段
	Driver  string     // 驱动名
}

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"device"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	PortMapping []string         `json:"portmapping"`
	Network     *Network
}

type NetworkDriver interface {
	Name() string                                         // DriverName
	Create(subNet string, name string) (*Network, error)  // Create Network
	Delete(network Network) error                         // Delete Network
	Connect(network *Network, endpoint *Endpoint) error   // Connect container point to network
	Disconnect(network Network, endpoint *Endpoint) error // Remove container network point from network.
}

func (n *Network) dump(dumpPath string) error {
	if err := os.MkdirAll(dumpPath, 0644); err != nil {
		return err
	}
	configPath := path.Join(dumpPath, n.Name)
	jsonData, err := json.Marshal(n)
	if err != nil {
		log.Errorf("Marshal network object error: %v", err)
		return err
	}
	if err := ioutil.WriteFile(configPath, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

func (n *Network) load(loadPath string) error {
	if _, err := container.PathExists(loadPath); err != nil {
		return err
	}
	content, err := ioutil.ReadFile(loadPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, n); err != nil {
		log.Errorf("Load config to object error: %v", err)
		return err
	}
	return nil
}

func (n *Network) remove(networkPath string) error {
	configPath := path.Join(networkPath, n.Name)
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		return os.Remove(configPath)
	}
}

func CreateNetwork(driver, subnet, name string) error {
	/*
		Parse string to IPNET object.
		ex: 192.168.7.4/24 cdir.IP -》 192.168.7.0; cdir.Mask -》 ffffff00; cdir.String() -》 192.168.7.0/24;
	*/
	_, cdir, _ := net.ParseCIDR(subnet)

	// Get gateway ip from IPAM, The first ip normal.
	gatewayIP, err := ipAllocator.Allocate(cdir)
	if err != nil {
		return err
	}
	cdir.IP = gatewayIP
	// Create network use specified driver.
	nw, err := drivers[driver].Create(cdir.String(), name)
	if err != nil {
		return err
	}
	return nw.dump(defaultNetworkPath)
}

func Connect(networkName string, containerInfo *container.ContainerInfo) error {
	network, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network %s", networkName)
	}
	ip, err := ipAllocator.Allocate(network.IPRange)
	if err != nil {
		return err
	}
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", containerInfo.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: containerInfo.PortMapping,
	}
	if err := drivers[network.Driver].Connect(network, ep); err != nil {
		return err
	}
	if err := configEndpointIpAddressAndRoute(ep, containerInfo); err != nil {
		return err
	}
	return configPortMapping(ep, containerInfo)
}

func configPortMapping(ep *Endpoint, info *container.ContainerInfo) error {
	return nil
}

func configEndpointIpAddressAndRoute(ep *Endpoint, info *container.ContainerInfo) error {
	return nil
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, nw := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			nw.Name,
			nw.IPRange.String(),
			nw.Driver,
		)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func DeleteNetwork(networkName string) error {
	nw, ok := networks[networkName]
	if !ok {
		return fmt.Errorf("no such network name: %v", networkName)
	}
	// remove gateway ip.
	if err := ipAllocator.Release(nw.IPRange, &nw.IPRange.IP); err != nil {
		return fmt.Errorf("error remove network gateway ip: %v", err)
	}

	// delete network config and device by use driver's function.
	if err := drivers[nw.Driver].Delete(*nw); err != nil {
		return fmt.Errorf("remove network driver error: %v", err)
	}

	return nw.remove(defaultNetworkPath)
}

func Init() error {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver
	// Judge network config path exist or no.
	exist, err := container.PathExists(defaultNetworkPath)
	if err != nil {
		return err
	}
	if !exist {
		if err := os.MkdirAll(defaultNetworkPath, 0644); err != nil {
			return err
		}
	}
	// Walk函数会遍历指定路径下的所有文件，并且指定第二参数的函数设定
	if err := filepath.Walk(defaultNetworkPath, func(nwPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, nwName := path.Split(nwPath)
		nw := &Network{
			Name: nwName,
		}
		if err := nw.load(nwPath); err != nil {
			log.Errorf("error load network: %s", err)
		}

		networks[nwName] = nw
		return nil
	}); err != nil {
		return err
	}

	return nil
}
