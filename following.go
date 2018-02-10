// Package igfollow helps you get Instagram following and followers users.
package igfollow

// Get all following users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const UrlFollowers = `https://i.instagram.com/api/v1/friendships/{{USERID}}/followers/`
const UrlFollowing = `https://i.instagram.com/api/v1/friendships/{{USERID}}/following/`

type RawFollow struct {
	Users     []FollowUser `json:"users"`
	BigList   bool         `json:"big_list"` // if false, no next_max_id in response
	PageSize  int64        `json:"page_size"`
	Status    string       `json:"status"`
	NextMaxId string       `json:"next_max_id"` // used for pagination if list is too big
}

type FollowUser struct {
	Pk       int64  `json:"pk"`
	Username string `json:"username"`
}

// id: Id of the Instagram user whose followers to be retrieved.
//
// ds_user_id, sessionid, csrftoken: Login Instagram and get these three cookie
// values from Chrome Developer Tools. See https://stackoverflow.com/a/44773079
func GetFollowers(id, ds_user_id, sessionid, csrftoken string) (rf RawFollow, err error) {
	url := strings.Replace(UrlFollowers, "{{USERID}}", id, 1)
	rf, err = GetFollowResponse(url, ds_user_id, sessionid, csrftoken)
	if err != nil {
		return
	}

	// If the list is too big and next_max_id is not ""
	for rf.NextMaxId != "" {
		urln := url + "?max_id=" + rf.NextMaxId
		rfn, err := GetFollowResponse(urln, ds_user_id, sessionid, csrftoken)
		if err != nil {
			return rf, err
		}
		rf.Users = append(rf.Users, rfn.Users...)
		rf.NextMaxId = rfn.NextMaxId
	}
	return
}

// id: Id of the Instagram user whose following to be retrieved.
//
// ds_user_id, sessionid, csrftoken: Login Instagram and get these three cookie
// values from Chrome Developer Tools. See https://stackoverflow.com/a/44773079
func GetFollowing(id, ds_user_id, sessionid, csrftoken string) (rf RawFollow, err error) {
	url := strings.Replace(UrlFollowing, "{{USERID}}", id, 1)
	rf, err = GetFollowResponse(url, ds_user_id, sessionid, csrftoken)
	if err != nil {
		return
	}

	// If the list is too big and next_max_id is not ""
	for rf.NextMaxId != "" {
		urln := url + "?max_id=" + rf.NextMaxId
		rfn, err := GetFollowResponse(urln, ds_user_id, sessionid, csrftoken)
		if err != nil {
			return rf, err
		}
		rf.Users = append(rf.Users, rfn.Users...)
		rf.NextMaxId = rfn.NextMaxId
	}
	return
}

func GetFollowResponse(url, ds_user_id, sessionid, csrftoken string) (rf RawFollow, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.AddCookie(&http.Cookie{Name: "ds_user_id", Value: ds_user_id})
	req.AddCookie(&http.Cookie{Name: "sessionid", Value: sessionid})
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: csrftoken})

	req.Header.Set("User-Agent", "Instagram 10.3.2 (iPhone7,2; iPhone OS 9_3_3; en_US; en-US; scale=2.00; 750x1334) AppleWebKit/420+")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("resp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&rf)
	return
}
