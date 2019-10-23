package types

import (
	"fmt"
	"time"
)

// MembershipInfo represents the gql MembershipInfo type.
type MembershipInfo struct {
	State             string               `json:"state"`
	PlanName          string               `json:"planName"`
	Period            int                  `json:"period"`
	StartDate         time.Time            `json:"startDate"`
	EndDate           time.Time            `json:"endDate"`
	MemberGPQuota     float64              `json:"memberGPQuota"`
	CoinPerGP         float64              `json:"coinPerGP"`
	AssetStorage      float64              `json:"assetStorage"`
	Visibility        []string             `json:"visibility"`
	Coupon            MembershipPlanCoupon `json:"coupon"`
	ModelPerProject   int                  `json:"modelPerProject"`
	CollaboratorQuota int                  `json:"collaboratorQuota"`
	ForceWatermark    bool                 `json:"forceWaterMark"`
}

// MembershipPlanCoupon represents the gql MEMBERSHIP_PLAN_COUPON type.
type MembershipPlanCoupon struct {
	Value      int `json:"value"`
	Repeat     int `json:"repeat"`
	ValidMonth int `json:"validMonth"`
}

func (c MembershipPlanCoupon) String() string {
	return fmt.Sprintf("Value: %d\nRepeat: %d\nMonth: %d", c.Value, c.Repeat, c.ValidMonth)
}
