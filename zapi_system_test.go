package zevenetlb

import (
	"os"
	"strings"
	"testing"
)

func createTestSession(t *testing.T) *ZapiSession {
	return createTestSessionEx(t, "")
}

func createTestSessionEx(t *testing.T, apiKey string) *ZapiSession {
	// retrieve api key if undefined
	if apiKey == "" {
		apiKey = os.Getenv("ZAPI_KEY")

		if apiKey == "" {
			t.Fatal("Failed to retrieve ZAPI key from environment variable ZAPI_KEY")
		}
	}

	// retrieve host name
	host := os.Getenv("ZAPI_HOSTNAME")

	if host == "" {
		host = "lb002.konsorten.net:444"
	}

	// create the session
	return NewSession(host, apiKey, nil)
}

func TestInvalidApiKey(t *testing.T) {
	session := createTestSessionEx(t, "inval1dAp1K3y")

	_, err := session.GetSystemVersion()

	if err == nil {
		t.Fatal("Error expected")
	}

	if !strings.Contains(err.Error(), "Authorization required") {
		t.Fatalf("Wrong error message returned: %v", err)
	}
}

func TestGetSystemVersion(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetSystemVersion()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Version: %v", res)
}

func TestIsCommunityEdition(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetSystemVersion()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Is Community Edition: %v", res.IsCommunityEdition())
}
