package transip

import (
	"math/rand"
	"strings"
	"time"

	"github.com/transip/gotransip"
)

var DefaultRetrySleep = 200 * time.Millisecond
var DefaultRetryJitter = 40 * time.Millisecond
var DefaultRetryAttempts = 5

type RetryClient struct {
	soapClient    gotransip.SOAPClient
	retrySleep    time.Duration
	retryJitter   time.Duration
	retryAttempts int
}

func NewRetryClient(c gotransip.ClientConfig) (RetryClient, error) {
	soapClient, err := gotransip.NewSOAPClient(c)
	if err != nil {
		return RetryClient{}, err
	}

	return RetryClient{
		soapClient:    soapClient,
		retrySleep:    DefaultRetrySleep,
		retryJitter:   DefaultRetryJitter,
		retryAttempts: DefaultRetryAttempts,
	}, nil
}

func (c RetryClient) Call(req gotransip.SoapRequest, result interface{}) error {
	return c.callWithRetry(req, result, c.retryAttempts)
}

func (c RetryClient) callWithRetry(req gotransip.SoapRequest, result interface{}, attempts int) error {
	err := c.soapClient.Call(req, result)
	if err != nil {
		hasSkippableError := false
		for _, v := range []string{
			"SOAP Fault 100",
		} {
			if strings.Contains(err.Error(), v) {
				hasSkippableError = true
			}
		}

		if !hasSkippableError {
			return err
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(c.retryJitter)))
			sleep := c.retrySleep + jitter

			time.Sleep(sleep)
			return c.callWithRetry(req, result, attempts)
		}
		return err
	}
	return err
}
