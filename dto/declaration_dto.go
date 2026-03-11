package dto

type DeclarationRequest struct {
	TraceID             string
	UndertakingName     string
	UndertakingDate     string
	UndertakingPlace    string
	UndertakingAccepted bool
}
