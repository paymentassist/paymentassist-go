package pasdk

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	testsAreRunning = true
	os.Exit(m.Run())
}

func Test_decodeResponseJSON_DoesntLosePrecision_WhenDeserialisingPlan(t *testing.T) {
	response, err := decodeResponseJSON[AccountResponse]([]byte(`
		{
			"status": "ok",
			"msg": null,
			"data": {
				"legal_name": "Test Dealer",
				"display_name": "Test Dealer",
				"plans": [
					{
						"plan_id": 6,
						"name": "3-Payment",
						"instalments": 3,
						"deposit": true,
						"apr": 57654574587358234524123402034972368234875238510615130945139456.89174635871346578136571365136571346587141398567348571347,
						"frequency": "monthly",
						"min_amount": null,
						"max_amount": 500000,
						"commission_rate": "8.50",
						"commission_fixed_fee": null
					}
				]
			}
		}`))

	if err != nil {
		t.Error()
	}

	if response.Plans[0].APR != "57654574587358234524123402034972368234875238510615130945139456.89174635871346578136571365136571346587141398567348571347" {
		t.Error()
	}
}

func Test_getRequestURL(t *testing.T) {
	if shouldRunIntegrationTests() {
		return
	}

	defer func() {
		os.Unsetenv("GO_PASDK_INTEGRATION_TESTS")
		testsAreRunning = true
	}()

	testsAreRunning = false

	url, err := getRequestURL(PAAuth{
		APIURL: "https://testurl"},
	)

	if url != "https://testurl/" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}

	url, err = getRequestURL(PAAuth{
		APIURL: "https://testurl/",
	})

	if url != "https://testurl/" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}

	url, err = getRequestURL(PAAuth{
		APIURL: "www.testurl",
	})

	if url != "" {
		t.Error()
	}
	if err.Error() != "the API URL must contain the string \"https:\"" {
		t.Error(err.Error())
	}
}

func Test_decodeResponseJSON_HandlesAPIRefusalsProperly(t *testing.T) {
	response, err := decodeResponseJSON[AccountResponse]([]byte(
		`{ "status": "error", "msg": null, "data": null }`))

	if response != nil {
		t.Error()
	}
	if err == nil {
		t.Error()
		return
	}
	if !err.IsRequestRefusedError {
		t.Error()
	}
	if err.GetErrorType() != "RequestRefusedError" {
		t.Error()
	}
	if err.Error() != `the API refused your request: { "status": "error", "msg": null, "data": null }` {
		t.Error(err.Error())
	}
}

func Test_catchGenericPanic(t *testing.T) {
	response := "test"

	responsePointer := &response
	var err *PASDKError

	func() {
		defer catchGenericPanic(&responsePointer, &err)

		panic("panic!")
	}()

	if responsePointer != nil {
		t.Error()
	}
	if !strings.Contains(err.Error(), "unexpected panic") {
		t.Error()
	}
	if !err.IsUnexpectedError {
		t.Error()
	}
	if err.GetErrorType() != "UnexpectedError" {
		t.Error()
	}
}

func Test_checkStatusCode(t *testing.T) {
	err := checkStatusCode(50, "fail")

	if !err.IsUnexpectedError {
		t.Error()
	}
	if !strings.Contains(err.Error(), "API request failed returning status code 50") {
		t.Error(err.Error())
	}

	err = checkStatusCode(150, "fail")

	if !err.IsUnexpectedError {
		t.Error()
	}
	if !strings.Contains(err.Error(), "API request failed returning status code 150") {
		t.Error(err.Error())
	}

	err = checkStatusCode(250, "success")

	if err != nil {
		t.Error()
	}

	err = checkStatusCode(350, "fail")

	if !err.IsUnexpectedError {
		t.Error()
	}
	if !strings.Contains(err.Error(), "API request failed returning status code 350") {
		t.Error(err.Error())
	}

	err = checkStatusCode(450, "fail")

	if !err.IsRequestRefusedError {
		t.Error()
	}
	if !strings.Contains(err.Error(), "API refused your request returning status code 450") {
		t.Error(err.Error())
	}

	err = checkStatusCode(550, "fail")

	if !err.IsUnexpectedError {
		t.Error()
	}
	if !strings.Contains(err.Error(), "API request failed returning status code 550") {
		t.Error(err.Error())
	}
}

func Test_decodeResponseJSON(t *testing.T) {
	// Test it handles empty bodies properly.
	response, err := decodeResponseJSON[BeginResponse]([]byte{})

	if response != nil {
		t.Error()
	}
	if err.Error() != "the response from the API was malformed: the response body was empty" {
		t.Error(err)
	}

	response, err = decodeResponseJSON[BeginResponse](nil)

	if response != nil {
		t.Error()
	}
	if err.Error() != "the response from the API was malformed: the response body was empty" {
		t.Error(err)
	}

	// Test it handles invalid JSON properly.
	response, err = decodeResponseJSON[BeginResponse]([]byte("{Test"))

	if response != nil {
		t.Error()
	}
	if !strings.Contains(err.Error(), "failed to parse API response") {
		t.Error(err)
	}
}

func Test_toString(t *testing.T) {
	stringValue := "test"
	intValue := 5
	floatValue := 5.555
	boolValue := true

	if toString(nil) != "" {
		t.Error()
	}
	if toString(stringValue) != "test" {
		t.Error()
	}
	if toString(&stringValue) != "test" {
		t.Error()
	}
	if toString(intValue) != "5" {
		t.Error()
	}
	if toString(&intValue) != "5" {
		t.Error()
	}
	if toString(floatValue) != "5.555" {
		t.Error()
	}
	if toString(&floatValue) != "5.555" {
		t.Error()
	}
	if toString(boolValue) != "true" {
		t.Error()
	}
	if toString(&boolValue) != "true" {
		t.Error()
	}
}

func Test_generateSignature(t *testing.T) {
	params := []string{
		"test=test",
		"test2=test2",
	}

	hash := "7eba7f616af343d16ff09e242362345e6cfb09d24b78a73c81d267f049fc47c2"
	output := generateSignature(params, "secret")

	if output != hash {
		t.Error(output)
	}
}

func Test_generateSignature_2(t *testing.T) {
	params := []string{
		"test1=test",
		"test2=test2",
	}

	hash := "8226de39365226038be9598213e480d22f4dfe7147f50d977087a8d4eb124f52"
	output := generateSignature(params, "demo_2ec4449ac4a7a86e2f79c4794e8")

	if output != hash {
		t.Error(output)
	}
}

func Test_capitaliseParamKeys(t *testing.T) {
	keys := []string{
		"test1=test test",
		"test2=test",
	}

	keys = capitaliseParamKeys(keys)

	if keys[0] != "TEST1=test test" {
		t.Error()
	}

	if keys[1] != "TEST2=test" {
		t.Error()
	}
}
