package middleware

import "net/http"

type Func func(server http.HandlerFunc) http.HandlerFunc
