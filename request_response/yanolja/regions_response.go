package yanolja

type SubRegionalInfo RegionalInfo

type RegionalInfo struct {
	RegionID          int               `json:"regionId"`
	RegionCode        string            `json:"regionCode"`
	ParentRegionID    int               `json:"parentRegionId"`
	RegionName        string            `json:"regionName"`
	RegionDescription string            `json:"regionDescription"`
	RegionLevel       int               `json:"regionLevel"`
	IsUsed            bool              `json:"isUsed"`
	SubRegions        []SubRegionalInfo `json:"subRegions"`
}
