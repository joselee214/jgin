package models

import (
	"time"
)

type Viewlog struct {
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	Gid       string    `xorm:"default '' index VARCHAR(200)"`
	Pid       int       `xorm:"not null pk unique(key) index INT(11)"`
	Uid       int       `xorm:"not null pk unique(key) index INT(11)"`
	UpdatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' index TIMESTAMP"`
}
