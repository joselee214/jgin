package models

import (
	"time"
)

type Posts struct {
	Id            int       `xorm:"not null pk autoincr INT(11)"`
	Allowreply    int       `xorm:"default 0 SMALLINT(2)"`
	CanSignUp     int       `xorm:"not null default 0 SMALLINT(2)"`
	Content       string    `xorm:"TEXT"`
	CreatedAt     time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	Deleted       int       `xorm:"not null default 0 index SMALLINT(1)"`
	Files         string    `xorm:"TEXT"`
	Gids          string    `xorm:"TEXT"`
	Imgs          string    `xorm:"TEXT"`
	Replytimes    int       `xorm:"not null default 0 INT(11)"`
	SignUpOptions string    `xorm:"TEXT"`
	Signtimes     int       `xorm:"not null default 0 INT(11)"`
	Title         string    `xorm:"TEXT"`
	Uid           int       `xorm:"index INT(11)"`
	Viewtimes     int       `xorm:"not null default 0 INT(11)"`
}
