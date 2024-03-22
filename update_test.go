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
		ApplicationToken: "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
		Amount:           &amount,
		OrderID:          &orderID,
		ExpiresIn:        &expiresIn,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.ApplicationToken != "aed3bd4e-c478-4d73-a6fa-3640a7155e4f" {
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

	if validateStatusRequest(request).Error() != "ApplicationToken cannot be empty" {
		t.Error()
	}

	request.ApplicationToken = "test"

	if validateStatusRequest(request) != nil {
		t.Error()
	}
}
