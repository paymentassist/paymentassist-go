package pasdk

// CaptureRequest allows you to finalise an application that's currently in a "pending_capture" state.
type CaptureRequest struct {
	ApplicationToken string // The token you received when calling the "begin" endpoint.
}

// CaptureResponse contains the data returned by a call to the "capture" endpoint. Unlike some other
// endpoints, "capture" can return a response even when unsuccessful.
type CaptureResponse struct {
	ApplicationToken            string  `json:"token"`            // The token representing this application.
	Status                      string  `json:"status"`           // The status of this application after the application was captured.
	DepositCaptured             *bool   `json:"deposit_captured"` // Indicates whether the deposit was successfully captured. This is always nil if the application does not include a deposit.
	DepositCaptureFailureReason *string `json:"deposit_reason"`   // If DepositCaptured is false, this contains the reason for capture failure. This is nil in all other situations.
}

// Fetch executes the request.
func (request CaptureRequest) Fetch() (response *CaptureResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateCaptureRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	// Alphabetically sorted.
	requestParams := []string{
		"token=" + toString(request.ApplicationToken),
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, userCredentials.APISecret)

	requestParams = append(requestParams, "api_key="+userCredentials.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL()

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
	if len(request.ApplicationToken) == 0 {
		return buildValidationFailedError("ApplicationToken cannot be empty")
	}

	return nil
}
