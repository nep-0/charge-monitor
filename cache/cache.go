package cache

type OutletInfo struct {
	Power       string `json:"power"`
	UsedMinutes int64  `json:"used_minutes"`
	UpdatedAt   int64  `json:"updated_at"`
}

type Cache interface {
	Get(outletId string) (OutletInfo, bool)
	Set(outletId string, info OutletInfo)
	JSON() []byte
	LoadFromJSON(data []byte) error
}
