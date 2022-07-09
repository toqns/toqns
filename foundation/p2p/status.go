package p2p

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusPaymentRequired     = 402
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// StatusText returns a default explanation for the provided status.
func StatusText(code int) string {
	switch code {
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not found"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""
	}
}
