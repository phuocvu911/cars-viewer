package cookies

const (
	// name under which the allowance is saved to the client. Update also go-backend/static/main.js if changed.
	ALLOW_TRACKING_COOKIE_NAME      string = "cars_viewer_allow_tracking"
	TRACKING_COOKIE_EXP_NOT_ALLOWED int    = 60 * 60 * 24 * 30  // Prompt tracking again in 30 days, by invalidating the cookie. Time measured in seconds.
	TRACKING_COOKIE_EXP_ALLOWED     int    = 60 * 60 * 24 * 365 // Prompt tracking again in 365 days ;) Time measured in seconds.

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
		Time measured in seconds.
	*/
	SHORT_TERM_USER_TRACKING_COOKIE_EXP int = 60 * 60 * 24 * 5
	DEFAULT_COOKIE_LEN                  int = 5 // 1 is equal to ABCD, 2 ABCD-EFGH, ...

	COOKIE_PATHS string = "/" // The client will send the cookies to these endpoints
)
