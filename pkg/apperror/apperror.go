package apperror

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm/logger"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

var (
	// ------------------------
	// Generic errors
	// ------------------------
	ErrInternalServer = errors.New("internal server error") // 500
	ErrUnknown        = errors.New("unknown error")         // 500
	ErrTimeout        = errors.New("timeout")               // 504
	ErrUnauthorized   = errors.New("unauthorized")          // 401
	ErrForbidden      = errors.New("forbidden")             // 403
	ErrNotImplemented = errors.New("not implemented")       // 501

	// ------------------------
	// GORM errors
	// ------------------------
	ErrRecordNotFound                = logger.ErrRecordNotFound                                          // 404
	ErrInvalidTransaction            = errors.New("invalid transaction")                                 // 400
	ErrMissingWhereClause            = errors.New("WHERE conditions required")                           // 400
	ErrUnsupportedRelation           = errors.New("unsupported relations")                               // 400
	ErrPrimaryKeyRequired            = errors.New("primary key required")                                // 400
	ErrModelValueRequired            = errors.New("model value required")                                // 400
	ErrModelAccessibleFieldsRequired = errors.New("model accessible fields required")                    // 400
	ErrSubQueryRequired              = errors.New("sub query required")                                  // 400
	ErrUnsupportData                 = errors.New("unsupported data")                                    // 400
	ErrUnsupportedDriver             = errors.New("unsupported driver")                                  // 400
	ErrRegistered                    = errors.New("registered")                                          // 409
	ErrInvalidField                  = errors.New("invalid field")                                       // 400
	ErrEmptySlice                    = errors.New("empty slice found")                                   // 400
	ErrDryRunModeUnsupported         = errors.New("dry run mode unsupported")                            // 400
	ErrInvalidDB                     = errors.New("invalid db")                                          // 400
	ErrInvalidValue                  = errors.New("invalid value, should be pointer to struct or slice") // 400
	ErrInvalidValueOfLength          = errors.New("invalid association values, length doesn't match")    // 400
	ErrPreloadNotAllowed             = errors.New("preload is not allowed when count is used")           // 400
	ErrDuplicatedKey                 = errors.New("duplicated key not allowed")                          // 409
	ErrForeignKeyViolated            = errors.New("violates foreign key constraint")                     // 409
	ErrCheckConstraintViolated       = errors.New("violates check constraint")                           // 409

	// ------------------------
	// Validation errors
	// ------------------------
	ErrInvalidData   = errors.New("invalid data")           // 400
	ErrInvalidID     = errors.New("invalid id")             // 400
	ErrRequiredField = errors.New("required field missing") // 400
	ErrInvalidFormat = errors.New("invalid format")         // 400
	ErrOutOfRange    = errors.New("value out of range")     // 400
	ErrUnprocessable = errors.New("unprocessable entity")   // 422

	// ------------------------
	// Business logic / domain-specific errors
	// ------------------------
	ErrAlreadyExists   = errors.New("already exists")   // 409
	ErrNotAvailable    = errors.New("not available")    // 409
	ErrLimitExceeded   = errors.New("limit exceeded")   // 429
	ErrOperationDenied = errors.New("operation denied") // 403

	// ------------------------
	// Other errors
	// ------------------------
	ErrConflict         = errors.New("conflict")            // 409
	ErrDependencyFail   = errors.New("dependency failure")  // 502
	ErrTransactionAbort = errors.New("transaction aborted") // 500
)

// StatusCode maps errors to Fiber HTTP status codes
func StatusCode(err error) int {
	switch {
	// Generic
	case errors.Is(err, ErrInternalServer), errors.Is(err, ErrUnknown), errors.Is(err, ErrTransactionAbort):
		return fiber.StatusInternalServerError
	case errors.Is(err, ErrTimeout):
		return fiber.StatusGatewayTimeout
	case errors.Is(err, ErrUnauthorized):
		return fiber.StatusUnauthorized
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrOperationDenied):
		return fiber.StatusForbidden
	case errors.Is(err, ErrNotImplemented):
		return fiber.StatusNotImplemented

	// Database / GORM errors
	case errors.Is(err, ErrRecordNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, ErrDuplicatedKey), errors.Is(err, ErrConflict), errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrNotAvailable):
		return fiber.StatusConflict
	case errors.Is(err, ErrDependencyFail):
		return fiber.StatusBadGateway
	case errors.Is(err, ErrInvalidTransaction), errors.Is(err, ErrMissingWhereClause),
		errors.Is(err, ErrUnsupportedRelation), errors.Is(err, ErrPrimaryKeyRequired),
		errors.Is(err, ErrModelValueRequired), errors.Is(err, ErrModelAccessibleFieldsRequired),
		errors.Is(err, ErrSubQueryRequired), errors.Is(err, ErrUnsupportData),
		errors.Is(err, ErrUnsupportedDriver), errors.Is(err, ErrEmptySlice),
		errors.Is(err, ErrDryRunModeUnsupported), errors.Is(err, ErrPreloadNotAllowed),
		errors.Is(err, ErrForeignKeyViolated), errors.Is(err, ErrCheckConstraintViolated):
		return fiber.StatusBadRequest

	// Validation / business logic
	case errors.Is(err, ErrInvalidData), errors.Is(err, ErrInvalidID), errors.Is(err, ErrRequiredField),
		errors.Is(err, ErrInvalidFormat), errors.Is(err, ErrOutOfRange), errors.Is(err, ErrInvalidValue),
		errors.Is(err, ErrInvalidValueOfLength), errors.Is(err, ErrInvalidField):
		return fiber.StatusBadRequest
	case errors.Is(err, ErrUnprocessable):
		return fiber.StatusUnprocessableEntity
	case errors.Is(err, ErrLimitExceeded):
		return fiber.StatusTooManyRequests

	// Default
	default:
		return fiber.StatusInternalServerError
	}
}

// GRPCCode maps errors to gRPC status codes
func GRPCCode(err error) codes.Code {
	switch {
	// Generic
	case errors.Is(err, ErrInternalServer), errors.Is(err, ErrUnknown), errors.Is(err, ErrTransactionAbort):
		return codes.Internal
	case errors.Is(err, ErrTimeout):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrUnauthorized):
		return codes.Unauthenticated
	case errors.Is(err, ErrForbidden), errors.Is(err, ErrOperationDenied):
		return codes.PermissionDenied
	case errors.Is(err, ErrNotImplemented):
		return codes.Unimplemented

	// Database / GORM errors
	case errors.Is(err, ErrRecordNotFound):
		return codes.NotFound
	case errors.Is(err, ErrDuplicatedKey), errors.Is(err, ErrConflict), errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrNotAvailable):
		return codes.AlreadyExists
	case errors.Is(err, ErrDependencyFail):
		return codes.FailedPrecondition
	case errors.Is(err, ErrInvalidTransaction), errors.Is(err, ErrMissingWhereClause),
		errors.Is(err, ErrUnsupportedRelation), errors.Is(err, ErrPrimaryKeyRequired),
		errors.Is(err, ErrModelValueRequired), errors.Is(err, ErrModelAccessibleFieldsRequired),
		errors.Is(err, ErrSubQueryRequired), errors.Is(err, ErrUnsupportData),
		errors.Is(err, ErrUnsupportedDriver), errors.Is(err, ErrEmptySlice),
		errors.Is(err, ErrDryRunModeUnsupported), errors.Is(err, ErrPreloadNotAllowed),
		errors.Is(err, ErrForeignKeyViolated), errors.Is(err, ErrCheckConstraintViolated),
		errors.Is(err, ErrInvalidData), errors.Is(err, ErrInvalidID), errors.Is(err, ErrRequiredField),
		errors.Is(err, ErrInvalidFormat), errors.Is(err, ErrOutOfRange), errors.Is(err, ErrInvalidValue),
		errors.Is(err, ErrInvalidValueOfLength), errors.Is(err, ErrInvalidField):
		return codes.InvalidArgument

	// Validation / business logic
	case errors.Is(err, ErrUnprocessable):
		return codes.FailedPrecondition
	case errors.Is(err, ErrLimitExceeded):
		return codes.ResourceExhausted

	// Default
	default:
		return codes.Unknown
	}
}
