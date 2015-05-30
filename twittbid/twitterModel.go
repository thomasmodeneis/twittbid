package twittbid

import (
	"net/http"
	"net/url"

	"github.com/ChimeraCoder/tokenbucket"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	_GET      = iota
	_Post     = iota
	baseURLV1 = "https://api.twitter.com/1"
	//BaseURL for twitter cfg
	BaseURL       = "https://api.twitter.com/1.1"
	uploadBaseURL = "https://upload.twitter.com/1.1"
)

//Tweet defines the twitt itself to enable transform from a json representation
type Tweet struct {
	Contributors         []Contributor          `json:"contributors"` // Not yet generally available to all, so hard to test
	Coordinates          *Coordinates           `json:"coordinates"`
	CreatedAt            string                 `json:"created_at"`
	Entities             Entities               `json:"entities"`
	FavoriteCount        int                    `json:"favorite_count"`
	Favorited            bool                   `json:"favorited"`
	FilterLevel          string                 `json:"filter_level"`
	ID                   int64                  `json:"id"`
	IDStr                string                 `json:"id_str"`
	InReplyToScreenName  string                 `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64                  `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string                 `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64                  `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string                 `json:"in_reply_to_user_id_str"`
	Lang                 string                 `json:"lang"`
	PossiblySensitive    bool                   `json:"possibly_sensitive"`
	RetweetCount         int                    `json:"retweet_count"`
	Retweeted            bool                   `json:"retweeted"`
	RetweetedStatus      *Tweet                 `json:"retweeted_status"`
	Source               string                 `json:"source"`
	Scopes               map[string]interface{} `json:"scopes"`
	Text                 string                 `json:"text"`
	Truncated            bool                   `json:"truncated"`
	User                 User                   `json:"user"`
	WithheldCopyright    bool                   `json:"withheld_copyright"`
	WithheldInCountries  []string               `json:"withheld_in_countries"`
	WithheldScope        string                 `json:"withheld_scope"`
}

//User define the twitter author user details
type User struct {
	ContributorsEnabled            bool     `json:"contributors_enabled"`
	CreatedAt                      string   `json:"created_at"`
	DefaultProfile                 bool     `json:"default_profile"`
	DefaultProfileImage            bool     `json:"default_profile_image"`
	Description                    string   `json:"description"`
	Entities                       Entities `json:"entities"`
	FavouritesCount                int      `json:"favourites_count"`
	FollowRequestSent              bool     `json:"follow_request_sent"`
	FollowersCount                 int      `json:"followers_count"`
	Following                      bool     `json:"following"`
	FriendsCount                   int      `json:"friends_count"`
	GeoEnabled                     bool     `json:"geo_enabled"`
	ID                             int64    `json:"id"`
	IDStr                          string   `json:"id_str"`
	IsTranslator                   bool     `json:"is_translator"`
	Lang                           string   `json:"lang"` // BCP-47 code of user defined language
	ListedCount                    int64    `json:"listed_count"`
	Location                       string   `json:"location"` // User defined location
	Name                           string   `json:"name"`
	Notifications                  bool     `json:"notifications"`
	ProfileBackgroundColor         string   `json:"profile_background_color"`
	ProfileBackgroundImageURL      string   `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHttps string   `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool     `json:"profile_background_tile"`
	ProfileBannerURL               string   `json:"profile_banner_url"`
	ProfileImageURL                string   `json:"profile_image_url"`
	ProfileImageURLHttps           string   `json:"profile_image_url_https"`
	ProfileLinkColor               string   `json:"profile_link_color"`
	ProfileSidebarBorderColor      string   `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string   `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string   `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool     `json:"profile_use_background_image"`
	Protected                      bool     `json:"protected"`
	ScreenName                     string   `json:"screen_name"`
	ShowAllInlineMedia             bool     `json:"show_all_inline_media"`
	Status                         *Tweet   `json:"status"` // Only included if the user is a friend
	StatusesCount                  int64    `json:"statuses_count"`
	TimeZone                       string   `json:"time_zone"`
	URL                            string   `json:"url"` // From UTC in seconds
	UtcOffset                      int      `json:"utc_offset"`
	Verified                       bool     `json:"verified"`
	WithheldInCountries            string   `json:"withheld_in_countries"`
	WithheldScope                  string   `json:"withheld_scope"`
}

//OauthClient defines oauth URLs
var OauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

//Contributor define the twitter contributor details
type Contributor struct {
	ID         int64  `json:"id"`
	IDStr      string `json:"id_str"`
	ScreenName string `json:"screen_name"`
}

//Coordinates define geo location
type Coordinates struct {
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
	Type        string     `json:"type"`
}

//SearchMetadata defines search metadata object
type SearchMetadata struct {
	CompletedIn   float32 `json:"completed_in"`
	MaxID         int64   `json:"max_id"`
	MaxIDString   string  `json:"max_id_str"`
	Query         string  `json:"query"`
	RefreshURL    string  `json:"refresh_url"`
	Count         int     `json:"count"`
	SinceID       int64   `json:"since_id"`
	SinceISString string  `json:"since_id_str"`
	NextResults   string  `json:"next_results"`
}

//SearchResponse defines the search response object with the twitter list and related metadata
type SearchResponse struct {
	Statuses []Tweet        `json:"statuses"`
	Metadata SearchMetadata `json:"search_metadata"`
}

//TwitterError json struct
type TwitterError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

//TwitterErrorResponse json struct
type TwitterErrorResponse struct {
	Errors []TwitterError `json:"errors"`
}

//URLEntity defines the url object
type URLEntity struct {
	Urls []struct {
		Indices     []int
		URL         string
		DisplayURL  string
		ExpandedURL string
	}
}

//Entities defines the entities object
type Entities struct {
	Hashtags []struct {
		Indices []int
		Text    string
	}
	Urls []struct {
		Indices     []int
		URL         string
		DisplayURL  string
		ExpandedURL string
	}
	URL          URLEntity
	UserMentions []struct {
		Name       string
		Indices    []int
		ScreenName string
		ID         int64
		IDstr      string
	}
	Media []struct {
		ID            int64
		IDstr         string
		MediaURL      string
		MediaURLhttps string
		URL           string
		DisplayURL    string
		ExpandedURL   string
		Sizes         MediaSizes
		Type          string
		Indices       []int
	}
}

//MediaSizes defines the sizes for media object
type MediaSizes struct {
	Medium MediaSize
	Thumb  MediaSize
	Small  MediaSize
	Large  MediaSize
}

//MediaSize defines media object
type MediaSize struct {
	W      int
	H      int
	Resize string
}

//TwitterAPI defines the twitter API fields eg: credentials, etc.
type TwitterAPI struct {
	Credentials          *oauth.Credentials
	queryQueue           chan Query
	bucket               *tokenbucket.Bucket
	returnRateLimitError bool

	HTTPClient *http.Client
}

//Query struct
type Query struct {
	url        string
	form       url.Values
	data       interface{}
	method     int
	responseCh chan Response
}

//Response struct
type Response struct {
	data interface{}
	err  error
}
