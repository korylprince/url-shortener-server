package httpapi

type contextKey int

const (
	contextKeyUser contextKey = iota
	contextKeyLogData
)
