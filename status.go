package pasdk

import (
	"time"
)

// StatusRequest allows you to check the status of an existing application.
type StatusRequest struct {
	ApplicationToken string // The token you received when calling the "begin" endpoint.
}

// StatusResponse contains the data returned by a successful call to the "status" endpoint.
type StatusResponse struct {
	ApplicationToken       string    `json:"token"`            // The token representing this application.
	Status                 string    `json:"status"`           // The status of this application.
	Amount                 int       `json:"amount"`           // The amount being applied for, in pence.
	ExpiresAt              time.Time `json:"expires_at"`       // The time this application expires.
	PaymentAssistReference string    `json:"pa_ref"`           // Payment Assist's reference for this application. This may be empty as a reference is not generated until the finance facility or payment is successfully created (once an application moves to a "completed" status).
	RequriesInvoice        bool      `json:"requires_invoice"` // Whether an invoice needs to be uploaded for this application before funds will be released to the merchant.
	HasInvoice             bool      `json:"has_invoice"`      // Whether an invoice has been uploaded for this application.
	LastAccessedAt         time.Time `json:"last_accessed_at"` // The last time the customer accessed the application.
}

// Fetch executes the request.
func (request StatusRequest) Fetch() (response *StatusResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateStatusRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	// Alphabetically sorted.
	requestParams := []string{
		"token=" + request.ApplicationToken,
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, userCredentials.APISecret)

	requestParams = append(requestParams, "api_key="+userCredentials.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL()

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
	if len(request.ApplicationToken) == 0 {
		return buildValidationFailedError("ApplicationToken cannot be empty")
	}

	return nil
}
