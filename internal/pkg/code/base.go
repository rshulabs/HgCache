package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ./error_code_generated.md

// common
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001
	// ErrUnknown - 500: Internal server error.
	ErrUnknown
	// ErrValidation - 400: Validation failed.
	ErrValidation
	// ErrNotFound - 404: Page not found.
	ErrNotFound
	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid
	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind
)

// database
const (
	//ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

// jwt authorization
const (
	// ErrEncrypt - 401: Error occurred while encrypting the user password.
	ErrEncrypt int = iota + 100201
	// ErrSignatureInvalid - 401: Signature invalid.
	ErrSignatureInvalid
	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader
	// ErrMissingHeader - 401: Missing authorization header.
	ErrMissingHeader
	// ErrPasswordIncorrect - 401: Password incorrect.
	ErrPasswordIncorrect
	// ErrPermissionDenied - 403: Permission denied.
	ErrPermissionDenied
	// ErrBlackListCheck - 401: Black list check failed.
	ErrBlackListCheck
	// ErrGuardTokenCheck - 401: Guard token check failed.
	ErrGuardTokenCheck
)

// encode/decode
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301
	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed
	// ErrInvalidJSON - 500: Invalid json data.
	ErrInvalidJSON
	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON
	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON
	// ErrInvalidYaml - 500: Invalid yaml data.
	ErrInvalidYaml
	// ErrEncodingYaml - 500: YAML data could not be encoded.
	ErrEncodingYaml
	// ErrDecodingYAML - 500: YAML data could not be decoded.
	ErrDecodingYAML
)
