package cookies

import (
	"errors"
	"net/http"
	"time"
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

// Cookie hierarchy = allow -> long -> short
// if long is missing, generate short too
func ParseCookieCtx(r *http.Request) (CookieCtx, error) {

	track, err := AllowedToTrack(r)

	// If the allowance does not exist at all, return error instead of defaulting to false.
	if err != nil {
		return CookieCtx{AllowTracking: &http.Cookie{Name: ALLOW_TRACKING_COOKIE_NAME, Value: "false", Expires: time.Now().Add(TRACKING_COOKIE_EXP_NOT_ALLOWED)}, ShortCookie: &http.Cookie{}}, err
	}

	// If the cookie is found, but the value is not "true"
	if track.Value != "true" {
		return CookieCtx{AllowTracking: &http.Cookie{Name: ALLOW_TRACKING_COOKIE_NAME, Value: "false", Expires: time.Now().Add(TRACKING_COOKIE_EXP_NOT_ALLOWED)}}, nil
	}

	long_cookie, err := GetLongTrackingCookie(r)

	var short_cookie *http.Cookie

	if err != nil {
		// Return error IF the error IS NOT ErrNoCookie.
		// ErrNoCookie is expected.
		if !errors.Is(err, http.ErrNoCookie) {
			return CookieCtx{}, err
		}

		// Tracking allowed, Cookie not found = SetNewCookie true
		long_cookie = &http.Cookie{Name: LONG_USER_TRACKING_COOKIE_NAME, Value: GenerateCookie(DEFAULT_COOKIE_LEN), Expires: time.Now().Add(TRACKING_COOKIE_EXP_ALLOWED)}
		short_cookie = &http.Cookie{Name: SHORT_USER_TRACKING_COOKIE_NAME, Value: GenerateCookie(DEFAULT_COOKIE_LEN), Expires: time.Now().Add(SHORT_TERM_USER_TRACKING_COOKIE_EXP)}

		return CookieCtx{AllowTracking: track, LongCookie: long_cookie, ShortCookie: short_cookie}, nil
	}

	short_cookie, err = GetShortTrackingCookie(r)

	if err != nil {
		// Return error IF the error IS NOT ErrNoCookie
		if !errors.Is(err, http.ErrNoCookie) {
			return CookieCtx{}, err
		}
		// Tracking allowed, Cookie not found = SetNewCookie true
		short_cookie = &http.Cookie{Name: SHORT_USER_TRACKING_COOKIE_NAME, Value: GenerateCookie(DEFAULT_COOKIE_LEN), Expires: time.Now().Add(SHORT_TERM_USER_TRACKING_COOKIE_EXP)}
	}

	return CookieCtx{AllowTracking: track, LongCookie: long_cookie, ShortCookie: short_cookie}, nil
}
