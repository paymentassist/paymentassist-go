package pasdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var requestClient *http.Client

func getAPIRequestClient() *http.Client {
	if requestClient == nil {
		requestClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return requestClient
}

func getRequestURL() (string, *PASDKError) {
	if testsAreRunning && !shouldRunIntegrationTests() {
		return "", nil
	}

	if !strings.Contains(userCredentials.APIURL, "https:") {
		return "", buildValidationFailedError("the API URL must contain the string \"https:\"")
	}

	if userCredentials.APIURL[len(userCredentials.APIURL)-1:] != "/" {
		return userCredentials.APIURL + "/", nil
	}

	return userCredentials.APIURL, nil
}

// Returns an error if there is an issue with the credentials.
func checkCredentialsExist() *PASDKError {
	if len(userCredentials.APIKey) == 0 {
		return buildValidationFailedError("APIKey cannot be empty - call pasdk.Initialise to pass in your credentials")
	}
	if len(userCredentials.APISecret) == 0 {
		return buildValidationFailedError("APISecret cannot be empty - call pasdk.Initialise to pass in your credentials")
	}

	if !testsAreRunning && len(userCredentials.APIURL) == 0 {
		return buildValidationFailedError("APIURL cannot be empty - call pasdk.Initialise to pass in the URL you want to send a request to")
	}

	return nil
}

func makeAPIPOSTRequest[T interface{}](formData []string, endpoint string) (*T, *PASDKError) {
	paErr := checkCredentialsExist()

	if paErr != nil {
		return nil, paErr
	}

	formValues := url.Values{}

	for _, data := range formData {
		parts := strings.Split(data, "=")
		formValues.Set(parts[0], parts[1])
	}

	if testsAreRunning && !shouldRunIntegrationTests() {
		response, err := getMockAPIResponse[T](endpoint)
		return response, err
	}

	request, err := http.NewRequest("POST", endpoint, strings.NewReader(formValues.Encode()))

	if err != nil {
		return nil, buildUnexpectedError("creating API request failed: " + err.Error())
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("X-Origin", "payment-assist-go-sdk")

	response, err := getAPIRequestClient().Do(request)

	if err != nil {
		return nil, buildUnexpectedError("API request failed: " + err.Error())
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, buildUnexpectedError("reading API response failed: " + err.Error())
	}

	paErr = checkStatusCode(response.StatusCode, string(body))

	if paErr != nil {
		return nil, paErr
	}

	output, paErr := decodeResponseJSON[T](body)

	if paErr != nil {
		return nil, paErr
	}

	return output, nil
}

// Returns an error if the status code indicated failure.
func checkStatusCode(statusCode int, requestBody string) *PASDKError {
	if (statusCode >= 0 && statusCode < 200) ||
		(statusCode >= 300 && statusCode < 400) ||
		(statusCode >= 500 && statusCode < 600) {
		return buildUnexpectedError("API request failed returning status code " + toString(statusCode) +
			": " + requestBody)
	}

	if statusCode >= 400 && statusCode < 500 {
		return buildRequestRefusedError("API refused your request returning status code " + toString(statusCode) +
			": " + requestBody)
	}

	return nil
}

func makeAPIGETRequest[T interface{}](formData []string, endpoint string) (*T, *PASDKError) {
	paErr := checkCredentialsExist()

	if paErr != nil {
		return nil, paErr
	}

	endpoint += "?"

	for _, data := range formData {
		parts := strings.Split(data, "=")

		endpoint += parts[0] + "=" + url.QueryEscape(parts[1]) + "&"
	}

	endpoint = endpoint[:len(endpoint)-1]

	if testsAreRunning && !shouldRunIntegrationTests() {
		response, err := getMockAPIResponse[T](endpoint)
		return response, err
	}

	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, buildUnexpectedError("creating API request failed: " + err.Error())
	}

	request.Header.Add("X-Origin", "payment-assist-go-sdk")

	response, err := getAPIRequestClient().Do(request)

	if err != nil {
		return nil, buildUnexpectedError("API request failed: " + err.Error())
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, buildUnexpectedError("reading API response failed: " + err.Error())
	}

	paErr = checkStatusCode(response.StatusCode, string(body))

	if paErr != nil {
		return nil, paErr
	}

	output, paErr := decodeResponseJSON[T](body)

	if paErr != nil {
		return nil, paErr
	}

	return output, nil
}

func catchGenericPanic[T interface{}](response **T, err **PASDKError) {
	data := recover()

	if data != nil {
		*response = nil

		*err = &PASDKError{
			errorMessage: "there was an unexpected panic; this may indicate a bug, " +
				"please contact support: " + fmt.Sprintf("%v", data),
			IsUnexpectedError: true,
		}
	}
}

func buildValidationFailedError(message string) *PASDKError {
	return &PASDKError{
		IsValidationFailedError: true,
		errorMessage:            message,
	}
}

func buildRequestRefusedError(message string) *PASDKError {
	return &PASDKError{
		IsRequestRefusedError: true,
		errorMessage:          message,
	}
}

func buildUnexpectedError(message string) *PASDKError {
	return &PASDKError{
		errorMessage:      message,
		IsUnexpectedError: true,
	}
}

// Returns an error if the request failed, or if something else went wrong.
func decodeResponseJSON[T interface{}](jsonData []byte) (*T, *PASDKError) {
	if len(jsonData) == 0 {
		return nil, buildUnexpectedError("the response from the API was malformed: the response body was empty")
	}

	// The data field can change depending on whether the request was successful or not,
	// so we need to first check for success before unmarshaling it into T.
	var statusResponseWrapper struct {
		Status string `json:"status"`
	}

	err := json.Unmarshal(jsonData, &statusResponseWrapper)

	if err != nil {
		return nil, buildUnexpectedError("failed to parse API response: " + err.Error())
	}

	if statusResponseWrapper.Status == "error" {
		return nil, buildRequestRefusedError("the API refused your request: " + string(jsonData))
	}

	if statusResponseWrapper.Status != "ok" {
		return nil, buildUnexpectedError("the API returned an unexpected response: " + string(jsonData))
	}

	// Now we can be sure the response was successful.
	var responseWrapper struct {
		Status  string  `json:"status"`
		Message *string `json:"msg"`
		Data    T       `json:"data"`
	}

	err = json.Unmarshal(jsonData, &responseWrapper)

	if err != nil {
		return nil, buildUnexpectedError("parsing JSON failed: " + err.Error())
	}

	return &responseWrapper.Data, nil
}

// The keys of requestParams should already be in alphabetical order.
func generateSignature(requestParams []string, secret string) string {
	requestParams = capitaliseParamKeys(requestParams)
	requestString := strings.Join(requestParams, "&")

	if len(requestString) > 0 {
		requestString += "&"
	}

	hasher := hmac.New(sha256.New, []byte(secret))
	hasher.Write([]byte(requestString))
	return hex.EncodeToString(hasher.Sum(nil))
}

func capitaliseParamKeys(params []string) []string {
	output := make([]string, 0, len(params))

	for _, param := range params {
		parts := strings.Split(param, "=")

		output = append(output, strings.ToUpper(parts[0])+"="+parts[1])
	}

	return output
}

func removeEmptyParams(params []string) []string {
	output := make([]string, 0, len(params))

	for _, param := range params {
		value := strings.Split(param, "=")[1]

		if len(value) > 0 {
			output = append(output, param)
		}
	}

	return output
}

func toString(input interface{}) string {
	if input == nil {
		return ""
	}

	switch input := input.(type) {
	case *string:
		if input == nil {
			return ""
		}

		return *input
	case string:
		return input
	case int:
		return strconv.Itoa(input)
	case *int:
		if input == nil {
			return ""
		}

		return strconv.Itoa(*input)
	case bool:
		if input {
			return "true"
		}

		return "false"
	case *bool:
		if input == nil {
			return ""
		}

		if *input {
			return "true"
		}

		return "false"
	case float64:
		return strconv.FormatFloat(input, 'f', -1, 64)
	case *float64:
		if input == nil {
			return ""
		}

		return strconv.FormatFloat(*input, 'f', -1, 64)
	default:
		// Be stricter when running tests.
		if testsAreRunning {
			panic("unrecognised type")
		}

		return fmt.Sprintf("%v", input)
	}
}
