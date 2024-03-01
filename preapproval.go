package pasdk

// Pre-approval allows you to check the eligibity of a customer in advance.
// Success simply means that the customer has passed our internal checks. They
// will still need to have funds available to cover any deposit payment for
// the application to be successful.
type PreapprovalRequest struct {
	APIKey            string // Your API key.
	APISecret         string // Your API secret.
	CustomerFirstName string // The customer's first name.
	CustomerLastName  string // The customer's last name.
	CustomerPostcode  string // The customer's postode.
	CustomerAddress1  string // The first line of the customer's address.
}

// The data returned by a successful call to the "preapproval" endpoint.
type PreapprovalResponse struct {
	Approved bool `json:"approved"` // Whether or not this customer passed the pre-approval checks.
}

// Execute the request.
func (request PreapprovalRequest) Fetch() (response *PreapprovalResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validatePreapprovalRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"addr1=" + request.CustomerAddress1,
		"f_name=" + request.CustomerFirstName,
		"postcode=" + request.CustomerPostcode,
		"s_name=" + request.CustomerLastName,
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, request.APISecret)

	requestParams = append(requestParams, "api_key="+request.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL(request.APISecret)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIPOSTRequest[PreapprovalResponse](requestParams, requestURL+"preapproval")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validatePreapprovalRequest(request PreapprovalRequest) (err *PASDKError) {
	if len(request.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty")
	}

	if len(request.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty")
	}

	if len(request.CustomerFirstName) == 0 {
		return buildValidationFailedError("CustomerFirstName cannot be empty")
	}

	if len(request.CustomerLastName) == 0 {
		return buildValidationFailedError("CustomerLastName cannot be empty")
	}

	if len(request.CustomerAddress1) == 0 {
		return buildValidationFailedError("CustomerAddress1 cannot be empty")
	}

	if len(request.CustomerPostcode) == 0 {
		return buildValidationFailedError("CustomerPostcode cannot be empty")
	}

	return nil
}