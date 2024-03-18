package pasdk

var userCredentials PAAuth

// Initialises the SDK with your API credentials as well as the API URL
// you want to make requests to.
func Initialise(credentials PAAuth) {
	userCredentials = credentials
}
