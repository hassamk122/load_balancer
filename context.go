package main

import "net/http"

type ContextKey int

const (
	Attempts ContextKey = iota
	Retry
)

// both functions allows us to extract attempts
// and retries from req context
func GetAttemptsFromContext(req *http.Request) int {
	attempts, ok := req.Context().Value(Attempts).(int)
	if ok {
		return attempts
	}
	return 0
}

func GetRetryFromContext(req *http.Request) int {
	retry, ok := req.Context().Value(Retry).(int)
	if ok {
		return retry
	}
	return 0
}
