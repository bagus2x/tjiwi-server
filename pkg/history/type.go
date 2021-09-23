package history

type Params struct {
	StorageID int64  `form:"storage_id"`
	Status    string `form:"status"`
	StartDate int64  `form:"start"`
	EndDate   int64  `form:"end"`
	Cursor    int64  `form:"cursor"`
	Direction string `form:"dir"`
	Limit     int64  `form:"limit"`
}

type Cursor struct {
	Next     int64 `json:"next"`
	Previous int64 `json:"previous"`
}

type BasePaper struct {
	ID             int64  `json:"id"`
	Gsm            int64  `json:"gsm"`
	Width          int64  `json:"width"`
	Io             int64  `json:"io"`
	MaterialNumber int64  `json:"materialNumber"`
	Quantity       int64  `json:"quantity"`
	Location       string `json:"location"`
}

type Member struct {
	ID       int64  `json:"id"`
	Photo    string `json:"photo"`
	Username string `json:"username"`
}

type GetHistoryResponse struct {
	ID        int64     `json:"id"`
	StorageID int64     `json:"storageID"`
	BasePaper BasePaper `json:"basePaper"`
	Member    Member    `json:"member"`
	Status    string    `json:"status"`
	Affected  int64     `json:"affected"`
	Location  string    `json:"location"`
	CreatedAt int64     `json:"createdAt"`
}

type GetHistoriesResponse struct {
	Cursor    Cursor                `json:"cursor"`
	Histories []*GetHistoryResponse `json:"histories"`
}
