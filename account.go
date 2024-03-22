package pasdk

// AccountRequest returns information about an account and its available plan types.
type AccountRequest struct{}

// AccountResponse contains the data returned by a successful call to the "account" endpoint.
type AccountResponse struct {
	LegalName   string `json:"legal_name"`   // The legal name of the merchant.
	DisplayName string `json:"display_name"` // The display name of the merchant.
	Plans       []Plan `json:"plans"`        // A list of available plan types for this merchant.
}

// Fetch executes the request.
func (request AccountRequest) Fetch() (response *AccountResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	signature := generateSignature([]string{}, userCredentials.APISecret)

	requestParams := []string{
		"api_key=" + userCredentials.APIKey,
		"signature=" + signature,
	}

	requestURL, err := getRequestURL()

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIGETRequest[AccountResponse](requestParams, requestURL+"account")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}
