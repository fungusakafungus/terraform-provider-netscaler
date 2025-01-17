/*
Copyright 2016 Citrix Systems, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package netscaler

import (
	"fmt"
	"github.com/chiradeep/go-nitro/netscaler"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccServicegroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServicegroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServicegroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServicegroupExist("netscaler_servicegroup.foo", nil),

					resource.TestCheckResourceAttr(
						"netscaler_servicegroup.foo", "servicegroupname", "test_servicegroup"),
					resource.TestCheckResourceAttr(
						"netscaler_servicegroup.foo", "servicetype", "HTTP"),
				),
			},
		},
	})
}

func testAccCheckServicegroupExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No lb vserver name is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client
		data, err := nsClient.FindResource(netscaler.Servicegroup.Type(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if data == nil {
			return fmt.Errorf("LB vserver %s not found", n)
		}

		return nil
	}
}

func testAccCheckServicegroupDestroy(s *terraform.State) error {
	nsClient := testAccProvider.Meta().(*NetScalerNitroClient).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netscaler_servicegroup" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No name is set")
		}

		_, err := nsClient.FindResource(netscaler.Servicegroup.Type(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("LB vserver %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

// TODO add testcase when we have servicegroupmembers_by_servername defined
const testAccServicegroup_basic = `

resource "netscaler_lbvserver" "foo1" {
  
  name = "foo_lb_1"
  ipv46 = "10.202.11.11"
  port = 80
  servicetype = "HTTP"
}

resource "netscaler_lbvserver" "foo2" {
  
  name = "foo_lb_2"
  ipv46 = "10.202.11.12"
  port = 80
  servicetype = "HTTP"
}


resource "netscaler_servicegroup" "foo" {
  
  servicegroupname = "test_servicegroup"
  servicetype = "HTTP"
  servicegroupmembers = ["172.20.0.9:80:10", "172.20.0.10:80:10", "172.20.0.11:8080:20"]
  lbvservers = ["foo_lb_1", "foo_lb_2"]
  depends_on = ["netscaler_lbvserver.foo1", "netscaler_lbvserver.foo2"]
}
`
