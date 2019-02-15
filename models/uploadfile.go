package models

import (
	"time"
)

type Uploadfile struct {
	CreatedAt time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	Deleted   int       `xorm:"not null default 0 index SMALLINT(2)"`
	Filename  string    `xorm:"VARCHAR(200)"`
	Filesize  int       `xorm:"INT(10)"`
	Filetype  string    `xorm:"default '' VARCHAR(100)"`
	Fromtype  int       `xorm:"not null default 0 comment('0附件1图片') SMALLINT(2)"`
	Gcode     string    `xorm:"VARCHAR(20)"`
	Id        int       `xorm:"not null pk autoincr INT(11)"`
	Path      string    `xorm:"VARCHAR(200)"`
	Pid       int       `xorm:"not null default 0 index INT(11)"`
	Uid       int       `xorm:"default 0 index INT(11)"`
}
