package cookies

import (
	"net/http"
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

	tracking_allowed_cookie := GenerateAllowTrackCookie(allow)

	http.SetCookie(w, tracking_allowed_cookie)

	// Its a bit dumb here but it works
	if tracking_allowed_cookie.Value == "true" {
		http.SetCookie(w, GenerateLongCookie(tracking_allowed_cookie))
	}

}
