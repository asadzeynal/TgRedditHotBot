package rdClient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func DecodeRedditResponse(body *io.ReadCloser) ([]*RedditPost, error) {
	r := RedditPostResponse{}
	err := json.NewDecoder(*body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("could not decode reddit response: %v", err)
	}

	if len(r.Data.Children) == 0 {
		return nil, fmt.Errorf("no results returned")
	}

	res := make([]*RedditPost, 0, len(r.Data.Children))
	for i := range r.Data.Children {
		child := r.Data.Children[i]
		var post RedditPost

		if len(child.Data.Preview.Images) > 0 && child.Data.Preview.Enabled {
			if child.Data.Preview.Images[0].Variants.Gif.Source.Url != "" {
				post.ImageUrl = strings.ReplaceAll(child.Data.Preview.Images[0].Variants.Gif.Source.Url, "&amp;", "&")
				post.ContentType = "gif"
			} else {
				post.ImageUrl = strings.ReplaceAll(child.Data.Preview.Images[0].Source.Url, "&amp;", "&")
				post.ContentType = "image"
			}
		} else if child.Data.IsVideo {
			err := setVideo(&post, &child)
			if err != nil {
				fmt.Printf("Unable to save with video, skipping item %v, %v", child.Data.Name, err)
				continue
			}
		} else {
			continue
		}

		post.Id = child.Data.Name
		post.Title = child.Data.Title
		post.Url = fmt.Sprintf("https://reddit.com%s", child.Data.Permalink)
		res = append(res, &post)
	}

	return res, nil
}

func setVideo(post *RedditPost, child *RedditPostResponseChild) error {
	videoResolutions := [5]string{"1080", "720", "480", "360", "240"}

	url := child.Data.Media.RedditVideo.FallbackUrl
	resIndex := strings.Index(url, "DASH_") + 5
	extIndex := strings.Index(url, ".mp4")
	if resIndex == 4 || extIndex == -1 {
		return fmt.Errorf("Could not parse resolution")
	}
	originalRes := url[resIndex:extIndex]
	for i := range videoResolutions {
		oResInt, err := strconv.Atoi(originalRes)
		if err != nil {
			continue
		}
		currResInt, err := strconv.Atoi(videoResolutions[i])
		if err != nil {
			continue
		}
		if oResInt < currResInt {
			continue
		}

		contentLength, err := fetchVideoSize(strings.Replace(url, "DASH_"+originalRes, "DASH_"+videoResolutions[i], 1))
		if err != nil || contentLength < 1 {
			continue
		}

		//20 mb is max allowed TG file size for content in URL
		if contentLength < 20971520 {
			width := child.Data.Media.RedditVideo.Width / (oResInt / currResInt)

			replacementString := "DASH_" + videoResolutions[i]
			post.Video = RedditVideo{
				Height:   currResInt,
				Width:    width,
				Duration: child.Data.Media.RedditVideo.Duration,
				Url:      strings.Replace(url, "DASH_"+originalRes, replacementString, 1),
			}

			audioUrl, err := checkAudioUrl(post.Video.Url, replacementString)
			if err != nil {
				continue
			}

			post.Video.AudioUrl = audioUrl
			post.ContentType = "video"
			return nil
		}
	}
	return fmt.Errorf("No video candidates\n")
}

func checkAudioUrl(videoUrl string, replacementString string) (string, error) {
	audioUrl := strings.Replace(videoUrl, replacementString, "DASH_audio", 1)
	req, err := http.NewRequest(http.MethodHead, audioUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error while creating request: %v %v", req, err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending request: %v %v", req, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", nil
	}
	return audioUrl, nil

}

func fetchVideoSize(url string) (int64, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, fmt.Errorf("error while creating request: %v %v", req, err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error while sending request: %v %v", req, err)
	}

	return res.ContentLength, nil
}
