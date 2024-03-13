package pasdk

import (
	"testing"
	"time"
)

func Test_Status(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := StatusRequest{
		ApplicationID: "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.ApplicationID != "aed3bd4e-c478-4d73-a6fa-3640a7155e4f" {
		t.Error()
	}
	if response.Status != "pending" {
		t.Error()
	}
	if response.Amount != 50000 {
		t.Error()
	}

	expires, _ := time.Parse(time.RFC3339, "2022-05-24T19:28:06+01:00")

	if !response.ExpiresAt.Equal(expires) {
		t.Error()
	}

	if !response.HasInvoice {
		t.Error()
	}
	if !response.RequriesInvoice {
		t.Error()
	}
	if response.PaymentAssistReference != "testreference" {
		t.Error()
	}
}

func Test_validateStatusRequest(t *testing.T) {
	request := StatusRequest{}

	if validateStatusRequest(request).Error() != "ApplicationID cannot be empty" {
		t.Error()
	}

	request.ApplicationID = "test"

	if validateStatusRequest(request) != nil {
		t.Error()
	}
}
