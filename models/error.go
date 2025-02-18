package models

// ErrorResponse represents the error response format
// @Description Standard error response format for all API errors
// @Schema
type ErrorResponse struct {
	// Message contains a human-readable error message
	// @example "Invalid product ID"
	Message string `json:"message"`

	// Details contains additional information (optional)
	// @example "The product ID provided was not an integer."
	Details string `json:"details,omitempty"`
}
