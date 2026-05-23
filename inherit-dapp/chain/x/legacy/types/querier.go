package types

// QueryPlanParams defines the params for querying a plan
type QueryPlanParams struct {
	PlanID uint64
}

// QueryPlansByCreatorParams defines the params for querying plans by creator
type QueryPlansByCreatorParams struct {
	Creator string
}

// QueryAssetParams defines the params for querying an asset
type QueryAssetParams struct {
	AssetID uint64
}

// QueryAssetsByPlanParams defines the params for querying assets by plan
type QueryAssetsByPlanParams struct {
	PlanID uint64
}

// PlanResponse is the response for a plan query
type PlanResponse struct {
	Plan    LegacyPlan  `json:"plan"`
	Assets  []Asset     `json:"assets"`
}

// PlansResponse is the response for a plans query
type PlansResponse struct {
	Plans []LegacyPlan `json:"plans"`
}
