package transip

import (
	"fmt"
	"os"
	"regexp"
	"testing"

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
	testNameServers = append(testNameServers, "ns01.test.example")
	testNameServers = append(testNameServers, "ns02.test.example")
	testNameServers = append(testNameServers, "ns03.test.example")

	defaultNameServers := make([]string, 0)
	defaultNameServers = append(defaultNameServers, "ns0.transip.net")
	defaultNameServers = append(defaultNameServers, "ns1.transip.nl")
	defaultNameServers = append(defaultNameServers, "ns2.transip.eu")

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
			{
				Config: testAccResourceDomainConfigWithNS(domainName, true, defaultNameServers),
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
