package zevenetlb

import (
	"runtime"
	"testing"

	ping "github.com/sparrc/go-ping"
)

const (
	unitTestVirtualInterfaceName = "eth0:unittest"
	unitTestVirtualInterfaceIP   = "10.209.0.31"
)

func TestGetAllNetworkInterfaces(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetAllNetworkInterfaces()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) <= 0 {
		t.Fatal("No NICs returned")
	}
}

func TestGetAllVirtualInterfaces(t *testing.T) {
	session := createTestSession(t)

	res, err := session.GetAllVirtualInterfaces()

	if err != nil {
		t.Fatal(err)
	}

	if len(res) <= 0 {
		t.Fatal("No NICs returned")
	}
}

func TestRoundtripVirtualInterface(t *testing.T) {
	session := createTestSession(t)

	// ensure the virtualInterface does not exist
	_, err := session.DeleteVirtualInterface(unitTestVirtualInterfaceName)

	if err != nil {
		t.Fatal(err)
	}

	// create the new virtualInterface
	vint, err := session.CreateVirtualInterface(unitTestVirtualInterfaceName, unitTestVirtualInterfaceIP)

	if err != nil {
		t.Fatal(err)
	}

	defer session.DeleteVirtualInterface(vint.Name)

	t.Logf("New Int: %v, Status: %v", vint.Name, vint.Status)

	// try to connect
	pinger, err := ping.NewPinger(vint.IP)
	if err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = 3
	pinger.Run() // blocks until finished

	if err != nil {
		t.Fatal(err)
	}
	// done, delete the virtualInterface
	deleted, err := session.DeleteVirtualInterface(vint.Name)

	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Fatal("Expected deleting the virtual Interface to succeed, but failed")
	}
}
