package resolver

import "errors"

// ResolutionSource indicates how a stock code was found.
type ResolutionSource string

const (
	SourceAlias  ResolutionSource = "Alias"
	SourceCode   ResolutionSource = "Code"
	SourceCache  ResolutionSource = "Cache"
	SourceSearch ResolutionSource = "Search"
	SourceNone   ResolutionSource = "None"
)

// ResolutionStatus represents the outcome of a resolution attempt.
type ResolutionStatus string

const (
	StatusSuccess   ResolutionStatus = "Success"
	StatusNotFound  ResolutionStatus = "NotFound"
	StatusAmbiguous ResolutionStatus = "Ambiguous"
)

type ResolutionResult struct {
	Input  string
	Code   string
	Name   string
	Source ResolutionSource
	Status ResolutionStatus
	Error  error
}

// Common errors for programmatic checking
var (
	ErrNotFound = errors.New("stock not found")
)
