package util

import (
	"encoding/json"
	"fmt"
	"io"
)

type RedditPostResponse struct {
	Kind string `json:"kind"`
	Data struct {
		After     string      `json:"after"`
		Dist      int         `json:"dist"`
		Modhash   string      `json:"modhash"`
		GeoFilter interface{} `json:"geo_filter"`
		Children  []struct {
			Kind string `json:"kind"`
			Data struct {
				ApprovedAtUtc     interface{} `json:"approved_at_utc"`
				Subreddit         string      `json:"subreddit"`
				Selftext          string      `json:"selftext"`
				AuthorFullname    string      `json:"author_fullname"`
				Saved             bool        `json:"saved"`
				ModReasonTitle    interface{} `json:"mod_reason_title"`
				Gilded            int         `json:"gilded"`
				Clicked           bool        `json:"clicked"`
				Title             string      `json:"title"`
				LinkFlairRichtext []struct {
					E string `json:"e"`
					T string `json:"t,omitempty"`
					A string `json:"a,omitempty"`
					U string `json:"u,omitempty"`
				} `json:"link_flair_richtext"`
				SubredditNamePrefixed      string      `json:"subreddit_name_prefixed"`
				Hidden                     bool        `json:"hidden"`
				Pwls                       int         `json:"pwls"`
				LinkFlairCssClass          *string     `json:"link_flair_css_class"`
				Downs                      int         `json:"downs"`
				ThumbnailHeight            *int        `json:"thumbnail_height"`
				TopAwardedType             interface{} `json:"top_awarded_type"`
				HideScore                  bool        `json:"hide_score"`
				Name                       string      `json:"name"`
				Quarantine                 bool        `json:"quarantine"`
				LinkFlairTextColor         string      `json:"link_flair_text_color"`
				UpvoteRatio                float64     `json:"upvote_ratio"`
				AuthorFlairBackgroundColor string      `json:"author_flair_background_color"`
				Ups                        int         `json:"ups"`
				TotalAwardsReceived        int         `json:"total_awards_received"`
				MediaEmbed                 struct {
				} `json:"media_embed"`
				ThumbnailWidth        *int          `json:"thumbnail_width"`
				AuthorFlairTemplateId *string       `json:"author_flair_template_id"`
				IsOriginalContent     bool          `json:"is_original_content"`
				UserReports           []interface{} `json:"user_reports"`
				SecureMedia           interface{}   `json:"secure_media"`
				IsRedditMediaDomain   bool          `json:"is_reddit_media_domain"`
				IsMeta                bool          `json:"is_meta"`
				Category              interface{}   `json:"category"`
				SecureMediaEmbed      struct {
				} `json:"secure_media_embed"`
				LinkFlairText       *string     `json:"link_flair_text"`
				CanModPost          bool        `json:"can_mod_post"`
				Score               int         `json:"score"`
				ApprovedBy          interface{} `json:"approved_by"`
				IsCreatedFromAdsUi  bool        `json:"is_created_from_ads_ui"`
				AuthorPremium       bool        `json:"author_premium"`
				Thumbnail           string      `json:"thumbnail"`
				Edited              interface{} `json:"edited"`
				AuthorFlairCssClass *string     `json:"author_flair_css_class"`
				AuthorFlairRichtext []struct {
					E string `json:"e"`
					T string `json:"t"`
				} `json:"author_flair_richtext"`
				Gildings struct {
				} `json:"gildings"`
				PostHint            string      `json:"post_hint,omitempty"`
				ContentCategories   interface{} `json:"content_categories"`
				IsSelf              bool        `json:"is_self"`
				SubredditType       string      `json:"subreddit_type"`
				Created             float64     `json:"created"`
				LinkFlairType       string      `json:"link_flair_type"`
				Wls                 int         `json:"wls"`
				RemovedByCategory   interface{} `json:"removed_by_category"`
				BannedBy            interface{} `json:"banned_by"`
				AuthorFlairType     string      `json:"author_flair_type"`
				Domain              string      `json:"domain"`
				AllowLiveComments   bool        `json:"allow_live_comments"`
				SelftextHtml        *string     `json:"selftext_html"`
				Likes               interface{} `json:"likes"`
				SuggestedSort       interface{} `json:"suggested_sort"`
				BannedAtUtc         interface{} `json:"banned_at_utc"`
				UrlOverriddenByDest string      `json:"url_overridden_by_dest,omitempty"`
				ViewCount           interface{} `json:"view_count"`
				Archived            bool        `json:"archived"`
				NoFollow            bool        `json:"no_follow"`
				IsCrosspostable     bool        `json:"is_crosspostable"`
				Pinned              bool        `json:"pinned"`
				Over18              bool        `json:"over_18"`
				Preview             struct {
					Images []struct {
						Source struct {
							Url    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"source"`
						Resolutions []struct {
							Url    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"resolutions"`
						Variants struct {
						} `json:"variants"`
						Id string `json:"id"`
					} `json:"images"`
					Enabled bool `json:"enabled"`
				} `json:"preview,omitempty"`
				AllAwardings             []interface{} `json:"all_awardings"`
				Awarders                 []interface{} `json:"awarders"`
				MediaOnly                bool          `json:"media_only"`
				LinkFlairTemplateId      string        `json:"link_flair_template_id,omitempty"`
				CanGild                  bool          `json:"can_gild"`
				Spoiler                  bool          `json:"spoiler"`
				Locked                   bool          `json:"locked"`
				AuthorFlairText          string        `json:"author_flair_text"`
				TreatmentTags            []interface{} `json:"treatment_tags"`
				Visited                  bool          `json:"visited"`
				RemovedBy                interface{}   `json:"removed_by"`
				ModNote                  interface{}   `json:"mod_note"`
				Distinguished            interface{}   `json:"distinguished"`
				SubredditId              string        `json:"subreddit_id"`
				AuthorIsBlocked          bool          `json:"author_is_blocked"`
				ModReasonBy              interface{}   `json:"mod_reason_by"`
				NumReports               interface{}   `json:"num_reports"`
				RemovalReason            interface{}   `json:"removal_reason"`
				LinkFlairBackgroundColor string        `json:"link_flair_background_color"`
				Id                       string        `json:"id"`
				IsRobotIndexable         bool          `json:"is_robot_indexable"`
				ReportReasons            interface{}   `json:"report_reasons"`
				Author                   string        `json:"author"`
				DiscussionType           interface{}   `json:"discussion_type"`
				NumComments              int           `json:"num_comments"`
				SendReplies              bool          `json:"send_replies"`
				WhitelistStatus          string        `json:"whitelist_status"`
				ContestMode              bool          `json:"contest_mode"`
				ModReports               []interface{} `json:"mod_reports"`
				AuthorPatreonFlair       bool          `json:"author_patreon_flair"`
				AuthorFlairTextColor     string        `json:"author_flair_text_color"`
				Permalink                string        `json:"permalink"`
				ParentWhitelistStatus    string        `json:"parent_whitelist_status"`
				Stickied                 bool          `json:"stickied"`
				Url                      string        `json:"url"`
				SubredditSubscribers     int           `json:"subreddit_subscribers"`
				CreatedUtc               float64       `json:"created_utc"`
				NumCrossposts            int           `json:"num_crossposts"`
				Media                    interface{}   `json:"media"`
				IsVideo                  bool          `json:"is_video"`
			} `json:"data"`
		} `json:"children"`
		Before interface{} `json:"before"`
	} `json:"data"`
}

type RedditPost struct {
	ImageUrl string
	Title    string
	Url      string
}

func DecodeRedditPost(body *io.ReadCloser, res *RedditPost) error {
	r := RedditPostResponse{}
	err := json.NewDecoder(*body).Decode(&r)
	if err != nil {
		return fmt.Errorf("could not decode reddit response: %v", err)
	}

	if len(r.Data.Children) == 0 {
		return fmt.Errorf("no results returned")
	}

	if len(r.Data.Children[0].Data.Preview.Images) == 0 {
		return fmt.Errorf("post does not contain any photos")
	}

	if !r.Data.Children[0].Data.Preview.Enabled {
		return fmt.Errorf("preview disabled")
	}

	res.Title = r.Data.Children[0].Data.Title
	res.ImageUrl = r.Data.Children[0].Data.Preview.Images[0].Source.Url
	res.Url = r.Data.Children[0].Data.Url

	return nil
}
