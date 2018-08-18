package zevenetlb

import "strings"

//
// Network Interfaces
//

type nicListResponse struct {
	Description string    `json:"description"`
	Interfaces  []NICInfo `json:"interfaces"`
}

// NICInfo contains the list of all available NICs.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-nic-interfaces
type NICInfo struct {
	IP      string `json:"ip"`
	HasVlan string `json:"has_vlan"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetAllNICs returns list os all available NICs.
func (s *ZapiSession) GetAllNICs() ([]NICInfo, error) {
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

type virtIntListResponse struct {
	Description string        `json:"description"`
	Interfaces  []VirtIntInfo `json:"interfaces"`
}

// VirtIntInfo contains the list of all available virtual Interfaces.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-virtual-interfaces
type VirtIntInfo struct {
	IP      string `json:"ip"`
	Parent  string `json:"parent"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetAllVirtInts returns list os all available NICs.
func (s *ZapiSession) GetAllVirtInts() ([]VirtIntInfo, error) {
	var result *virtIntListResponse

	err := s.getForEntity(&result, "interfaces", "virtual")

	if err != nil {
		return nil, err
	}

	return result.Interfaces, nil
}

type virtIntDetailsResponse struct {
	Description string         `json:"description"`
	Interface   VirtIntDetails `json:"interface"`
}

// VirtIntDetails contains all information regarding a virtual Interface.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#retrieve-virtual-interface
type VirtIntDetails struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	MAC     string `json:"mac"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// GetVirtInt returns details on a specific virtual Interface.
func (s *ZapiSession) GetVirtInt(virtIntName string) (*VirtIntDetails, error) {
	var result *virtIntDetailsResponse

	err := s.getForEntity(&result, "interfaces", "virtual", virtIntName)

	if err != nil {
		// virtInt not found?
		if strings.Contains(err.Error(), "VirtInt not found") {
			return nil, nil
		}

		return nil, err
	}

	return &result.Interface, nil
}

// DeleteVirtInt will delete an existing virtual Interface (or do nothing if missing)
func (s *ZapiSession) DeleteVirtInt(virtIntName string) (bool, error) {
	// retrieve virtInt details
	virtInt, err := s.GetVirtInt(virtIntName)

	if err != nil {
		return false, err
	}

	// farm does not exist?
	if virtInt == nil {
		return false, nil
	}

	// delete the farm
	return true, s.delete("interfaces", "virtual", virtIntName)
}

type virtIntCreate struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// CreateVirtInt creates a new virtual Interface.
func (s *ZapiSession) CreateVirtInt(virtIntName string, virtualIP string) (*VirtIntDetails, error) {

	req := virtIntCreate{
		IP:   virtualIP,
		Name: virtIntName,
	}

	err := s.post(req, "interfaces", "virtual")

	if err != nil {
		return nil, err
	}

	// retrieve status
	return s.GetVirtInt(virtIntName)
}
