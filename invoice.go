package pasdk

import (
	"encoding/base64"
)

// InvoiceRequest allows you to upload an invoice for a completed application.
type InvoiceRequest struct {
	AuthInfo      PAAuth // Authentication information including your API credentials.
	ApplicationID string // The application ID (token) you received when calling the "begin" endpoint.
	FileType      string // The file type. Some supported options are "pdf", "html", "txt", "doc" and "xls".
	FileData      []byte // The file as a slice of bytes.
}

// InvoiceResponse contains the data returned by a call to the "invoice" endpoint. Unlike some
// endpoints, "invoice" can return a response even if the upload was unsuccessful.
type InvoiceResponse struct {
	ApplicationID string `json:"token"`         // The ID (token) of this application.
	UploadStatus  string `json:"upload_status"` // The status of the upload ("success" or "failed").
}

// Fetch executes the request.
func (request InvoiceRequest) Fetch() (response *InvoiceResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateInvoiceRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"filedata=" + base64.StdEncoding.EncodeToString(request.FileData),
		"filetype=" + request.FileType,
		"token=" + request.ApplicationID,
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, request.AuthInfo.APISecret)

	requestParams = append(requestParams, "api_key="+request.AuthInfo.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL(request.AuthInfo)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIPOSTRequest[InvoiceResponse](requestParams, requestURL+"invoice")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validateInvoiceRequest(request InvoiceRequest) (err *PASDKError) {
	if len(request.AuthInfo.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty")
	}

	if len(request.AuthInfo.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty")
	}

	if len(request.ApplicationID) == 0 {
		return buildValidationFailedError("ApplicationID cannot be empty")
	}

	if len(request.FileType) == 0 {
		return buildValidationFailedError("FileType cannot be empty")
	}

	if len(request.FileData) == 0 {
		return buildValidationFailedError("FileData cannot be empty")
	}

	return nil
}
