package pasdk

// PlanRequest accepts a transaction amount and an optional plan ID and term length,
// returning a full payment schedule including amounts and dates.
type PlanRequest struct {
	AuthInfo PAAuth // Authentication information including your API credentials.
	Amount   int    // The invoice amount in pence.
	PlanID   *int   // The plan ID. If empty, the account's default plan is used.
}

// PlanResponse contains the data returned by a successful call to the "plan" endpoint.
type PlanResponse struct {
	PlanName        string      `json:"plan"`      // The name of this plan.
	Amount          int         `json:"amount"`    // The amount you requested, in pence.
	Interest        int         `json:"interest"`  // The amount of interest payable, in pence.
	TotalRepayable  int         `json:"repayable"` // The total amount that would be repayable under this plan, in pence.
	PaymentSchedule []Repayment `json:"schedule"`  // A breakdown of what the repayments would look like under this plan.
}

// Fetch executes the request.
func (request PlanRequest) Fetch() (response *PlanResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validatePlanRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"amount=" + toString(request.Amount),
		"plan_id=" + toString(request.PlanID),
	}

	requestParams = removeEmptyParams(requestParams)

	signature := generateSignature(requestParams, request.AuthInfo.APISecret)

	requestParams = append(requestParams, "api_key="+request.AuthInfo.APIKey)
	requestParams = append(requestParams, "signature="+signature)

	requestURL, err := getRequestURL(request.AuthInfo)

	if err != nil {
		return nil, err.Wrap("failed determining request URL: ")
	}

	response, err = makeAPIPOSTRequest[PlanResponse](requestParams, requestURL+"plan")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validatePlanRequest(request PlanRequest) (err *PASDKError) {
	if len(request.AuthInfo.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty")
	}

	if len(request.AuthInfo.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty")
	}

	if request.Amount == 0 {
		return buildValidationFailedError("field Amount must be greater than 0")
	}

	return nil
}
