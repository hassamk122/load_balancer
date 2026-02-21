package main

import "net/http"

const (
	Attempts int = iota
	Retry
)

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
