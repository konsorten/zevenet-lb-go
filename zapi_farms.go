package zevenetlb

import (
	"fmt"
	"strings"
)

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

// GetAllFarms returns list os all available farms.
func (b *ZapiSession) GetAllFarms() ([]FarmInfo, error) {
	var result *farmListResponse

	err, _ := b.getForEntity(&result, "farms")

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
	CiphersCustom            string            `json:"cipherc"`
	Ciphers                  string            `json:"ciphers"`
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
func (sv *FarmDetails) String() string {
	return fmt.Sprintf("%v (%v)", sv.FarmName, sv.Listener)
}

// IsHTTP checks if the farm has HTTP or HTTPS support enabled.
func (sv *FarmDetails) IsHTTP() bool {
	return strings.HasPrefix(sv.Listener, "http")
}

// GetFarmDetails returns details on a specific farm.
func (b *ZapiSession) GetFarmDetails(farmName string) (*FarmDetails, error) {
	var result *farmDetailsResponse

	err, _ := b.getForEntity(&result, "farms", farmName)

	if err != nil {
		return nil, err
	}

	// inject values
	if result != nil {
		result.Params.FarmName = farmName
		result.Params.Services = result.Services
	}

	return &result.Params, nil
}

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
