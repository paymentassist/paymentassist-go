package pasdk

// Account returns information about an account and its available plan types.
type AccountRequest struct {
	APIKey    string // Your API key.
	APISecret string // Your API secret.
}

// The data returned by a successful call to the "account" endpoint.
type AccountResponse struct {
	LegalName   string `json:"legal_name"`   // The legal name of the dealer.
	DisplayName string `json:"display_name"` // The display name of the dealer.
	Plans       []Plan `json:"plans"`        // A list of available plan types for this dealer.
}

// Execute the request.
func (request AccountRequest) Fetch() (response *AccountResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateAccountRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	signature := generateSignature([]string{}, request.APISecret)

	requestParams := []string{
		"api_key=" + request.APIKey,
		"signature=" + signature,
	}

	requestURL, err := getRequestURL(request.APISecret)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIGETRequest[AccountResponse](requestParams, requestURL+"account")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validateAccountRequest(request AccountRequest) (err *PASDKError) {
	if len(request.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty")
	}

	if len(request.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty")
	}

	return nil
}
