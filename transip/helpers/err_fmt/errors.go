package err_fmt

import (
	"fmt"

	"github.com/transip/gotransip/domain"
)

//CreateDomainUnavailableError returns an error specifying why the domain is not available for registration
func CreateDomainUnavailableError(name string, status domain.Status) error {
	msg := "Domain is currently unavailable and can not be registered due to unknown reasons."

	switch status {
	case domain.StatusInYourAccount:
		msg = "Domain is already in your account. run `terraform import ADDR ID` to mange the domain."
	case domain.StatusNotFree:
		msg = "Domain has already been registered and is not available."
	case domain.StatusInternalPull:
		msg = "Domain is available to be pulled from another account to yours."
	case domain.StatusInternalPush:
		msg = "Domain is available to be pushed to another account from yours."
	}

	return fmt.Errorf("Domain [%s] is not available for registration: %s", name, msg)
}
