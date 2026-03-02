package model

type FundConfigPO struct {
	Id            int64  `db:"ID"`
	FundCode      string `db:"FUND_CODE"`
	FundType      int32  `db:"FUND_TYPE"`
	FundName      string `db:"FUND_NAME"`
	Sort          int32  `db:"SORT"`
	CreateAt      string `db:"CREATED_AT"`
	UpdateAt      string `db:"UPDATED_AT"`
	RecordVersion int32  `db:"RECORD_VERSION"`
	IsDeleted     int32  `db:"IS_DELETED"`
}
