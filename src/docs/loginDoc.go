package docs

import (
	security "codestep/security"
)

// swagger:route POST /api/login login-tag idOfLoginEdnpoint
// Authenticate user
// responses:
//   200: loginResponse

// swagger:parameters idOfLoginEdnpoint
type loginRequestWrapper struct {
	// Login and password
	// in:body
	Body security.LoginRequest
}

// User authentication success
// swagger:response loginResponse
type loginResponseWrapper struct {
	// in:body
	Body security.LoginResponse
}
