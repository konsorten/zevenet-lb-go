package zevenetlb

import (
	"testing"
)

func TestGetAllFarms(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetAllFarms()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) <= 0 {
		t.Fatal("No farms returned")
	}

	t.Logf("Farms: %v", res)

	// get the farm details
	for _, f := range res {
		farm, err := session.GetFarmDetails(f.FarmName)

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Farm: %v", farm)

		for _, c := range farm.Certificates {
			t.Logf("  Certificate: %v", c)
		}

		for _, s := range farm.Services {
			t.Logf("  Service: %v", s)

			for _, b := range s.Backends {
				t.Logf("    Backend: %v", b)
			}
		}
	}
}
