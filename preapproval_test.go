package pasdk

import (
	"testing"
)

func Test_Preapproval(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := PreapprovalRequest{
		AuthInfo: PAAuth{
			APISecret: getTestAPISecret(),
			APIKey:    getTestAPIKey(),
		},
		CustomerFirstName: "Test",
		CustomerLastName:  "Testington",
		CustomerAddress1:  "Test House",
		CustomerPostcode:  "TEST TES",
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if !response.Approved {
		t.Error()
	}
}

func Test_validatePreapprovalRequest(t *testing.T) {
	request := PreapprovalRequest{}

	if validatePreapprovalRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.AuthInfo.APIKey = "test"

	if validatePreapprovalRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.AuthInfo.APISecret = "test"

	if validatePreapprovalRequest(request).Error() != "CustomerFirstName cannot be empty" {
		t.Error()
	}

	request.CustomerFirstName = "test"

	if validatePreapprovalRequest(request).Error() != "CustomerLastName cannot be empty" {
		t.Error()
	}

	request.CustomerLastName = "test"

	if validatePreapprovalRequest(request).Error() != "CustomerAddress1 cannot be empty" {
		t.Error()
	}

	request.CustomerAddress1 = "test"

	if validatePreapprovalRequest(request).Error() != "CustomerPostcode cannot be empty" {
		t.Error()
	}

	request.CustomerPostcode = "test"

	if validatePreapprovalRequest(request) != nil {
		t.Error()
	}
}
