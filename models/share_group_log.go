package models

import (
	"time"
)

type ShareGroupLog struct {
	CreatedAt time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	Gid       string    `xorm:"VARCHAR(200)"`
	Pid       int       `xorm:"index INT(11)"`
	Sglid     int       `xorm:"not null pk autoincr INT(11)"`
	Uid       int       `xorm:"INT(11)"`
}
