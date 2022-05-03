package middlewarex

import "net/http"

type Func func(server http.Handler) http.Handler
