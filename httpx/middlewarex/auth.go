package middlewarex

// import (
// 	"context"
// 	"net/http"
// 	"strings"
// )
//
// type AuthConfig struct {
// 	Authorizer
// 	AuthProvider
// 	Optional bool
// }
//
// func Auth(cfg AuthConfig) Func {
// 	return func(next http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			token, ok := cfg.Token(r)
// 			if ok && cfg.Authorize(token) {
// 				next(w, AuthorizeRequest(r))
// 			} else if cfg.Optional {
// 				next(w, r)
// 			} else {
// 				statusCode := http.StatusUnauthorized
// 				http.Error(w, http.StatusText(statusCode), statusCode)
// 			}
// 		}
// 	}
// }
//
// func AuthorizationHeaderAuth(optional bool, tokens ...string) Func {
// 	return Auth(AuthConfig{
// 		Authorizer:   NewTokenAuthorizer(tokens...),
// 		AuthProvider: NewAuthorizationHeaderAuthProvider(),
// 		Optional:     optional,
// 	})
// }
//
// //// optional auth helpers
//
// func IsAuthenticated(ctx context.Context) bool {
// 	return ctx.Value(authContextKey) != nil
// }
//
// func AuthorizeRequest(r *http.Request) *http.Request {
// 	ctx := context.WithValue(r.Context(), authContextKey, struct{}{})
// 	return r.WithContext(ctx)
// }
//
// //// Authorizers
//
// // Authorizer takes a token string and returns a bool.
// type Authorizer interface {
// 	Authorize(string) bool
// }
//
// func NewTokenAuthorizer(tokens ...string) Authorizer {
// 	tokenMap := make(map[string]struct{}, len(tokens))
// 	for _, token := range tokens {
// 		tokenMap[strings.TrimSpace(token)] = struct{}{}
// 	}
//
// 	return tokenAuthorizer{tokens: tokenMap}
// }
//
// type tokenAuthorizer struct {
// 	tokens map[string]struct{}
// }
//
// func (a tokenAuthorizer) Authorize(token string) bool {
// 	if token == "" {
// 		return false
// 	}
//
// 	_, ok := a.tokens[token]
// 	return ok
// }
//
// //// AuthProviders
//
// // AuthProvider takes in a request and returns the auth token and a bool
// // which will be true if the token was found.
// type AuthProvider interface {
// 	Token(*http.Request) (string, bool)
// }
//
// func NewHeaderAuthProvider(key, prefix string) AuthProvider {
// 	return HeaderAuthProvider{
// 		Key:    key,
// 		Prefix: prefix,
// 	}
// }
//
// func NewAuthorizationHeaderAuthProvider() AuthProvider {
// 	return NewHeaderAuthProvider("Authorization", "")
// }
//
// type HeaderAuthProvider struct {
// 	Key    string
// 	Prefix string
// }
//
// func (p HeaderAuthProvider) Token(r *http.Request) (string, bool) {
// 	headers, ok := r.Header[p.Key]
// 	if !ok || len(headers) == 0 {
// 		return "", false
// 	}
//
// 	header := headers[0]
// 	if p.Prefix != "" {
// 		header = strings.TrimPrefix(header, p.Prefix)
// 		header = strings.TrimSpace(header)
// 	}
//
// 	return header, true
// }
