package backend

import "time"

// CampaignType for campaign's type
type CampaignType int

const (
	// CampaignTypeMerchant for campaign merchant
	CampaignTypeMerchant CampaignType = 1

	// CampaignTypeBank for campaign bank
	CampaignTypeBank CampaignType = 2

	// CampaignTypePrivate for campaign private
	CampaignTypePrivate CampaignType = 3
)

// HelloInput input for hello
type HelloInput struct {
	Type      CampaignType
	CreatedAt time.Time
}
