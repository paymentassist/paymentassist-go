package pasdk

import (
	"strings"
	"testing"
	"time"
)

// Run the suite of integration tests, if enabled.
func Test_IntegrationTests(t *testing.T) {
	if !shouldRunIntegrationTests() {
		return
	}

	Initialise(PAAuth{
		APISecret: getTestAPISecret(),
		APIKey:    getTestAPIKey(),
		APIURL:    "https://api.demo.payassi.st/",
	})

	accountResponse := testAccount(t)
	testPlan(t, *accountResponse)
	testPreapproval(t)
	beginResponse := testBegin(t)
	testStatus(t, beginResponse.ApplicationToken)
	testUpdate(t, beginResponse.ApplicationToken)
	testCapture(t, beginResponse.ApplicationToken)
	testInvoice(t, beginResponse.ApplicationToken)
}

func testUpdate(t *testing.T, token string) {
	// Test only updating some fields.
	amount := 80000

	request := UpdateRequest{
		ApplicationToken: token,
		Amount:           &amount,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if *response.Amount != 80000 {
		t.Error()
	}
	if response.ExpiresIn != nil {
		t.Error()
	}
	if response.OrderID != nil {
		t.Error()
	}

	// Test updating most fields. We can't easily test updating order ID
	// because that requires a completed application.
	amount = 70000
	expiresIn := 60 * 10

	request = UpdateRequest{
		ApplicationToken: token,
		Amount:           &amount,
		ExpiresIn:        &expiresIn,
	}

	response, err = request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if *response.Amount != 70000 {
		t.Error()
	}
	if *response.ExpiresIn != 600 {
		t.Error()
	}
	if response.OrderID != nil {
		t.Error()
	}
}

func testStatus(t *testing.T, token string) {
	request := StatusRequest{
		ApplicationToken: token,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.Amount != 100000 {
		t.Error()
	}
	if response.ExpiresAt.Before(time.Now()) ||
		response.ExpiresAt.After(time.Now().Add(time.Hour*25)) {
		t.Error(response.ExpiresAt)
	}
	if response.ApplicationToken != token {
		t.Error()
	}
	if response.HasInvoice != false {
		t.Error()
	}
	if response.RequriesInvoice != false {
		t.Error()
	}
	if response.Status != "pending" {
		t.Error(response.Status)
	}
}

func testPlan(t *testing.T, accountInfo AccountResponse) {
	request := PlanRequest{
		Amount: 50000,
		PlanID: &accountInfo.Plans[0].ID,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if response.Amount != 50000 {
		t.Error()
	}
	if response.Interest != 0 {
		t.Error()
	}
	if response.PlanName != accountInfo.Plans[0].Name {
		t.Error()
	}
	if response.TotalRepayable != 50000 {
		t.Error()
	}

	if len(response.PaymentSchedule) != 4 {
		t.Error()
	}
	if response.PaymentSchedule[3].Amount <= 0 {
		t.Error()
	}

	if response.PaymentSchedule[3].Date.Before(time.Now()) {
		t.Error()
	}
	if response.PaymentSchedule[3].Date.After(time.Now().AddDate(0, 5, 0)) {
		t.Error()
	}
}

func testPreapproval(t *testing.T) {
	request := PreapprovalRequest{
		CustomerFirstName: "Test",
		CustomerLastName:  "Testington",
		CustomerAddress1:  "Test House",
		CustomerPostcode:  "TEST TES",
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if !response.Approved {
		t.Error()
	}
}

func testCapture(t *testing.T, token string) {
	request := CaptureRequest{
		ApplicationToken: token,
	}

	_, err := request.Fetch()

	// This is the closest we can get to testing it because only a completed application
	// can be captured.
	if !strings.Contains(err.Error(), `{"status":"error","msg":"Application is not awaiting capture","data":[]}`) {
		t.Error()
	}
}

func testInvoice(t *testing.T, token string) {
	request := InvoiceRequest{
		ApplicationToken: token,
		FileType:         "txt",
		FileData:         []byte("Test invoice"),
	}

	response, err := request.Fetch()

	// Only a completed application can be invoiced so we are expecting this to fail.
	if !strings.Contains(err.Error(), "Application is not yet completed") {
		t.Error(err.Error())
	}
	if response != nil {
		t.Error()
	}
}

func testBegin(t *testing.T) *BeginResponse {
	falseValue := false

	request := BeginRequest{
		OrderID:           getRandomID(),
		Amount:            100000,
		CustomerFirstName: "Test",
		CustomerLastName:  "Testington",
		CustomerAddress1:  "Test House",
		CustomerPostcode:  "TEST TES",
		EnableAutoCapture: &falseValue,
	}

	response, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if len(response.ApplicationToken) != 36 {
		t.Error()
	}
	if len(response.ContinuationURL) < 10 {
		t.Error()
	}

	return response
}

func testAccount(t *testing.T) *AccountResponse {
	request := AccountRequest{}

	accountResponse, err := request.Fetch()

	if err != nil {
		t.Error(err)
	}

	if len(accountResponse.DisplayName) == 0 {
		t.Error()
	}
	if len(accountResponse.LegalName) == 0 {
		t.Error()
	}
	if len(accountResponse.Plans[0].APR) == 0 {
		t.Error()
	}
	if len(accountResponse.Plans[0].CommissionRate) == 0 {
		t.Error()
	}
	if len(accountResponse.Plans[0].Frequency) == 0 {
		t.Error()
	}
	if len(accountResponse.Plans[0].Name) == 0 {
		t.Error()
	}
	if accountResponse.Plans[0].DepositRequired == false {
		t.Error()
	}
	if accountResponse.Plans[0].ID == 0 {
		t.Error()
	}
	if accountResponse.Plans[0].Instalments == 0 {
		t.Error()
	}
	if accountResponse.Plans[0].MaxAmount == nil ||
		*accountResponse.Plans[0].MaxAmount <= 0 {
		t.Error()
	}

	return accountResponse
}
