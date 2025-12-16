package models

type Campaign struct {
	CampaignName string  `json:"campaignName"`
	Clicks       int     `json:"clicks"`
	Conversions  float64 `json:"conversions"`
	Cost         float64 `json:"cost"`
	Impressions  int     `json:"impressions"`
}
type GoogleAdwAccount struct {
	AccountId              string     `json:"account_id"`
	AccountName            string     `json:"account_name"`
	AccountCost            float64    `json:"account_cost"`
	AccountOptimazionScore float64    `json:"account_optimazion_score"`
	Campaigns              []Campaign `json:"account_campaigns"`
}
