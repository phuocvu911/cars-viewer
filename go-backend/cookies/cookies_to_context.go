package cookies

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type CookieCtxKey struct{}

type CookieCtx struct {
	AllowTracking *http.Cookie
	ShortCookie   *http.Cookie
	LongCookie    *http.Cookie
}

// Parses cookies from the request headers. Cookie Context is ONLY used for tracking purposes.
// Routes not outputting user related content do not need the middleware.
// If request needs to be logged on user level or output is determined by the previous behaviour, then wrap the request.
// Middleware handles updating short term cookie.
//
//	AddCookieContext(http.HandlerFunc(Handler))
//
// In Handler, access the context by using:
//
//	r.Context()
//
// Cookies context is the following:
//
//	type CookieCtx struct {
//		AllowTracking bool
//		Cookie        *string
//		SetNewCookie  bool
//	}
func AddCookieContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieCtx, err := ParseCookieCtx(r)

		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("ParseCookieCtx error: %v", err)
				http.Error(w, "Middleware failed. ", http.StatusInternalServerError)
				return
			}
			// Accept ErrNoCookie
		}

		// Add the context to the current context by using context.WithValue()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CookieCtxKey{}, cookieCtx)))

	})
}
