package types

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User represents the gql user type.
type User struct {
	Email       string
	Name        string
	Username    string
	Balance     float64
	FreeGPQuota float64
	Country     string
	Stats       struct {
		Star      int
		Project   int
		Planet    int
		Follower  int // Number of fans following me
		Following int // Number of users I am following
	}
	MembershipState string // NA, ACTIVE, SUSPENDED, TRIAL, EXPIRED, STOPPED
	Membership      MembershipInfo
	Developer       struct {
		Status string // NA, UNDER_REVIEW, TRIAL, ACTIVE, INACTIVE, EXPIRED
	}
	Joined time.Time
}

// NameOrEmail gives the email if username is not set.
func (u *User) NameOrEmail() string {
	ret := u.Username
	if bson.IsObjectIdHex(ret) {
		ret = u.Email
	}
	return ret
}

// UserHeaderString gives a row of string for the table header.
func UserHeaderString() []string {
	return []string{
		"Username/Email",
		"Coins",
		"GP Quota",
		"Membership",
		"Developership",
		"Star",
		"Project",
		"Planet",
		"Fans",
		"Following",
		"Joined",
		"Country",
	}
}

// RowString gives a row of string for the table output.
func (u *User) RowString() []string {
	return []string{
		u.NameOrEmail(),
		fmt.Sprintf("%.2f", u.Balance),
		fmt.Sprintf("%.2f", u.FreeGPQuota),
		u.MembershipState,
		u.Developer.Status,
		strconv.Itoa(u.Stats.Star),
		strconv.Itoa(u.Stats.Project),
		strconv.Itoa(u.Stats.Planet),
		strconv.Itoa(u.Stats.Follower),
		strconv.Itoa(u.Stats.Following),
		u.Joined.Format("2006-01-02 15:04:05"),
		u.Country,
	}
}
