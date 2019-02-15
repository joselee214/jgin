package models

type UserLoginLog struct {
	Created int    `xorm:"not null default 0 INT(10)"`
	Ip      string `xorm:"VARCHAR(50)"`
	Uid     int    `xorm:"not null default 0 index INT(11)"`
	Ulid    int64  `xorm:"not null pk autoincr BIGINT(20)"`
}
