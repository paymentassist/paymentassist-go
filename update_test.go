package pasdk

import (
	"testing"
)

func Test_Update(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	amount := 100000
	orderID := "neworderid"
	expiresIn := 600

	request := UpdateRequest{
		APISecret:     getTestAPISecret(),
		APIKey:        getTestAPIKey(),
		ApplicationID: "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
		Amount:        &amount,
		OrderID:       &orderID,
		ExpiresIn:     &expiresIn,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.ApplicationID != "aed3bd4e-c478-4d73-a6fa-3640a7155e4f" {
		t.Error()
	}
	if *response.Amount != 100000 {
		t.Error()
	}
	if *response.OrderID != "neworderid" {
		t.Error()
	}
	if *response.ExpiresIn != 600 {
		t.Error()
	}
}

func Test_validateUpdateRequest(t *testing.T) {
	request := StatusRequest{}

	if validateStatusRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.APIKey = "test"

	if validateStatusRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.APISecret = "test"

	if validateStatusRequest(request).Error() != "ApplicationID cannot be empty" {
		t.Error()
	}

	request.ApplicationID = "test"

	if validateStatusRequest(request) != nil {
		t.Error()
	}
}
