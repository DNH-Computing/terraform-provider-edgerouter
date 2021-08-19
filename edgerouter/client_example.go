package edgerouter

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/DNH-Computing/terraform-provider-edgerouter/edgerouter/client"
)

type GetZonePolicyInput struct {
	Get *GetZonePolicyNode `json:"GET,omitempty"` // omitempty to not put anything in the JSON if the field is `nil`
}
type GetZonePolicyOutput struct {
	Get     *GetZonePolicyNode `json:"GET"`
	Success bool               `json:"success"`
}

type GetZonePolicyNode struct {
	ZonePolicy *ZonePolicyNode `json:"zone-policy"`
}
type ZonePolicyNode struct {
	Zone map[string]*ZoneNode `json:"zone"`
}
type ZoneNode struct {
	DefaultAction string                   `json:"default-action"`
	From          map[string]*ZoneNodeFrom `json:"from"`
	Interface     []string                 `json:"interface"`
}
type ZoneNodeFrom struct {
	Firewall ZoneNodeFromFirewall `json:"firewall"`
}
type ZoneNodeFromFirewall struct {
	Name     string `json:"name"`
	Ipv6Name string `json:"ipv6-name"`
}

func getZonePolicy(client *client.Client, zone string) (*ZoneNode, error) {
	getSingleInput := GetZonePolicyInput{
		Get: &GetZonePolicyNode{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					zone: nil,
				},
			},
		},
	}
	var getSingleOutput GetZonePolicyOutput
	err := client.Post(context.Background(), "/api/edge/batch.json", &getSingleInput, &getSingleOutput)
	if err != nil {
		return nil, err
	}
	if !getSingleOutput.Success {
		return nil, fmt.Errorf("unsuccessful response: %+v", getSingleOutput)
	}
	return getSingleOutput.Get.ZonePolicy.Zone[zone], nil
}

func clientExample() {
	// Build a new client and log in
	client, err := client.NewClient(&tls.Config{
		// Accept expired/non-trusted/invalid CN certificates
		InsecureSkipVerify: true,
	}, "https://192.168.1.1", "ubnt", "ubnt")
	if err != nil {
		panic(err)
	}

	// Get all zone policies like `show zone-policy` on the CLI
	getAllInput := GetZonePolicyInput{
		Get: &GetZonePolicyNode{
			ZonePolicy: nil,
		},
	}
	var getAllOutput GetZonePolicyOutput
	err = client.Post(context.Background(), "/api/edge/batch.json", &getAllInput, &getAllOutput)
	if err != nil {
		panic(err)
	}
	if !getAllOutput.Success {
		panic(fmt.Errorf("unsuccessful response"))
	}
	fmt.Printf("%+v\n", getAllOutput)
	fmt.Printf("%+v\n", getAllOutput.Get.ZonePolicy.Zone["nhtest"])
	fmt.Printf("%+v\n", getAllOutput.Get.ZonePolicy.Zone["nhtest"].From["nhtest2"])
	fmt.Printf("%+v\n", getAllOutput.Get.ZonePolicy.Zone["nhtest2"])
	fmt.Printf("%+v\n", getAllOutput.Get.ZonePolicy.Zone["nhtest2"].From["nhtest"])

	// Get a single zone policy like `show zone-policy zone nhtest2` on the CLI
	getSingleInput := GetZonePolicyInput{
		Get: &GetZonePolicyNode{
			ZonePolicy: &ZonePolicyNode{
				Zone: map[string]*ZoneNode{
					"nhtest2": nil,
				},
			},
		},
	}
	var getSingleOutput GetZonePolicyOutput
	err = client.Post(context.Background(), "/api/edge/batch.json", &getSingleInput, &getSingleOutput)
	if err != nil {
		panic(err)
	}
	if !getSingleOutput.Success {
		panic(fmt.Errorf("unsuccessful response"))
	}
	fmt.Printf("%+v\n", getSingleOutput)
	fmt.Printf("%+v\n", getSingleOutput.Get.ZonePolicy.Zone["nhtest"]) // expect this one to be nil/absent, since it wasn't requested
	fmt.Printf("%+v\n", getSingleOutput.Get.ZonePolicy.Zone["nhtest2"])
	fmt.Printf("%+v\n", getSingleOutput.Get.ZonePolicy.Zone["nhtest2"].From["nhtest"])

	// Get a single zone policy using a helper method
	zone, err := getZonePolicy(client, "nhtest")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", zone)
	fmt.Printf("%+v\n", zone.From["nhtest2"])
}
