package pasdk

import (
	"strings"
	"testing"
)

func Test_Begin(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := BeginRequest{
		OrderID:           "111",
		Amount:            50000,
		CustomerFirstName: "Test",
		CustomerLastName:  "Testington",
		CustomerAddress1:  "Test House",
		CustomerPostcode:  "TEST TES",
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if len(response.ApplicationID) != 36 {
		t.Error()
	}

	if !strings.Contains(response.ContinuationURL, "https://") {
		t.Error(response.ContinuationURL)
	}
}

func Test_validateBeginRequest(t *testing.T) {
	request := BeginRequest{}

	if validateBeginRequest(request).Error() != "OrderID cannot be empty" {
		t.Error()
	}

	request.OrderID = "test"

	if validateBeginRequest(request).Error() != "field Amount must be greater than 0" {
		t.Error()
	}

	request.Amount = 50000

	if validateBeginRequest(request).Error() != "CustomerFirstName cannot be empty" {
		t.Error()
	}

	request.CustomerFirstName = "test"

	if validateBeginRequest(request).Error() != "CustomerLastName cannot be empty" {
		t.Error()
	}

	request.CustomerLastName = "test"

	if validateBeginRequest(request).Error() != "CustomerAddress1 cannot be empty" {
		t.Error()
	}

	request.CustomerAddress1 = "test"

	if validateBeginRequest(request).Error() != "CustomerPostcode cannot be empty" {
		t.Error()
	}

	request.CustomerPostcode = "test"

	if validateBeginRequest(request) != nil {
		t.Error()
	}

	trueValue := true
	test := "test"

	request.SendEmail = &trueValue

	if validateBeginRequest(request).Error() != "CustomerEmail cannot be empty if SendEmail is true" {
		t.Error()
	}

	request.CustomerEmail = &test

	if validateBeginRequest(request) != nil {
		t.Error()
	}

	request.SendSMS = &trueValue

	if validateBeginRequest(request).Error() != "CustomerTelephone cannot be empty if SendSMS is true" {
		t.Error()
	}

	request.CustomerTelephone = &test

	if validateBeginRequest(request) != nil {
		t.Error()
	}
}

func Test_applyBeginDefaults(t *testing.T) {
	request := BeginRequest{}

	request = applyBeginDefaults(request)

	if *request.SendEmail != false {
		t.Error()
	}
	if *request.SendSMS != false {
		t.Error()
	}
	if *request.EnableMultiPlan != false {
		t.Error()
	}
	if *request.ReturnQRCode != false {
		t.Error()
	}
	if *request.EnableAutoCapture != true {
		t.Error()
	}
}
