package models

import (
	"time"
)

type SignUpLog struct {
	CreatedAt time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	Pid       int       `xorm:"not null pk index unique(uid_2) INT(11)"`
	Signup    string    `xorm:"TEXT"`
	Uid       int       `xorm:"not null pk index unique(uid_2) INT(11)"`
}
