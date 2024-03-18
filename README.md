# paymentassist-go

The official Go SDK for the Payment Assist Merchant API.

## Dependencies

Officially, Go v1.19 and beyond is supported, but the SDK might also work on v1.18.

## Installation

`go get github.com/paymentassist/paymentassist-go`

```
import (
	pasdk "github.com/paymentassist/paymentassist-go"
)
```

## Workflow

![Payment Assist API Workflow](https://raw.githubusercontent.com/paymentassist/paymentassist-php/master/api-workflow.png "API Workflow")

## Usage

The full API reference can be found here: https://api-docs.payment-assist.co.uk/.

To use this SDK, first start by initialising it with the `Initialise` function, which takes your API credentials as well as the PaymentAssist API URL you want to make requests to.

```
pasdk.Initialise(pasdk.AuthInfo{
	APIKey:    "my_api_key",
	APISecret: "my_api_secret",
	APIURL: "https://api.demo.payassi.st/",
})
```

Note that it is not recommended to hard-code your API credentials like in the above example, this is just for illustration purposes.

After this, you can create a request object for the action you want to perform, followed by calling the `Fetch()` method on it. `Fetch()` returns a response object and an error.

If an error is returned, the request was unsucessful and the response object will be `nil`. The error is a custom type that contains detailed information about what happened. Of note are the fields `IsRequestRefusedError`, `IsValidationFailedError` and `IsUnexpectedError`. In the case of failure you may want to use these to decide whether or not to retry the request. However, you don't have to use these, and there is no harm in retrying all errors. See the code comments for more information on what these error types mean.

Note that `InvoiceRequest` and `CaptureRequest` may return a response and no error even if the request was unsuccessful; specific error data for these is provided in the response.

Example:

```
request := pasdk.AccountRequest{}
accountResponse, err := request.Fetch()

if err != nil {
    fmt.Println("There was an error: "+err.Error())
	return
}

// Print the dealer's display name.
fmt.Println(accountResponse.DisplayName)
```

The following actions are available:

| Action | Description |
|--------|-------------|
| __AccountRequest__ | Returns information about the dealer's account. |
| __PlanRequest__ | Returns what the repayments would look like under a hypothetical repayment plan. |
| __PreapprovalRequest__ | Checks whether a customer would pass the basic pre-approval checks. |
| __BeginRequest__ | Begins an application. |
| __StatusRequest__ | Returns information about an ongoing application. |
| __UpdateRequest__ | Updates an existing application. |
| __CaptureRequest__ | Finalises an application that's in pending_capture state (used only when auto-capture is disabled). |
| __InvoiceRequest__ | Uploads an invoice for a completed application. |

## Notes


As virtually all requests to the API should return immediately (apart from /capture, which can take a few seconds to process deposits), there is currently no support for cancelling an ongoing request. There is a hard-coded timeout of 30 seconds per request which should be sufficient in all scenarios.

## Support

For technical support, please email [itsupport@payment-assist.co.uk](mailto:itsupport@payment-assist.co.uk).

If you encounter any issues or find that a particular part of the SDK isn't meeting your requirements, feel free to contact support and we will do our best to accommodate where we can.