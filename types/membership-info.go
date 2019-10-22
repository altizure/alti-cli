package types

import "time"

// MembershipInfo represents the gql MembershipInfo type.
type MembershipInfo struct {
	State             string
	PlanName          string
	Period            int
	StartDate         time.Time
	EndDate           time.Time
	MemberGPQuota     float64
	CoinPerGP         float64
	AssetStorage      float64
	Visibility        []string
	Coupon            MembershipPlanCoupon
	ModelPerProject   int
	CollaboratorQuota int
	ForceWatermark    bool
}

// MembershipPlanCoupon represents the gql MEMBERSHIP_PLAN_COUPON type.
type MembershipPlanCoupon struct {
	Value      int
	Repeat     int
	ValidMonth int
}
