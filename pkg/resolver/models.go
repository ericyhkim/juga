package resolver

import (
	"errors"

	"github.com/ericyhkim/juga/pkg/models"
)

type ResolutionSource string

const (
	SourceAlias  ResolutionSource = "Alias"
	SourceCode   ResolutionSource = "Code"
	SourceCache  ResolutionSource = "Cache"
	SourceSearch ResolutionSource = "Search"
	SourceNone   ResolutionSource = "None"
)

type ResolutionStatus string

const (
	StatusSuccess   ResolutionStatus = "Success"
	StatusNotFound  ResolutionStatus = "NotFound"
	StatusAmbiguous ResolutionStatus = "Ambiguous"
)

type ResolutionResult struct {
	Input       string
	Code        string
	Name        string
	Source      ResolutionSource
	Status      ResolutionStatus
	IsAmbiguous bool
	Candidates  []models.Ticker
	Trace       string
	Error       error
}

var (
	ErrNotFound = errors.New("stock not found")
)
