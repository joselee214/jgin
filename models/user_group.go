package models

import (
	"time"
)

type UserGroup struct {
	CreatedAt time.Time `xorm:"default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	Gid       string    `xorm:"not null pk unique(uid) VARCHAR(200)"`
	Mark      string    `xorm:"TEXT"`
	Uid       int       `xorm:"not null pk unique(uid) INT(11)"`
}
