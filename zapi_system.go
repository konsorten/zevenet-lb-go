package zevenetlb

import (
	"fmt"
	"strings"
)

// SystemVersion contains information about the system version.
// See https://www.zevenet.com/zapidoc_ce_v3.1/#show-version
type SystemVersion struct {
	Description string `json:"description,omitempty"`
	Params      struct {
		ApplianceVersion string `json:"appliance_version,omitempty"`
		Hostname         string `json:"hostname,omitempty"`
		KernelVersion    string `json:"kernel_version,omitempty"`
		SystemDate       string `json:"system_date,omitempty"`
		ZevenetVersion   string `json:"zevenet_version,omitempty"`
	} `json:"params,omitempty"`
}

// String returns the version number of the system, e.g. "ZCE 5 (v5.0)"
func (sv *SystemVersion) String() string {
	return fmt.Sprintf("%v (v%v)", sv.Params.ApplianceVersion, sv.Params.ZevenetVersion)
}

// IsCommunityEdition checks if the Zevenet loadbalancer is the Community Edition (vs Enterprise Edition)
func (sv *SystemVersion) IsCommunityEdition() bool {
	return strings.HasPrefix(sv.Params.ApplianceVersion, "ZCE")
}

// GetSystemVersion returns system version information.
func (b *ZapiSession) GetSystemVersion() (result *SystemVersion, err error) {
	err, _ = b.getForEntity(&result, "system", "version")
	return
}
