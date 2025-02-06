package healthcheck

// Check is a healthcheck/readiness check.
type Check func() error

// StatusCheck holds the response the of the check and the error mapping
type StatusCheck struct {
	Name  string
	Error *string
}
