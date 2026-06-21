package cookies

import (
	"net/http"
	"time"
)

// Return the long cookie without changing it. This is to avoid a bug where a client can have only short term cookie with permission.
func WriteLongCookieHeader(w http.ResponseWriter, cookie *http.Cookie) {

	if cookie != nil {
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}

}

// Creates a tracking allowance cookie to the request
// Call after prompting agreement for data collection.
func SetTrackingAllowanceAndLongtermIDCookies(w http.ResponseWriter, allow bool) {

	var allowed_str string
	expr := time.Now()

	if allow {
		allowed_str = "true"
		expr = expr.Add(TRACKING_COOKIE_EXP_ALLOWED)
	} else {
		allowed_str = "false"
		expr = expr.Add(TRACKING_COOKIE_EXP_NOT_ALLOWED)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     ALLOW_TRACKING_COOKIE_NAME,
		Value:    allowed_str,
		Path:     "/", // Pass to all paths
		Expires:  expr,
		SameSite: http.SameSiteLaxMode,
	})

	if !allow {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     LONG_USER_TRACKING_COOKIE_NAME,
		Value:    GenerateCookie(DEFAULT_COOKIE_LEN),
		Path:     "/", // Pass to all paths
		Expires:  expr,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

}

// Used for updating the short term tracking cookie
func UpdateShortCookie(w http.ResponseWriter, cookie *http.Cookie) {

	expr := time.Now().Add(SHORT_TERM_USER_TRACKING_COOKIE_EXP)

	if cookie == nil {
		cookie = &http.Cookie{Name: SHORT_USER_TRACKING_COOKIE_NAME, Expires: expr, Value: GenerateCookie(DEFAULT_COOKIE_LEN)}
	}

	if cookie.Value == "" {
		cookie.Value = GenerateCookie(DEFAULT_COOKIE_LEN)
	}

	cookie.Expires = expr
	cookie.Path = "/"

	http.SetCookie(w, cookie)
}
