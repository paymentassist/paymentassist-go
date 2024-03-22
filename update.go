package pasdk

import (
	"encoding/json"
	"errors"
	"strconv"
)

// UpdateRequest allows you to update an existing application.
type UpdateRequest struct {
	ApplicationToken string  // The token you received when calling the "begin" endpoint.
	OrderID          *string // Your new order ID. You can only change this if the application's status is "completed".
	ExpiresIn        *int    // The new expiry time for this appication in seconds from now. Setting this to 0 will instantly expire the application. You can only change this if the application's status is "pending", "in_progress" or "pending_capture".
	Amount           *int    // The new amount for this application in pence. You can only change this if the application's status is "pending", "in_progress" or "pending_capture". The new amount must be less than the current amount.
}

// UpdateResponse contains the data returned by a successful call to the "update" endpoint.
type UpdateResponse struct {
	ApplicationToken string  `json:"token"`    // The token representing this application.
	OrderID          *string `json:"order_id"` // The new order ID you requested, if any.
	ExpiresIn        *int    `json:"expiry"`   // The new expiry time you requested in seconds, if any.
	Amount           *int    `json:"amount"`   // The new amount you requested in pence, if any.
}

func (response *UpdateResponse) UnmarshalJSON(data []byte) error {
	type Alias UpdateResponse

	tmp := struct {
		ExpiresIn interface{} `json:"expiry"`
		Amount    interface{} `json:"amount"`
		*Alias
	}{
		Alias: (*Alias)(response), // Cast plan to Alias type, to unmarshal other fields normally.
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return errors.New("couldn't unmarshal UpdateResponse: " + err.Error())
	}

	// API returns it as a string but it's more helpful to have it as a number.
	if tmp.ExpiresIn != nil {
		expiresIn, err := strconv.Atoi(tmp.ExpiresIn.(string))

		if err != nil {
			return errors.New("failed to convert string to integer: " + err.Error())
		}

		response.ExpiresIn = &expiresIn
	}

	if tmp.Amount != nil {
		amount, err := strconv.Atoi(tmp.Amount.(string))

		if err != nil {
			return errors.New("failed to convert string to integer: " + err.Error())
		}

		response.Amount = &amount
	}

	return nil
}

// Fetch executes the request.
func (request UpdateRequest) Fetch() (response *UpdateResponse, err *PASDKError) {
	defer catchGenericPanic(&response, &err)

	err = validateUpdateRequest(request)

	if err != nil {
		return nil, err.Wrap("request is invalid: ")
	}

	requestParams := []string{
		"amount=" + toString(request.Amount),
		"expiry=" + toString(request.ExpiresIn),
		"order_id=" + toString(request.OrderID),
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

	response, err = makeAPIPOSTRequest[UpdateResponse](requestParams, requestURL+"update")

	if err != nil {
		return nil, err.Wrap("API request failed: ")
	}

	return response, nil
}

func validateUpdateRequest(request UpdateRequest) (err *PASDKError) {
	if len(request.ApplicationToken) == 0 {
		return buildValidationFailedError("ApplicationToken cannot be empty")
	}

	return nil
}
