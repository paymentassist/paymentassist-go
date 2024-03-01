package pasdk

import (
	"testing"
)

func Test_Invoice(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	request := InvoiceRequest{
		APISecret:     getTestAPISecret(),
		APIKey:        getTestAPIKey(),
		ApplicationID: "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
		FileType:      "txt",
		FileData:      []byte("Test invoice for £100"),
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.ApplicationID != "aed3bd4e-c478-4d73-a6fa-3640a7155e4f" {
		t.Error()
	}
	if response.UploadStatus != "success" {
		t.Error()
	}
}

func Test_validateInvoiceRequest(t *testing.T) {
	request := InvoiceRequest{}

	if validateInvoiceRequest(request).Error() != "APIKey cannot be empty" {
		t.Error()
	}

	request.APIKey = "test"

	if validateInvoiceRequest(request).Error() != "APISecret cannot be empty" {
		t.Error()
	}

	request.APISecret = "test"

	if validateInvoiceRequest(request).Error() != "ApplicationID cannot be empty" {
		t.Error()
	}

	request.ApplicationID = "test"

	if validateInvoiceRequest(request).Error() != "FileType cannot be empty" {
		t.Error(validateInvoiceRequest(request).Error())
	}

	request.FileType = "test"

	if validateInvoiceRequest(request).Error() != "FileData cannot be empty" {
		t.Error()
	}

	request.FileData = []byte("test")

	if validateInvoiceRequest(request) != nil {
		t.Error()
	}
}
