package model

import "database/sql"

type FundInfoPO struct {
	Id            int64           `db:"ID"`
	FundCode      string          `db:"FUND_CODE"`
	Name          string          `db:"NAME"`
	Jzrq          string          `db:"JZRQ"`
	Dwjz          float64         `db:"DWJZ"`
	Gsz           sql.NullFloat64 `db:"GSZ"`
	Gszzl         sql.NullFloat64 `db:"GSZZL"`
	Gztime        sql.NullString  `db:"GZTIME"`
	CreateAt      string          `db:"CREATED_AT"`
	UpdateAt      string          `db:"UPDATED_AT"`
	RecordVersion int32           `db:"RECORD_VERSION"`
	IsDeleted     int32           `db:"IS_DELETED"`
}

type FundInfo struct {
	FundCode string `json:"fundcode"`
	Name     string `json:"name"`
	Jzrq     string `json:"jzrq"`
	Dwjz     string `json:"dwjz"`
	Gsz      string `json:"gsz"`
	Gszzl    string `json:"gszzl"`
	Gztime   string `json:"gztime"`
}
