package zevenetlb

type nicListResponse struct {
	Description string    `json:"description"`
	Params      []NICInfo `json:"params"`
}

// NICInfo contains the list of all available farms.
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

	err := s.getForEntity(&result, "interfaces/nic")

	if err != nil {
		return nil, err
	}

	return result.Params, nil
}
