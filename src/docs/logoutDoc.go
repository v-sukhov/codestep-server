package docs

import (
	security "codestep/security"
)

// swagger:route POST /logout logout-tag idOfLogoutEdnpoint
// Logout user
// responses:
//   200: logoutResponse

// swagger:parameters idOfLogoutEdnpoint
type logoutRequestWrapper struct {
	// Empty request
	// in:body
	Body security.LogoutRequest
}

// User logout success
// swagger:response logoutResponse
type logoutResponseWrapper struct {
	// in:body
	Body security.LogoutResponse
}
