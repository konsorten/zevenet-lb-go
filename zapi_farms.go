package zevenetlb

import (
	"fmt"
	"strings"
)

type genericResponse struct {
	Description string `json:"description"`
}

type farmListResponse struct {
	Description string     `json:"description"`
	Params      []FarmInfo `json:"params"`
}

// FarmInfo contains the list of all available farms.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#list-all-farms
type FarmInfo struct {
	FarmName    string `json:"farmname"`
	Profile     string `json:"profile"`
	Status      string `json:"status"`
	VirtualIP   string `json:"vip"`
	VirtualPort int    `json:"vport,string"`
}

// String returns the farm's name and profile.
func (fi *FarmInfo) String() string {
	return fmt.Sprintf("%v (%v)", fi.FarmName, fi.Profile)
}

// GetAllFarms returns list os all available farms.
func (s *ZapiSession) GetAllFarms() ([]FarmInfo, error) {
	var result *farmListResponse

	err := s.getForEntity(&result, "farms")

	if err != nil {
		return nil, err
	}

	return result.Params, nil
}

type farmDetailsResponse struct {
	Description string           `json:"description"`
	Params      FarmDetails      `json:"params"`
	Services    []ServiceDetails `json:"services"`
}

// FarmDetails contains all information regarding a farm and the services.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#retrieve-farm-by-name
type FarmDetails struct {
	Certificates             []CertificateInfo `json:"certlist"`
	FarmName                 string            `json:"farmname"`
	CiphersCustom            string            `json:"cipherc,omitempty"`
	Ciphers                  string            `json:"ciphers,omitempty"`
	ConnectionTimeoutSeconds int               `json:"contimeout"`
	DisableSSLv2             bool              `json:"disable_sslv2,string"`
	DisableSSLv3             bool              `json:"disable_sslv3,string"`
	DisableTLSv1             bool              `json:"disable_tlsv1,string"`
	DisableTLSv11            bool              `json:"disable_tlsv1_1,string"`
	DisableTLSv12            bool              `json:"disable_tlsv1_2,string"`
	ErrorString414           string            `json:"error414"`
	ErrorString500           string            `json:"error500"`
	ErrorString501           string            `json:"error501"`
	ErrorString503           string            `json:"error503"`
	HTTPVerbs                string            `json:"httpverb"`
	Listener                 string            `json:"listener"`
	RequestTimeoutSeconds    int               `json:"reqtimeout"`
	ResponseTimeoutSeconds   int               `json:"restimeout"`
	ResurrectIntervalSeconds int               `json:"resurrectime"`
	RewriteLocation          string            `json:"rewritelocation"`
	Status                   string            `json:"status"`
	VirtualIP                string            `json:"vip"`
	VirtualPort              int               `json:"vport"`
	Services                 []ServiceDetails  `json:"services"`
}

// String returns the farm's name and listener.
func (fd *FarmDetails) String() string {
	return fmt.Sprintf("%v (%v)", fd.FarmName, fd.Listener)
}

// IsHTTP checks if the farm has HTTP or HTTPS support enabled.
func (fd *FarmDetails) IsHTTP() bool {
	return strings.HasPrefix(fd.Listener, "http")
}

// IsRunning checks if the farm is up and running.
func (fd *FarmDetails) IsRunning() bool {
	return fd.Status == "up"
}

// IsRestartRequired checks if the farm needs to be restartet.
func (fd *FarmDetails) IsRestartRequired() bool {
	return fd.Status == "needed restart"
}

// GetFarm returns details on a specific farm.
func (s *ZapiSession) GetFarm(farmName string) (*FarmDetails, error) {
	var result *farmDetailsResponse

	err := s.getForEntity(&result, "farms", farmName)

	if err != nil {
		// farm not found?
		if strings.Contains(err.Error(), "Farm not found") {
			return nil, nil
		}

		return nil, err
	}

	// inject values
	if result != nil {
		result.Params.FarmName = farmName
		result.Params.Services = result.Services
	}

	return &result.Params, nil
}

// DeleteFarm will delete an existing farm (or do nothing if missing)
func (s *ZapiSession) DeleteFarm(farmName string) (bool, error) {
	// retrieve farm details
	farm, err := s.GetFarm(farmName)

	if err != nil {
		return false, err
	}

	// farm does not exist?
	if farm == nil {
		return false, nil
	}

	// delete the farm
	return true, s.delete("farms", farmName)
}

type farmCreate struct {
	FarmName    string `json:"farmname"`
	Profile     string `json:"profile"`
	VirtualIP   string `json:"vip"`
	VirtualPort int    `json:"vport"`
}

// CreateFarmAsHTTP creates a new HTTP farm.
func (s *ZapiSession) CreateFarmAsHTTP(farmName string, virtualIP string, virtualPort int) (*FarmDetails, error) {
	// set default HTTP port
	if virtualPort <= 0 {
		virtualPort = 80
	}

	// create the farm
	req := farmCreate{
		FarmName:    farmName,
		Profile:     "http",
		VirtualIP:   virtualIP,
		VirtualPort: virtualPort,
	}

	err := s.post(req, "farms")

	if err != nil {
		return nil, err
	}

	// retrieve status
	return s.GetFarm(farmName)
}

// CreateFarmAsHTTPS creates a new HTTPS farm.
func (s *ZapiSession) CreateFarmAsHTTPS(farmName string, virtualIP string, virtualPort int, certFilename string) (*FarmDetails, error) {
	// set default HTTPS port
	if virtualPort <= 0 {
		virtualPort = 443
	}

	// create the farm
	farm, err := s.CreateFarmAsHTTP(farmName, virtualIP, virtualPort)

	if err != nil {
		return nil, err
	}

	// update the farm
	farm.Listener = "https"
	farm.Ciphers = "highsecurity"
	farm.DisableSSLv2 = false
	farm.DisableSSLv3 = false
	farm.DisableTLSv1 = false

	s.UpdateFarm(farm)

	return farm, nil
}

// UpdateFarm updates the HTTP/S farm.
func (s *ZapiSession) UpdateFarm(farm *FarmDetails) error {
	return s.put(farm, "farms", farm.FarmName)
}

type farmAction struct {
	Action string `json:"action"`
}

// StartFarm will start a stopped farm.
func (s *ZapiSession) StartFarm(farmName string) error {
	req := farmAction{Action: "start"}

	return s.put(req, "farms", farmName, "actions")
}

// StopFarm will stop a running farm.
func (s *ZapiSession) StopFarm(farmName string) error {
	req := farmAction{Action: "stop"}

	return s.put(req, "farms", farmName, "actions")
}

// RestartFarm will restart a running farm.
func (s *ZapiSession) RestartFarm(farmName string) error {
	req := farmAction{Action: "restart"}

	return s.put(req, "farms", farmName, "actions")
}

// CertificateInfo contains reference information on a certificate.
type CertificateInfo struct {
	Filename string `json:"file"`
	ID       int    `json:"id"`
}

// String returns the certificate's filename.
func (ci CertificateInfo) String() string {
	return ci.Filename
}

type serviceDetailsResponse struct {
	Description string         `json:"description"`
	Params      ServiceDetails `json:"params"`
}

// ServiceDetails contains all information regarding a single service.
type ServiceDetails struct {
	ServiceName                         string           `json:"id"`
	FarmGuardianEnabled                 bool             `json:"fgenabled,string"`
	FarmGuardianLogsEnabled             string           `json:"fglog"`
	FarmGuardianScriptEnabled           string           `json:"fgscript"`
	FarmGuardianCheckIntervalSeconds    int              `json:"fgtimecheck"`
	EncryptedBackends                   bool             `json:"httpsb,string"`
	LastResponseBalancingEnabled        bool             `json:"leastresp,string"`
	ConnectionPersistenceMode           string           `json:"persistence"`
	ConnectionPersistenceID             string           `json:"sessionid"`
	ConnectionPersistenceTimeoutSeconds int              `json:"ttl"`
	RedirectURL                         string           `json:"redirect"`
	RedirectType                        string           `json:"redirecttype"`
	URLPattern                          string           `json:"urlp"`
	HostPattern                         string           `json:"vhost"`
	Backends                            []BackendDetails `json:"backends"`
}

// String returns the services' name.
func (sd ServiceDetails) String() string {
	return sd.ServiceName
}

type backendDetailsResponse struct {
	Description string         `json:"description"`
	Params      BackendDetails `json:"params"`
}

// BackendDetails contains all information regarding a single backend server.
type BackendDetails struct {
	ID             int    `json:"id"`
	IPAddress      string `json:"ip"`
	Port           int    `json:"port"`
	Status         string `json:"status"`
	TimeoutSeconds *int   `json:"timeout"`
	Weight         *int   `json:"weight"`
}

// String returns the backend's IP, port, and ID.
func (bd BackendDetails) String() string {
	return fmt.Sprintf("%v:%v (ID: %v)", bd.IPAddress, bd.Port, bd.ID)
}
