package cookies

import (
	"net/http"
)

// read permission for tracking
func AllowedToTrack(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(ALLOW_TRACKING_COOKIE_NAME)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

// Return error IF the token does not exist
func GetShortTrackingCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(SHORT_USER_TRACKING_COOKIE_NAME)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

// Return error IF the token does not exist
func GetLongTrackingCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie(LONG_USER_TRACKING_COOKIE_NAME)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func CookieNameOk(s string) bool {

	if s == ALLOW_TRACKING_COOKIE_NAME {
		return true
	}

	if s == LONG_USER_TRACKING_COOKIE_NAME {
		return true
	}

	if s == SHORT_USER_TRACKING_COOKIE_NAME {
		return true
	}

	return false
}

func CookieIsOk(cookie *http.Cookie) bool {

	if cookie == nil {
		return false
	}

	if cookie.Value == "" || cookie.Name == "" {
		return false
	}

	if !CookieNameOk(cookie.Name) {
		return false
	}

	return true
}

// Recreate cookieCtx and revoke all cookies by setting .Expires to the past
func SetDeleteAllCookiesHeader(w http.ResponseWriter) {
	// Generate all new cookies to make sure they exists without complicated logic
	var cookies CookieCtx

	cookies.AllowTracking = GenerateAllowTrackCookie(false)
	cookies.LongCookie = GenerateLongCookie(cookies.AllowTracking)
	cookies.ShortCookie = GenerateShortCookie()

	// Revoke the cookies by expiring them
	cookies.AllowTracking.MaxAge = -1
	cookies.LongCookie.MaxAge = -1
	cookies.ShortCookie.MaxAge = -1

	http.SetCookie(w, cookies.AllowTracking)
	http.SetCookie(w, cookies.LongCookie)
	http.SetCookie(w, cookies.ShortCookie)

}

// Only used if the Allow track is permitted
func GenerateLongCookie(allow_track *http.Cookie) *http.Cookie {

	return &http.Cookie{Name: LONG_USER_TRACKING_COOKIE_NAME,
		Value:    GenerateCookie(DEFAULT_COOKIE_LEN),
		HttpOnly: true,
		Secure:   true,
		MaxAge:   allow_track.MaxAge,
		Path:     COOKIE_PATHS,
	}
}

func GenerateShortCookie() *http.Cookie {
	return &http.Cookie{Name: SHORT_USER_TRACKING_COOKIE_NAME,
		Value:    GenerateCookie(DEFAULT_COOKIE_LEN),
		HttpOnly: true,
		Secure:   true,
		MaxAge:   SHORT_TERM_USER_TRACKING_COOKIE_EXP,
		Path:     COOKIE_PATHS,
	}

}

// Pass in a short life cookie and the function will check if all is good.
// If it is, its life is extended. If not the cookie is regenerated.
func WriteCookieLifeExtension(short_cookie *http.Cookie) *http.Cookie {

	if CookieIsOk(short_cookie) {

		short_cookie.MaxAge = SHORT_TERM_USER_TRACKING_COOKIE_EXP
		short_cookie.HttpOnly = true
		short_cookie.Path = COOKIE_PATHS
		return short_cookie

	}

	return &http.Cookie{Name: SHORT_USER_TRACKING_COOKIE_NAME,
		Value:    GenerateCookie(DEFAULT_COOKIE_LEN),
		HttpOnly: true,
		MaxAge:   SHORT_TERM_USER_TRACKING_COOKIE_EXP,
		Path:     COOKIE_PATHS,
	}

}

func GenerateAllowTrackCookie(allow bool) *http.Cookie {

	if allow {
		return &http.Cookie{Name: ALLOW_TRACKING_COOKIE_NAME,
			Value:    "true",
			HttpOnly: false,
			Path:     "/", // Use "/ path for allowance"
			MaxAge:   TRACKING_COOKIE_EXP_ALLOWED,
		}
	}

	return &http.Cookie{Name: ALLOW_TRACKING_COOKIE_NAME,
		Value:    "false",
		HttpOnly: false,
		Path:     "/", // Use "/ path for allowance"
		MaxAge:   TRACKING_COOKIE_EXP_NOT_ALLOWED,
	}

}

// Cookie hierarchy = allow -> long -> short
// if long is missing, generate short too
func ParseCookieCtx(r *http.Request) (CookieCtx, error) {

	// --------
	// Tracking
	// --------
	track, _ := AllowedToTrack(r)

	if !CookieIsOk(track) {
		return CookieCtx{AllowTracking: GenerateAllowTrackCookie(false), DeleteAllCookies: true}, nil
	}

	if track.Value != "true" && track.Value != "false" {
		return CookieCtx{AllowTracking: GenerateAllowTrackCookie(false), DeleteAllCookies: true}, nil
	}

	if track.Value == "false" {
		return CookieCtx{AllowTracking: GenerateAllowTrackCookie(false)}, nil
	}

	// --------
	// Long cookie
	// --------

	ctx := CookieCtx{AllowTracking: GenerateAllowTrackCookie(true)}

	long_cookie, _ := GetLongTrackingCookie(r)

	if !CookieIsOk(long_cookie) {
		ctx.LongCookie = GenerateLongCookie(ctx.AllowTracking)
		ctx.ReturnAllCookies = true
	} else {
		ctx.LongCookie = long_cookie
	}

	short_cookie, _ := GetShortTrackingCookie(r)

	if !CookieIsOk(short_cookie) {
		ctx.ShortCookie = GenerateShortCookie()
	} else {
		ctx.ShortCookie = WriteCookieLifeExtension(short_cookie)
	}

	return ctx, nil
}
