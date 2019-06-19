package transip

import (
	"testing"

	"github.com/transip/gotransip/domain"
)

func TestExpandNameServers(t *testing.T) {
	testCases := []struct {
		desc        string
		nameServers []string
	}{
		{
			desc:        "empty",
			nameServers: []string{},
		},
		{
			desc:        "many",
			nameServers: []string{"ns01.test", "ns02.test"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			//logic
			nSRaw := make([]interface{}, len(tC.nameServers))
			for i, v := range tC.nameServers {
				nSRaw[i] = v
			}
			nameServerObjs := expandNameServers(nSRaw)

			//test length
			nsCount := len(tC.nameServers)
			nsObjCount := len(nameServerObjs)
			if nsCount != nsObjCount {
				t.Errorf("expected %d name servers objects got %d", nsCount, nsObjCount)
			}

			//test name comparison
			for i, nSO := range nameServerObjs {
				if nSO.Hostname != tC.nameServers[i] {
					t.Errorf("expected %s at index %d got %s", tC.nameServers[i], i, nSO.Hostname)
				}
			}
		})
	}
}

func TestFlattenNameServers(t *testing.T) {
	testCases := []struct {
		desc        string
		nameServers []string
	}{
		{
			desc:        "empty",
			nameServers: []string{},
		},
		{
			desc:        "many",
			nameServers: []string{"ns01.test", "ns02.test"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			//logic
			nameServerObjs := make([]domain.Nameserver, 0)
			for _, v := range tC.nameServers {
				nsObj := domain.Nameserver{
					Hostname: v,
				}
				nameServerObjs = append(nameServerObjs, nsObj)
			}
			nameServers := flattenNameServers(nameServerObjs)

			//test length
			nsCount := len(tC.nameServers)
			nsStrsCount := len(nameServers)
			if nsCount != nsStrsCount {
				t.Errorf("expected %d strings got %d", nsCount, nsStrsCount)
			}

			//test name comparison
			for i, nS := range nameServers {
				if nS != tC.nameServers[i] {
					t.Errorf("expected string %s at index %d got %s", tC.nameServers[i], i, nS)
				}
			}
		})
	}
}
