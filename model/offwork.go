package model

type CheckinRequest struct {
	Lng float64 `json:"lng" binding:"required"`
	Lat float64 `json:"lat" binding:"required"`
}

type NearbyGridItem struct {
	Lng   float64 `json:"lng"`
	Lat   float64 `json:"lat"`
	Count int     `json:"count"`
}

type NearbyResponse struct {
	AllCount    int              `json:"all_count"`
	OffworkData []NearbyGridItem `json:"offwork_data"`
}
