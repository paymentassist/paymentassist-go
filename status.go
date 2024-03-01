package pasdk

import (
	"time"
)

// Status allows you to check the status of an existing application.
type StatusRequest struct {
	APIKey        string // Your API key.
	APISecret     string // Your API secret.
	ApplicationID string // The application ID (token) you received when calling the "begin" endpoint.
}

// The data returned by a successful call to the "status" endpoint.
type StatusResponse struct {
	ApplicationID          string    `json:"token"`            // The ID (token) of this application.
	Status                 string    `json:"status"`           // The status of this application.
	Amount                 int       `json:"amount"`           // The amount being applied for, in pence.
	ExpiresAt              time.Time `json:"expires_at"`       // The time this application expires.
	PaymentAssistReference string    `json:"pa_ref"`           // Payment Assist's internal reference code for this application. This might be empty as an internal reference is not generated as soon as the application is started.
	RequriesInvoice        bool      `json:"requires_invoice"` // Whether an invoice needs to be uploaded for this application before payment can be made to the dealer.
	HasInvoice             bool      `json:"has_invoice"`      // Whether an invoice has been uploaded for this application.
}

// Execute the request.
func (request StatusRequest) Fetch() (response *StatusResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateStatusRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"token=" + request.ApplicationID,
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, request.APISecret)

	requestParams = append(requestParams, "api_key="+request.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL(request.APISecret)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIGETRequest[StatusResponse](requestParams, requestURL+"status")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validateStatusRequest(request StatusRequest) (err *PASDKError) {
	if len(request.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty")
	}

	if len(request.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty")
	}

	if len(request.ApplicationID) == 0 {
		return buildValidationFailedError("ApplicationID cannot be empty")
	}

	return nil
}