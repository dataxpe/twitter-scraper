package twitterscraper_test

import (
	"strings"
	"testing"
	"time"

	twitterscraper "github.com/dataxpe/twitter-scraper"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetProfileOK(t *testing.T) {
	loc := time.FixedZone("UTC", 0)
	joined := time.Date(2010, 01, 18, 8, 49, 30, 0, loc)
	sample := twitterscraper.Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/436075027193004032/XlDa2oaz_normal.jpeg",
		Banner:    "https://pbs.twimg.com/profile_banners/106037940/1541084318",
		Biography: "nothing",
		//	Birthday:   "March 21",
		IsPrivate:      false,
		IsVerified:     false,
		Joined:         &joined,
		Location:       "Ukraine",
		Name:           "Nomadic",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/nomadic_ua",
		UserID:         "106037940",
		Username:       "nomadic_ua",
		Website:        "https://nomadic.name",
	}

	/*if e := testScraper.LoadTokenStateFromFile("x_tid_state.json"); e != nil {
		t.Error(e)
		return
	}*/

	testScraper.SetTokenState(twitterscraper.TidState{
		Key: "BXzbpo6bveCCmzg0OZygYVlpHchvM4MNqOrgzqKlNkWgzU3niFjfr6fuXbkgVX74",
		KeyBytes: []byte{
			5, 124, 219, 166, 142, 155, 189, 224, 130, 155, 56, 52,
			57, 156, 160, 97, 89, 105, 29, 200, 111, 51, 131, 13,
			168, 234, 224, 206, 162, 165, 54, 69, 160, 205, 77, 231,
			136, 88, 223, 175, 167, 238, 93, 185, 32, 85, 126, 248,
		},
		AnimationKey:  "bf663100100",
		RandomKeyword: "obfiowerehiring",
		RandomNumber:  3,
	})

	profile, err := testScraper.GetProfile("nomadic_ua")
	if err != nil {
		t.Error(err)
		return
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "TweetsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "MediaCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "NormalFollowersCount"),
	}
	if diff := cmp.Diff(sample, profile, cmpOptions...); diff != "" {
		t.Error("Resulting profile does not match the sample", diff)
	}

	if profile.FollowersCount == 0 {
		t.Error("Expected FollowersCount is greater than zero")
	}
	if profile.FollowingCount == 0 {
		t.Error("Expected FollowingCount is greater than zero")
	}
	if profile.LikesCount == 0 {
		t.Error("Expected LikesCount is greater than zero")
	}
	if profile.TweetsCount == 0 {
		t.Error("Expected TweetsCount is greater than zero")
	}
}

func TestGetProfilePrivate(t *testing.T) {
	joined := time.Date(2016, 8, 23, 16, 9, 17, 0, time.UTC)
	sample := twitterscraper.Profile{
		Avatar:    "https://pbs.twimg.com/profile_images/768120429143744512/D-tFFNO8_normal.jpg",
		Biography: "The best job in the world is hobby which get paid. Sometimes logic technique can be more dangerous then using tools",
		//	Birthday:   "March 21",
		IsPrivate:      true,
		IsVerified:     false,
		Joined:         &joined,
		Location:       "in-front LCD Monitor",
		Name:           "▓▬▓ eidelweiss ▓▬▓",
		PinnedTweetIDs: []string{},
		URL:            "https://twitter.com/dummysystems",
		UserID:         "768117919930691586",
		Username:       "dummysystems",
		Website:        "http://eidelweiss-advisories.blogspot.com",
	}

	// some random private profile (found via google)
	profile, err := testScraper.GetProfile("dummysystems")
	if err != nil {
		t.Error(err)
	}

	cmpOptions := cmp.Options{
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowersCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FollowingCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "FriendsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "LikesCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "ListedCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "TweetsCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "MediaCount"),
		cmpopts.IgnoreFields(twitterscraper.Profile{}, "NormalFollowersCount"),
	}
	if diff := cmp.Diff(sample, profile, cmpOptions...); diff != "" {
		t.Error("Resulting profile does not match the sample", diff)
	}

	if profile.FollowingCount == 0 {
		t.Error("Expected FollowingCount is greater than zero")
	}
	if profile.LikesCount == 0 {
		t.Error("Expected LikesCount is greater than zero")
	}
	if profile.TweetsCount == 0 {
		t.Error("Expected TweetsCount is greater than zero")
	}
}

func TestGetProfileErrorSuspended(t *testing.T) {
	_, err := testScraper.GetProfile("1")
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if !strings.Contains(err.Error(), "suspended") {
			t.Error("Expected error to contain 'suspended', got", err)
		}
	}
}

func TestGetProfileErrorNotFound(t *testing.T) {
	neUser := "sample3123131"
	expectedError := "user not found"
	_, err := testScraper.GetProfile(neUser)
	if err == nil {
		t.Error("Expected Error, got success")
	} else {
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err)
		}
	}
}

func TestGetProfileByID(t *testing.T) {
	profile, err := testScraper.GetProfileByID("768117919930691586")
	if err != nil {
		t.Error(err)
	}

	if profile.Username != "dummysystems" {
		t.Errorf("Expected username 'tomdumont', got '%s'", profile.Username)
	}
}

func TestGetUserIDByScreenName(t *testing.T) {
	userID, err := testScraper.GetUserIDByScreenName("X")
	if err != nil {
		t.Errorf("getUserByScreenName() error = %v", err)
	}
	if userID == "" {
		t.Error("Expected non-empty user ID")
	}
}
