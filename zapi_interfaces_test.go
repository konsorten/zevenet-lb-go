package zevenetlb

import (
	"testing"

	ping "github.com/sparrc/go-ping"
)

const (
	unitTestVirtIntName = "eth0:test"
	unitTestVirtIntIP   = "10.209.0.31"
)

func TestGetAllNICs(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetAllNICs()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) <= 0 {
		t.Fatal("No NICs returned")
	}

}

func TestRoundtripVirtInt(t *testing.T) {
	session := createTestSession(t)

	// ensure the virtInt does not exist
	_, err := session.DeleteVirtInt(unitTestVirtIntName)

	if err != nil {
		t.Fatal(err)
	}

	// create the new virtInt
	vint, err := session.CreateVirtInt(unitTestVirtIntName, unitTestVirtIntIP)

	if err != nil {
		t.Fatal(err)
	}

	defer session.DeleteVirtInt(vint.Name)

	t.Logf("New Int: %v, Status: %v", vint.Name, vint.Status)

	// try to connect
	pinger, err := ping.NewPinger(vint.IP)
	if err != nil {
		panic(err)
	}

	pinger.Count = 3
	pinger.Run() // blocks until finished

	if err != nil {
		t.Fatal(err)
	}
	// done, delete the virtInt
	deleted, err := session.DeleteVirtInt(vint.Name)

	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("Expected deleting the virtual Interface to succeed, but failed")
	}
}
