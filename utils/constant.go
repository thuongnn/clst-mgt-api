package utils

import "time"

const (
	StatusSuccess = 0
	StatusError   = 1
	StatusPending = 2
	StatusUnknown = 3

	TimeoutScan = time.Second

	StatusSuccessScan = "success"
	StatusErrorScan   = "error"

	BasicAuth  = "basic"
	Oauth2Auth = "oauth2"

	AdminRole = "admin"
	UserRole  = "user"
)
