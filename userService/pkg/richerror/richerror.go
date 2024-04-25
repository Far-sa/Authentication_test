package richerror

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindNotFound
	KindUnexpected
	KindForbidden
)

type Op string

type RichError struct {
	operation     Op
	internalError error
	message       string
	status        Kind
	meta          map[string]interface{}
}

func New(op Op) RichError {
	return RichError{operation: op}
}

func (r RichError) Error() string {
	return r.message
}

func (r RichError) WithErr(err error) RichError {
	r.internalError = err
	return r
}

func (r RichError) WithOp(op Op) RichError {
	r.operation = op
	return r
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) WithKind(status Kind) RichError {
	r.status = status
	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r
}

// * Recursive fn
func (r RichError) Status() Kind {
	if r.status != 0 {
		return r.status
	}

	re, ok := r.internalError.(RichError)
	if !ok {
		return 0
	}

	return re.Status()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	re, ok := r.internalError.(RichError)
	if !ok {
		return r.internalError.Error()
	}

	return re.Message()
}

// func New(args ...interface{}) RichError {
// 	r := RichError{}

// 	for _, arg := range args {
// 		switch arg.(type) {
// 		case string:
// 			r.message = arg.(string)
// 		case Op:
// 			r.operation = arg.(Op)
// 		case error:
// 			r.internalError = arg.(error)
// 		case Kind:
// 			r.status = arg.(Kind)
// 		case map[string]interface{}:
// 			r.meta = arg.(map[string]interface{})
// 		}

// 	}
// 	return r
// }

// func New(err error, operation, message string, status Kind, meta map[string]interface{}) RichError {
// 	return RichError{
// 		operation:     operation,
// 		internalError: err,
// 		message:       message,
// 		status:        status,
// 		meta:          meta,
// 	}
// }
