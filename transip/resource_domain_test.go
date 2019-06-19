package transip

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccResourceDomain_Basic_Invalid(t *testing.T) {
	managedDomainName := os.Getenv("TRANSIP_TEST_DOMAIN")
	invalidDomainName := "test.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceDomainConfig(managedDomainName, false),
				ExpectError: regexp.MustCompile(`Domain \[[a-z.]*\] is not available for registration: Domain is already in your account.*`),
			},
			{
				Config:      testAccResourceDomainConfig(invalidDomainName, false),
				ExpectError: regexp.MustCompile(`Domain \[[a-z.]*\] is not available for registration: Domain is currently unavailable.*`),
			},
		},
	})
}

func TestAccResourceDomain_Import(t *testing.T) {
	domainName := os.Getenv("TRANSIP_TEST_DOMAIN")
	resourceName := "transip_domain.test_domain"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: domainName,
			},
		},
	})
}

func TestAccResourceDomain_ManageSettings_Basic(t *testing.T) {
	domainName := os.Getenv("TRANSIP_TEST_DOMAIN")
	resourceName := "transip_domain.test_domain"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainConfig(domainName, true),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manage_settings_only"},
			},
		},
	})
}

func TestAccResourceDomain_ManageSettings_Basic_Invalid(t *testing.T) {
	invalidDomainName := "test.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceDomainConfig(invalidDomainName, true),
				ExpectError: regexp.MustCompile(`Error managing domain \[[a-z.]*\]: Domain is not in your account.*`),
			},
		},
	})
}

func TestAccResourceDomain_ManageSettings_NameServers(t *testing.T) {
	domainName := os.Getenv("TRANSIP_TEST_DOMAIN")
	resourceName := "transip_domain.test_domain"

	testNameServers := make([]string, 0)
	nsID := acctest.RandStringFromCharSet(4, "0123456789")

	testNameServers = append(testNameServers, fmt.Sprintf("ns1%s.test.example", nsID))
	testNameServers = append(testNameServers, fmt.Sprintf("ns2%s.test.example", nsID))
	testNameServers = append(testNameServers, fmt.Sprintf("ns3%s.test.example", nsID))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDomainConfigWithNS(domainName, true, testNameServers),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manage_settings_only"},
			},
		},
	})
}

func testAccResourceDomainConfig(name string, imp bool) string {
	return fmt.Sprintf(`
	resource "transip_domain" "test_domain" {
	  name                 = "%s"
	  manage_settings_only = %t
	}
	`, name, imp)
}

func testAccResourceDomainConfigWithNS(name string, imp bool, nS []string) string {
	var nSString string
	for _, v := range nS {
		formatString := `%s
		"%s",`
		if nSString == "" {
			formatString = `"%s%s",`
		}

		nSString = fmt.Sprintf(formatString, nSString, v)
	}
	return fmt.Sprintf(`
	resource "transip_domain" "test_domain" {
	  name                 = "%s"
	  manage_settings_only = %t
	  
	  name_servers = [
		%s
	  ]

	}
	`, name, imp, nSString)
}
