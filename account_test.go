package pasdk

import (
	"testing"
)

func Test_Account(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := AccountRequest{
		APISecret: getTestAPISecret(),
		APIKey:    getTestAPIKey(),
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.DisplayName != "Test Dealer" {
		t.Error()
	}
	if response.LegalName != "Test Dealer" {
		t.Error()
	}

	if response.Plans[0].ID != 6 {
		t.Error()
	}
	if response.Plans[0].Name != "3-Payment" {
		t.Error()
	}
	if response.Plans[0].Instalments != 3 {
		t.Error()
	}
	if response.Plans[0].Deposit != true {
		t.Error()
	}
	if response.Plans[0].APR != "0" {
		t.Error()
	}
	if response.Plans[0].Frequency != "monthly" {
		t.Error()
	}
	if response.Plans[0].MinAmount != nil {
		t.Error()
	}
	if *response.Plans[0].MaxAmount != 500000 {
		t.Error()
	}
	if response.Plans[0].CommissionRate != "8.50" {
		t.Error()
	}
	if response.Plans[0].CommissionFixedFee != nil {
		t.Error()
	}

	if response.Plans[1].ID != 1 {
		t.Error()
	}
	if response.Plans[1].Name != "4-Payment" {
		t.Error()
	}
	if response.Plans[1].Instalments != 4 {
		t.Error()
	}
	if response.Plans[1].Deposit != false {
		t.Error()
	}
	if response.Plans[1].APR != "5.5" {
		t.Error()
	}
	if response.Plans[1].Frequency != "monthly" {
		t.Error()
	}
	if *response.Plans[1].MinAmount != 10000 {
		t.Error()
	}
	if *response.Plans[1].MaxAmount != 300000 {
		t.Error(*response.Plans[1].MaxAmount)
	}
	if response.Plans[1].CommissionRate != "0" {
		t.Error()
	}
	if *response.Plans[1].CommissionFixedFee != 5000 {
		t.Error()
	}
}

func Test_validateAccountRequest(t *testing.T) {
	request := AccountRequest{}

	if validateAccountRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.APIKey = "test"

	if validateAccountRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.APISecret = "test"

	if validateAccountRequest(request) != nil {
		t.Error()
	}
}
