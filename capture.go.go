package pasdk

// Capture allows you to capture an application that's currently in a "pending_capture" state.
type CaptureRequest struct {
	APIKey        string // Your API key.
	APISecret     string // Your API secret.
	ApplicationID string // The application ID (token) you received when calling the "begin" endpoint.
}

// The data returned by a call to the "capture" endpoint. Unlike some other
// endpoints, "capture" can return a response even when unsuccessful.
type CaptureResponse struct {
	ApplicationID               string  `json:"token"`            // The ID (token) of this application.
	Status                      string  `json:"status"`           // The status of this application after the application was captured.
	DepositCaptured             *bool   `json:"deposit_captured"` // Indicates whether the deposit was successfully captured. This is always nil if the application does not include a deposit.
	DepositCaptureFailureReason *string `json:"deposit_reason"`   // If DepositCaptured is false, this contains the reason for capture failure. This is nil in all other situations.
}

// Execute the request.
func (request CaptureRequest) Fetch() (response *CaptureResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateCaptureRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"token=" + toString(request.ApplicationID),
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, request.APISecret)

	requestParams = append(requestParams, "api_key="+request.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL(request.APISecret)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIPOSTRequest[CaptureResponse](requestParams, requestURL+"capture")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validateCaptureRequest(request CaptureRequest) (err *PASDKError) {
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