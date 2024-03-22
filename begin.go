package pasdk

// BeginRequest begins the application process. Nullable fields are generally optional.
type BeginRequest struct {
	OrderID                  string  // A unique invoice ID or order ID.
	Amount                   int     // The invoice amount in pence.
	CustomerFirstName        string  // The customer's first name.
	CustomerLastName         string  // The customer's last name.
	CustomerAddress1         string  // The first line of the customer's address.
	CustomerAddress2         *string // The second line of the customer's address.
	CustomerAddress3         *string // The third line of the customer's address.
	CustomerTown             *string // The customer's town.
	CustomerCounty           *string // The customer's county.
	CustomerPostcode         string  // The customer's postcode.
	CustomerEmail            *string // The customer's email address. This is required if SendEmail is true.
	CustomerTelephone        *string // The customer's telephone number. This is required if SendSMS is true.
	SendEmail                *bool   // Whether to send the application link to the customer via email. Defaults to false.
	SendSMS                  *bool   // Whether to send the application link to the customer via SMS. Defaults to false.
	EnableMultiPlan          *bool   // If true, the customer will see a list of all available payment plans and will be able to select one themselves. Defaults to false.
	ReturnQRCode             *bool   // If true, a base64-encoded QR code will be returned, which the customer can scan with a mobile device to continue the application. Defaults to false.
	EnableAutoCapture        *bool   // Enables auto-capture (see https://api-docs.payment-assist.co.uk/auto-capture). Defaults to true.
	FailureURL               *string // A URL you want the customer to be redirected to when the application is denied.
	SuccessURL               *string // A URL you want the customer to be redirected to when the application is approved.
	WebhookURL               *string // A callback URL for receiving webhooks (see https://api-docs.payment-assist.co.uk/webhooks).
	PlanID                   *int    // The ID of the application's plan type. This is required if the account has access to multiple plan types and EnableMultiPlan is false.
	VehicleRegistrationPlate *string // The vehicle's registration plate, where relevant.
	Description              *string // A description of the services or goods being sold.
	Expiry                   *int    // The amount of time before the application expires, in seconds. This is 24 hours by default.
}

// BeginResponse contains the data returned by a successful call to the "begin" endpoint.
type BeginResponse struct {
	ApplicationToken string `json:"token"` // A token representing the application that was created. You should save this for later use.
	ContinuationURL  string `json:"url"`   // The URL you should direct the customer to so that they can complete the rest of the signup process.
}

// Fetch executes the request.
func (request BeginRequest) Fetch() (response *BeginResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	request = applyBeginDefaults(request)
	err = validateBeginRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"addr1=" + request.CustomerAddress1,
		"addr2=" + toString(request.CustomerAddress2),
		"addr3=" + toString(request.CustomerAddress3),
		"amount=" + toString(request.Amount),
		"auto_capture=" + toString(request.EnableAutoCapture),
		"county=" + toString(request.CustomerCounty),
		"description=" + toString(request.Description),
		"email=" + toString(request.CustomerEmail),
		"expiry=" + toString(request.Expiry),
		"f_name=" + request.CustomerFirstName,
		"failure_url=" + toString(request.FailureURL),
		"multi_plan=" + toString(request.EnableMultiPlan),
		"order_id=" + request.OrderID,
		"plan_id=" + toString(request.PlanID),
		"postcode=" + request.CustomerPostcode,
		"qr_code=" + toString(request.ReturnQRCode),
		"reg_no=" + toString(request.VehicleRegistrationPlate),
		"s_name=" + request.CustomerLastName,
		"send_email=" + toString(request.SendEmail),
		"send_sms=" + toString(request.SendSMS),
		"success_url=" + toString(request.SuccessURL),
		"telephone=" + toString(request.CustomerTelephone),
		"town=" + toString(request.CustomerTown),
		"webhook_url=" + toString(request.WebhookURL),
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, userCredentials.APISecret)

	requestParams = append(requestParams, "api_key="+userCredentials.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL()

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIPOSTRequest[BeginResponse](requestParams, requestURL+"begin")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func applyBeginDefaults(params BeginRequest) BeginRequest {
	falseValue := false
	trueValue := true

	if params.SendEmail == nil {
		params.SendEmail = &falseValue
	}

	if params.SendSMS == nil {
		params.SendSMS = &falseValue
	}

	if params.EnableMultiPlan == nil {
		params.EnableMultiPlan = &falseValue
	}

	if params.ReturnQRCode == nil {
		params.ReturnQRCode = &falseValue
	}

	if params.EnableAutoCapture == nil {
		params.EnableAutoCapture = &trueValue
	}

	return params
}

func validateBeginRequest(request BeginRequest) (err *PASDKError) {
	if len(request.OrderID) == 0 {
		return buildValidationFailedError("OrderID cannot be empty")
	}

	if request.Amount <= 0 {
		return buildValidationFailedError("field Amount must be greater than 0")
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

	if request.SendEmail != nil &&
		*request.SendEmail &&
		(request.CustomerEmail == nil || len(*request.CustomerEmail) == 0) {
		return buildValidationFailedError("CustomerEmail cannot be empty if SendEmail is true")
	}

	if request.SendSMS != nil &&
		*request.SendSMS &&
		(request.CustomerTelephone == nil || len(*request.CustomerTelephone) == 0) {
		return buildValidationFailedError("CustomerTelephone cannot be empty if SendSMS is true")
	}

	return nil
}
