package models

import (
	"time"
)

type Replys struct {
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Content   string    `xorm:"TEXT"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	Deleted   int       `xorm:"not null default 0 index SMALLINT(2)"`
	Pid       int       `xorm:"not null index INT(11)"`
	Sort      int       `xorm:"not null default 0 index INT(11)"`
	Uid       int       `xorm:"not null INT(11)"`
}
