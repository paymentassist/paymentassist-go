package pasdk

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
)

var testsAreRunning = false

func shouldRunIntegrationTests() bool {
	_, exists := os.LookupEnv("GO_PASDK_INTEGRATION_TESTS")
	return exists
}

func getTestAPISecret() string {
	if shouldRunIntegrationTests() {
		return os.Getenv("GO_PASDK_TEST_SECRET")
	}

	return "testsecret"
}

func getTestAPIKey() string {
	if shouldRunIntegrationTests() {
		return os.Getenv("GO_PASDK_TEST_API_KEY")
	}

	return "testapikey"
}

// Returns a random 10 character string.
func getRandomID() string {
	randomBytes := make([]byte, 5)

	rand.Read(randomBytes)

	return hex.EncodeToString(randomBytes)[:10]
}

func getMockAPIResponse[T interface{}](endpoint string) (*T, *PASDKError) {
	// If this is a GET request then the endpoint will have parameters on it. Take them
	// off so we can match on the actual endpoint.
	if strings.Contains(endpoint, "?") {
		endpoint = strings.Split(endpoint, "?")[0]
	}

	switch endpoint {
	case "begin":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"token": "0138ef43-f703-41cb-8f08-f36f41b47560",
					"url": "https://example.com"
				}
			}`))
	case "preapproval":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"approved": true
				}
			}`))
	case "update":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"token": "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
					"order_id": "neworderid",
					"expiry": "600",
					"amount": "100000"
				}
			}`))
	case "plan":
		return decodeResponseJSON[T]([]byte(`
			{  
				"status": "ok",
				"msg": null,
				"data": {  
					"plan": "4-Payment",
					"amount": 50000,
					"interest": 0,
					"repayable": 50000,
					"schedule": [  
						{  
							"date": "2019-03-12",
							"amount": 12500
						},
						{  
							"date": "2019-04-12",
							"amount": 12500
						},
						{  
							"date": "2019-05-12",
							"amount": 12500
						},
						{  
							"date": "2019-06-12",
							"amount": 12500
						}
					]
				}
			}`))
	case "capture":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"token": "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
					"status": "completed",
					"deposit_captured": true
				}
			}`))
	case "invoice":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"token": "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
					"upload_status": "success"
				}
			}`))
	case "status":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"token": "aed3bd4e-c478-4d73-a6fa-3640a7155e4f",
					"status": "pending",
					"amount": 50000,
					"expires_at": "2022-05-24T19:28:06+01:00",
					"pa_ref": "testreference",
					"requires_invoice": true,
					"has_invoice": true,
					"last_accessed_at": "2025-11-12T12:00:00+00:00"
				}
			}`))
	case "account":
		return decodeResponseJSON[T]([]byte(`
			{
				"status": "ok",
				"msg": null,
				"data": {
					"legal_name": "Test Merchant",
					"display_name": "Test Merchant",
					"plans": [
						{
							"plan_id": 6,
							"name": "3-Payment",
							"instalments": 3,
							"deposit": true,
							"apr": 0,
							"frequency": "monthly",
							"min_amount": null,
							"max_amount": 500000,
							"commission_rate": "8.50",
							"commission_fixed_fee": null
						},
						{
							"plan_id": 1,
							"name": "4-Payment",
							"instalments": 4,
							"deposit": false,
							"apr": 5.5,
							"frequency": "monthly",
							"min_amount": 10000,
							"max_amount": 300000,
							"commission_rate": "0",
							"commission_fixed_fee": 5000
						}
					]
				}
			}`))
	default:
		panic("unrecognised endpoint " + endpoint)
	}
}
