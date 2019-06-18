package err_fmt

import (
	"regexp"
	"testing"

	"github.com/transip/gotransip/domain"
)

func TestCreateDomainUnavailableError(t *testing.T) {
	testCases := []struct {
		desc           string
		domainStatus   domain.Status
		domainName     string
		errorMsgRegExp string
	}{
		{
			desc:           "inyouraccount",
			domainStatus:   domain.StatusInYourAccount,
			domainName:     "test.local",
			errorMsgRegExp: `Domain \[test.local\] is not available for registration: Domain is already in your account.*`,
		},
		{
			desc:           "unavailable",
			domainStatus:   domain.StatusUnavailable,
			domainName:     "test.local",
			errorMsgRegExp: `Domain \[test.local\] is not available for registration: Domain is currently unavailable.*`,
		},
		{
			desc:           "notfree",
			domainStatus:   domain.StatusNotFree,
			domainName:     "test.local",
			errorMsgRegExp: `Domain \[test.local\] is not available for registration: Domain has already been registered.*`,
		},
		{
			desc:           "internalpull",
			domainStatus:   domain.StatusInternalPull,
			domainName:     "test.local",
			errorMsgRegExp: `Domain \[test.local\] is not available for registration: Domain is available to be pulled.*`,
		},
		{
			desc:           "internalpush",
			domainStatus:   domain.StatusInternalPush,
			domainName:     "test.local",
			errorMsgRegExp: `Domain \[test.local\] is not available for registration: Domain is available to be pushed.*`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			//logic
			err := CreateDomainUnavailableError("test.local", tC.domainStatus)

			// test errorstring
			matched, _ := regexp.MatchString(tC.errorMsgRegExp, err.Error())
			if matched == false {
				t.Errorf("expected a match with '%s' got '%s'", tC.errorMsgRegExp, err.Error())
			}
		})
	}
}
