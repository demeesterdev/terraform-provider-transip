package transip

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDomain(t *testing.T) {
	domainName := os.Getenv("TRANSIP_TEST_DOMAIN")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.transip_domain.test_domain", "name", domainName),
				),
			},
			{
				Config:      testAccDataSourceDomainBasic("test.local"),
				ExpectError: regexp.MustCompile("Could not request domain info \\[test.local\\]: .*"),
			},
		},
	})
}

func testAccDataSourceDomainBasic(name string) string {
	return fmt.Sprintf(`
	data "transip_domain" "test_domain" {
	  name        = "%s"
	}
	`, name)
}
