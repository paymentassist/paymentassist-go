package pasdk

import (
	"testing"
)

func Test_Capture(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := CaptureRequest{
		APISecret:     getTestAPISecret(),
		APIKey:        getTestAPIKey(),
		ApplicationID: "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.ApplicationID != "aed3bd4e-c478-4d73-a6fa-3640a7155e4f" {
		t.Error()
	}
	if response.Status != "completed" {
		t.Error()
	}
	if *response.DepositCaptured != true {
		t.Error()
	}
	if response.DepositCaptureFailureReason != nil {
		t.Error()
	}
}

func Test_validateCaptureRequest(t *testing.T) {
	request := CaptureRequest{}

	if validateCaptureRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.APIKey = "test"

	if validateCaptureRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.APISecret = "test"

	if validateCaptureRequest(request).Error() != "ApplicationID cannot be empty" {
		t.Error()
	}

	request.ApplicationID = "test"

	if validateCaptureRequest(request) != nil {
		t.Error()
	}
}
