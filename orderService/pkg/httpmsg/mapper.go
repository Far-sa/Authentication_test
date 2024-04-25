package httpmsg

import (
	"net/http"
	"user-svc/pkg/errmsg"
	"user-svc/pkg/richerror"
)

// Error Interpreter err
func Error(err error) (message string, code int) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()
		code := MapKindToHttpStatusCode(re.Status())
		if code > 500 {
			msg = errmsg.ErrorMsgInternalServerError
		}
		return msg, code
		//return re.Message(), MapKindToHttpStatusCode(re.Status())

	default:
		return err.Error(), http.StatusBadRequest
	}
}

func MapKindToHttpStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
