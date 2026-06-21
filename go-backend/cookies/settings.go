package cookies

import "time"

const (
	// name under which the allowance is saved to the client. Update also go-backend/static/main.js if changed.
	ALLOW_TRACKING_COOKIE_NAME      string        = "cars_viewer_allow_tracking"
	TRACKING_COOKIE_EXP_NOT_ALLOWED time.Duration = 30 * 24 * time.Hour  // Prompt tracking again in 30 days, by invalidating the cookie.
	TRACKING_COOKIE_EXP_ALLOWED     time.Duration = 365 * 24 * time.Hour // Prompt tracking again in 365 days ;)

	/*
		Long term tracking cookie. It's supposed to last as long as the permission to track lasts.
		It will keep track of the user and will be logged with short term cookie.
		Long term cookie serves the purpose to track user over trends. Short term cookie is where the recommendations are aggregated from.
	*/
	LONG_USER_TRACKING_COOKIE_NAME string = "cars_viewer_long_term_id"
	/*
		Short term cookie lasts just a while. If the user is inactive for more than SHORT_TERM_USER_TRACKING_COOKIE_EXP,
		tracking will start fresh again. Long term cookie makes it possible to track users preferences over time.
	*/
	SHORT_USER_TRACKING_COOKIE_NAME string = "cars_viewer_short_term_id"

	/*
		Essentially new short term cookie will be created to the user after this time period of inactivity.
		Intrests change over time so you don't want to aggregate whole year analytics into a single data set.
		Recommendations will be based on this user identity. The long term token can be used to identify trends over time.
	*/
	SHORT_TERM_USER_TRACKING_COOKIE_EXP time.Duration = 3 * 24 * time.Hour
	DEFAULT_COOKIE_LEN                  int           = 5
)
