package zevenetlb

import (
	"strings"
)

//
// Network Interfaces
//

type nicListResponse struct {
	Description string                 `json:"description"`
	Interfaces  []NetworkInterfaceInfo `json:"interfaces"`
}

// NetworkInterfaceInfo contains the list of all available NICs.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-nic-interfaces
type NetworkInterfaceInfo struct {
	IP      string `json:"ip"`
	HasVlan string `json:"has_vlan"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetAllNetworkInterfaces returns list os all available NICs.
func (s *ZapiSession) GetAllNetworkInterfaces() ([]NetworkInterfaceInfo, error) {
	var result *nicListResponse

	err := s.getForEntity(&result, "interfaces", "nic")

	if err != nil {
		return nil, err
	}

	return result.Interfaces, nil
}

//
// Virtual Interfaces
//

type virtualInterfaceListResponse struct {
	Description string                 `json:"description"`
	Interfaces  []VirtualInterfaceInfo `json:"interfaces"`
}

// VirtualInterfaceInfo contains the list of all available virtual Interfaces.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-virtual-interfaces
type VirtualInterfaceInfo struct {
	IP      string `json:"ip"`
	Parent  string `json:"parent"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetAllVirtualInterfaces returns list os all available NICs.
func (s *ZapiSession) GetAllVirtualInterfaces() ([]VirtualInterfaceInfo, error) {
	var result *virtualInterfaceListResponse

	err := s.getForEntity(&result, "interfaces", "virtual")

	if err != nil {
		return nil, err
	}

	return result.Interfaces, nil
}

type virtualInterfaceDetailsResponse struct {
	Description string                  `json:"description"`
	Interface   VirtualInterfaceDetails `json:"interface"`
}

// VirtualInterfaceDetails contains all information regarding a virtual Interface.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#retrieve-virtual-interface
type VirtualInterfaceDetails struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetVirtualInterface returns details on a specific virtual Interface.
func (s *ZapiSession) GetVirtualInterface(virtualInterfaceName string) (*VirtualInterfaceDetails, error) {
	var result *virtualInterfaceDetailsResponse

	err := s.getForEntity(&result, "interfaces", "virtual", virtualInterfaceName)

	if err != nil {
		// virtualInterface not found?
		if v, ok := err.(RequestError); ok {
			if strings.Contains(v.Message, "not found") {
				return nil, nil
			}
		}

		return nil, err
	}

	return &result.Interface, nil
}

// DeleteVirtualInterface will delete an existing virtual Interface (or do nothing if missing)
func (s *ZapiSession) DeleteVirtualInterface(virtualInterfaceName string) (bool, error) {
	// retrieve virtualInterface details
	virtualInterface, err := s.GetVirtualInterface(virtualInterfaceName)

	if err != nil {
		return false, err
	}

	// farm does not exist?
	if virtualInterface == nil {
		return false, nil
	}

	// delete the farm
	return true, s.delete("interfaces", "virtual", virtualInterfaceName)
}

type virtualInterfaceCreate struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// CreateVirtualInterface creates a new virtual Interface.
func (s *ZapiSession) CreateVirtualInterface(virtualInterfaceName string, virtualIP string) (*VirtualInterfaceDetails, error) {

	req := virtualInterfaceCreate{
		IP:   virtualIP,
		Name: virtualInterfaceName,
	}

	err := s.post(req, "interfaces", "virtual")

	if err != nil {
		return nil, err
	}

	// retrieve status
	return s.GetVirtualInterface(virtualInterfaceName)
}
