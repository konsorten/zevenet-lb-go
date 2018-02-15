package zevenetlb

import "fmt"

// ExampleConnect shows how to connect to the Zevenet loadbalancer.
func ExampleConnect() {
	session := Connect("myloadbalancer:444", "zapi-key", nil)

	version, _ := session.GetSystemVersion()

	fmt.Println(version)
}
