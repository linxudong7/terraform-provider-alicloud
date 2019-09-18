package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

var physicalConnectionId = os.Getenv("ALICLOUD_PHYSICAL_CONNECTION_ID")

// At present, the provider does not support creating physical connection resource, so you should create manually a physical connection
// by web console and set it by environment variable ALICLOUD_PHYSICAL_CONNECTION_ID before running the following test case.
func TestAccAlicloudVirtualBorderRouter_basic(t *testing.T) {
	var v vpc.VirtualBorderRouterType
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d", rand)
	var basicMap = map[string]string{
		"physical_connection_id": physicalConnectionId,
		"vlan_id":                "2500",
		"local_gateway_ip":       "10.0.0.1",
		"peer_gateway_ip":        "10.0.0.2",
		"peering_subnet_mask":    "255.255.255.0",
		"name":                   name,
		"description":            "Acc vbr test",
	}
	resourceId := "alicloud_virtual_border_router.default"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVirtualBorderRouterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id": physicalConnectionId,
					"vlan_id":                "2500",
					"local_gateway_ip":       "10.0.0.1",
					"peer_gateway_ip":        "10.0.0.2",
					"peering_subnet_mask":    "255.255.255.0",
					"name":                   name,
					"description":            "Acc vbr test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id": physicalConnectionId,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_id": physicalConnectionId,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vlan_id": "2501",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vlan_id": "2501",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_gateway_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"local_gateway_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_gateway_ip": "10.0.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_gateway_ip": "10.0.0.4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peering_subnet_mask": "255.255.0.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peering_subnet_mask": "255.255.0.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_change", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d_change", rand),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Acc vbr test change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Acc vbr test change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_id": physicalConnectionId,
					"vlan_id":                "2500",
					"local_gateway_ip":       "10.0.0.1",
					"peer_gateway_ip":        "10.0.0.2",
					"peering_subnet_mask":    "255.255.255.0",
					"name":                   name,
					"description":            "Acc vbr test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(basicMap),
				),
			},
		},
	})
}

func TestAccAlicloudVirtualBorderRouter_multi(t *testing.T) {
	var v vpc.VirtualBorderRouterType
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccVirtualBorderRouterBasic%d", rand)
	var basicMap = map[string]string{
		"physical_connection_id": physicalConnectionId,
		"vlan_id":                "2501",
		"local_gateway_ip":       "10.0.0.3",
		"peer_gateway_ip":        "10.0.0.4",
		"peering_subnet_mask":    "255.255.255.0",
		"name":                   name,
		"description":            "Acc vbr test",
	}
	resourceId := "alicloud_virtual_border_router.default.1"
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVirtualBorderRouterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                  "2",
					"physical_connection_id": physicalConnectionId,
					"vlan_id":                "${element(var.vlan_id_list,count.index)}",
					"local_gateway_ip":       "${element(var.local_gateway_ip_list,count.index)}",
					"peer_gateway_ip":        "${element(var.peer_gateway_ip_list,count.index)}",
					"peering_subnet_mask":    "255.255.255.0",
					"name":                   name,
					"description":            "Acc vbr test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceVirtualBorderRouterConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	variable "vlan_id_list" {
		type = "list"
		default = [ "2500", "2501" ]
	}
	variable "local_gateway_ip_list" {
		type = "list"
		default = [ "10.0.0.1", "10.0.0.3" ]
	}
	variable "peer_gateway_ip_list" {
		type = "list"
		default = [ "10.0.0.2", "10.0.0.4" ]
	}
	`, name)
}
