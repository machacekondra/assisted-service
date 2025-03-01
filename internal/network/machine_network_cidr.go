package network

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/openshift/assisted-service/internal/common"
	"github.com/openshift/assisted-service/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getVIPInterfaceNetwork(vip net.IP, addresses []string) *net.IPNet {
	for _, addr := range addresses {
		_, ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		if ipnet.Contains(vip) {
			return ipnet
		}
	}
	return nil
}

/*
 * Calculate the machine network CIDR from the one of (ApiVip, IngressVip) and the ip addresses of the hosts.
 * The ip addresses of the host appear with CIDR notation. Therefore, the network can be calculated from it.
 * The goal of this function is to find the first network that one of the vips belongs to it.
 * This network is returned as a result.
 */
func CalculateMachineNetworkCIDR(apiVip string, ingressVip string, hosts []*models.Host) (string, error) {
	var ip string
	if apiVip != "" {
		ip = apiVip
	} else if ingressVip != "" {
		ip = ingressVip
	} else {
		return "", nil
	}
	isIPv4 := IsIPv4Addr(ip)
	parsedVipAddr := net.ParseIP(ip)
	if parsedVipAddr == nil {
		return "", errors.Errorf("Could not parse VIP ip %s", ip)
	}
	for _, h := range hosts {
		if swag.StringValue(h.Status) == models.HostStatusDisabled {
			continue
		}
		var inventory models.Inventory
		err := json.Unmarshal([]byte(h.Inventory), &inventory)
		if err != nil {
			continue
		}
		for _, intf := range inventory.Interfaces {
			var ipnet *net.IPNet
			if isIPv4 {
				ipnet = getVIPInterfaceNetwork(parsedVipAddr, intf.IPV4Addresses)
			} else {
				ipnet = getVIPInterfaceNetwork(parsedVipAddr, intf.IPV6Addresses)
			}
			if ipnet != nil {
				return ipnet.String(), nil
			}
		}
	}
	return "", errors.Errorf("No suitable matching CIDR found for VIP %s", ip)
}

func ipInCidr(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}
	return ipnet.Contains(ip)
}

func VerifyVip(hosts []*models.Host, machineNetworkCidr string, vip string, vipName string, mustExist bool, log logrus.FieldLogger) error {
	if !mustExist && vip == "" {
		return nil
	}
	if !ipInCidr(vip, machineNetworkCidr) {
		return errors.Errorf("%s <%s> does not belong to machine-network-cidr <%s>", vipName, vip, machineNetworkCidr)
	}
	if !IpInFreeList(hosts, vip, machineNetworkCidr, log) {
		return errors.Errorf("%s <%s> is already in use in cidr %s", vipName, vip, machineNetworkCidr)
	}
	return nil
}

func verifyDifferentVipAddresses(apiVip string, ingressVip string) error {
	if apiVip == ingressVip && apiVip != "" {
		return errors.Errorf("api-vip and ingress-vip cannot have the same value: %s", apiVip)
	}
	return nil
}

func VerifyVips(hosts []*models.Host, machineNetworkCidr string, apiVip string, ingressVip string, mustExist bool, log logrus.FieldLogger) error {
	err := VerifyVip(hosts, machineNetworkCidr, apiVip, "api-vip", mustExist, log)
	if err == nil {
		err = VerifyVip(hosts, machineNetworkCidr, ingressVip, "ingress-vip", mustExist, log)
	}
	if err == nil {
		err = verifyDifferentVipAddresses(apiVip, ingressVip)
	}
	return err
}

func VerifyMachineCIDR(machineCidr string, hosts []*models.Host, log logrus.FieldLogger) error {
	ip, ipNet, err := net.ParseCIDR(machineCidr)
	if err != nil {
		return err
	}
	if ipNet.IP.To4().String() != ip.To4().String() {
		return common.NewApiError(http.StatusBadRequest, errors.Errorf("%s is not a valid machine CIDR", machineCidr))
	}
	for _, h := range hosts {
		if belongsToNetwork(log, h, ipNet) {
			return nil
		}
	}
	return common.NewApiError(http.StatusBadRequest, errors.Errorf("%s does not belong to any of the host networks", machineCidr))
}

func findMatchingIPForFamily(ipnet *net.IPNet, addresses []string) (bool, string) {
	for _, addr := range addresses {
		ip, _, err := net.ParseCIDR(addr)
		if err != nil {
			continue
		}
		if ipnet.Contains(ip) {
			return true, addr
		}
	}
	return false, ""
}

func findMatchingIP(ipnet *net.IPNet, intf *models.Interface, isIPv4 bool) (bool, string) {
	if isIPv4 {
		return findMatchingIPForFamily(ipnet, intf.IPV4Addresses)
	} else {
		return findMatchingIPForFamily(ipnet, intf.IPV6Addresses)
	}
}

func getMachineCIDRObj(host *models.Host, machineNetworkCidr string, obj string) (string, error) {
	var inventory models.Inventory
	var err error
	isIPv4 := IsIPV4CIDR(machineNetworkCidr)
	if err = json.Unmarshal([]byte(host.Inventory), &inventory); err != nil {
		return "", err
	}
	_, ipNet, err := net.ParseCIDR(machineNetworkCidr)
	if err != nil {
		return "", err
	}
	for _, intf := range inventory.Interfaces {
		found, addr := findMatchingIP(ipNet, intf, isIPv4)
		if found {
			switch obj {
			case "interface":
				return intf.Name, nil
			case "ip":
				return strings.Split(addr, "/")[0], nil
			default:
				return "", errors.Errorf("obj %s not supported", obj)
			}
		}
	}
	return "", errors.Errorf("No matching interface found for host %s", host.ID.String())
}

func GetMachineCIDRInterface(host *models.Host, cluster *common.Cluster) (string, error) {
	return getMachineCIDRObj(host, cluster.MachineNetworkCidr, "interface")
}

func GetMachineCIDRIP(host *models.Host, cluster *common.Cluster) (string, error) {
	return getMachineCIDRObj(host, cluster.MachineNetworkCidr, "ip")
}

func IpInCidr(ipAddr, cidr string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return false, errors.New("IP is nil")
	}
	return ipNet.Contains(ip), nil
}

func belongsToNetwork(log logrus.FieldLogger, h *models.Host, machineIpnet *net.IPNet) bool {
	var inventory models.Inventory
	err := json.Unmarshal([]byte(h.Inventory), &inventory)
	if err != nil {
		log.WithError(err).Warnf("Error unmarshalling host %s inventory %s", h.ID, h.Inventory)
		return false
	}
	isIPv4 := IsIPV4CIDR(machineIpnet.String())
	for _, intf := range inventory.Interfaces {
		if found, _ := findMatchingIP(machineIpnet, intf, isIPv4); found {
			return true
		}
	}
	return false
}

func GetMachineCIDRHosts(log logrus.FieldLogger, cluster *common.Cluster) ([]*models.Host, error) {
	if cluster.MachineNetworkCidr == "" {
		return nil, errors.New("Machine network CIDR was not set in cluster")
	}
	_, machineIpnet, err := net.ParseCIDR(cluster.MachineNetworkCidr)
	if err != nil {
		return nil, err
	}
	ret := make([]*models.Host, 0)
	for _, h := range cluster.Hosts {
		if belongsToNetwork(log, h, machineIpnet) {
			ret = append(ret, h)
		}
	}
	return ret, nil
}

// GetMachineCidrForUserManagedNetwork used to get machine cidr in case of none platform and sno
func GetMachineCidrForUserManagedNetwork(cluster *common.Cluster, log logrus.FieldLogger) string {
	if cluster.MachineNetworkCidr != "" {
		return cluster.MachineNetworkCidr
	}

	bootstrap := common.GetBootstrapHost(cluster)
	if bootstrap == nil {
		return ""
	}

	networks := GetClusterNetworks([]*models.Host{bootstrap}, log)
	if len(networks) > 0 {
		// if there is ipv4 network, return it or return the first one
		for _, network := range networks {
			if IsIPV4CIDR(network) {
				return network
			}
		}
		return networks[0]
	}
	return ""
}

func GetIpForSingleNodeInstallation(cluster *common.Cluster, log logrus.FieldLogger) (string, error) {
	bootstrap := common.GetBootstrapHost(cluster)
	if bootstrap == nil {
		return "", errors.Errorf("no bootstrap host were found in cluster")
	}
	cidr := GetMachineCidrForUserManagedNetwork(cluster, log)
	hostIp, err := getMachineCIDRObj(bootstrap, cidr, "ip")
	if hostIp == "" || err != nil {
		msg := "failed to get ip for bootstrap in place dnsmasq manifest"
		if err != nil {
			msg = errors.Wrapf(err, msg).Error()
		}
		return "", errors.Errorf(msg)
	}

	return hostIp, nil
}

func GetClusterNetworks(hosts []*models.Host, log logrus.FieldLogger) []string {
	var err error
	cidrs := make(map[string]bool)
	for _, h := range hosts {
		if h.Inventory != "" {
			var inventory models.Inventory
			err = json.Unmarshal([]byte(h.Inventory), &inventory)
			if err != nil {
				log.WithError(err).Warnf("Unmarshal inventory %s", h.Inventory)
				continue
			}
			for _, inf := range inventory.Interfaces {

				for _, ipv4 := range inf.IPV4Addresses {
					_, cidr, err := net.ParseCIDR(ipv4)
					if err != nil {
						log.WithError(err).Warnf("Parse CIDR %s", ipv4)
						continue
					}
					cidrs[cidr.String()] = true
				}

				for _, ipv6 := range inf.IPV6Addresses {
					_, cidr, err := net.ParseCIDR(ipv6)
					if err != nil {
						log.WithError(err).Warnf("Parse CIDR %s", ipv6)
						continue
					}
					cidrs[cidr.String()] = true
				}
			}
		}
	}
	ret := make([]string, 0)
	for cidr := range cidrs {
		ret = append(ret, cidr)
	}
	return ret
}

func IsHostInMachineNetCidr(log logrus.FieldLogger, cluster *common.Cluster, host *models.Host) bool {
	_, machineIpnet, err := net.ParseCIDR(cluster.MachineNetworkCidr)
	if err != nil {
		return false
	}
	return belongsToNetwork(log, host, machineIpnet)
}

type IPSet map[strfmt.IPv4]struct{}

func (s IPSet) Add(str strfmt.IPv4) {
	s[str] = struct{}{}
}

func (s IPSet) Intersect(other IPSet) IPSet {
	ret := make(IPSet)
	for k := range s {
		if v, ok := other[k]; ok {
			ret[k] = v
		}
	}
	return ret
}

func freeAddressesUnmarshal(network, freeAddressesStr string, prefix *string) (IPSet, error) {
	var unmarshaled models.FreeNetworksAddresses
	err := json.Unmarshal([]byte(freeAddressesStr), &unmarshaled)
	if err != nil {
		return nil, err
	}
	for _, f := range unmarshaled {
		if f.Network == network {
			ret := make(IPSet)
			for _, a := range f.FreeAddresses {
				if prefix == nil || strings.HasPrefix(a.String(), *prefix) {
					ret.Add(a)
				}
			}
			return ret, nil
		}
	}
	return nil, errors.Errorf("No network %s found", network)
}

func MakeFreeAddressesSet(hosts []*models.Host, network string, prefix *string, log logrus.FieldLogger) IPSet {
	var (
		availableFreeAddresses []string
		sets                   = make([]IPSet, 0)
		resultingSet           = make(IPSet)
	)
	for _, h := range hosts {
		if swag.StringValue(h.Status) != models.HostStatusDisabled && h.FreeAddresses != "" {
			availableFreeAddresses = append(availableFreeAddresses, h.FreeAddresses)
		}
	}
	if len(availableFreeAddresses) == 0 {
		return resultingSet
	}
	// Create IP sets from each of the hosts free-addresses
	for _, a := range availableFreeAddresses {
		s, err := freeAddressesUnmarshal(network, a, prefix)
		if err != nil {
			log.WithError(err).Debugf("Unmarshal free addresses for network %s", network)
			continue
		}
		// TODO: Have to decide if we want to filter empty sets
		sets = append(sets, s)
	}
	if len(sets) == 0 {
		return resultingSet
	}

	// Perform set intersection between all valid sets
	resultingSet = sets[0]
	for _, s := range sets[1:] {
		resultingSet = resultingSet.Intersect(s)
	}
	return resultingSet
}

// This is best effort validation.  Therefore, validation will be done only if there are IPs in free list
func IpInFreeList(hosts []*models.Host, vipIPStr, network string, log logrus.FieldLogger) bool {
	isFree := true
	freeSet := MakeFreeAddressesSet(hosts, network, nil, log)
	if len(freeSet) > 0 {
		_, isFree = freeSet[strfmt.IPv4(vipIPStr)]
	}
	return isFree
}

func CreateIpWithCidr(ip, cidr string) (string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	ipNet.IP = net.ParseIP(ip)
	return ipNet.String(), nil
}
