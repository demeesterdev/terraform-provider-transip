package transip

import (
	"github.com/transip/gotransip/domain"
)

func expandNameServers(in []interface{}) domain.Nameservers {

	out := make(domain.Nameservers, 0)
	for _, vRaw := range in {
		v := vRaw.(string)
		entry := domain.Nameserver{
			Hostname: v,
		}
		out = append(out, entry)
	}
	return out
}

func flattenNameServers(in domain.Nameservers) []string {
	out := make([]string, 0)
	for _, ns := range in {
		out = append(out, ns.Hostname)
	}
	return out
}
