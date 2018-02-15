package zevenetlb

import "fmt"

// This is how to connect to the Zevenet loadbalancer.
func ExampleConnect() {
	session := Connect("myloadbalancer:444", "zapi-key", nil)

	version, _ := session.GetSystemVersion()

	fmt.Println(version)
	// Output: ZCE 5 (v5.0)
}

// This is how to retrieve a specific farm. In this case the first farm that exists.
func ExampleZapiSession_GetFarm() {
	session := Connect("myloadbalancer:444", "zapi-key", nil)

	farms, _ := session.GetAllFarms()

	farm, _ := session.GetFarm(farms[0].FarmName)

	fmt.Println(farm)
}

// This is how to create a new HTTP farm *without* SSL support.
func ExampleZapiSession_CreateFarmAsHTTP() {
	session := Connect("myloadbalancer:444", "zapi-key", nil)

	farm, _ := session.CreateFarmAsHTTP("mynewfarm", "10.10.10.10", 80)

	fmt.Println(farm)
	// Output: mynewfarm (http)
}

// This is how to create a new HTTP farm *with* SSL support, using the Zevenet default certificate.
func ExampleZapiSession_CreateFarmAsHTTPS() {
	session := Connect("myloadbalancer:444", "zapi-key", nil)

	farm, _ := session.CreateFarmAsHTTPS("mynewfarm", "10.10.10.10", 443, "zencert.pem")

	fmt.Println(farm)
	// Output: mynewfarm (https)
}