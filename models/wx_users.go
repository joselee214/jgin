package models

import (
	"time"
)

type WxUsers struct {
	Id         int       `xorm:"not null pk autoincr INT(12)"`
	Appid      string    `xorm:"unique(appidopenid) VARCHAR(65)"`
	Avatarurl  string    `xorm:"TEXT"`
	City       string    `xorm:"VARCHAR(65)"`
	Country    string    `xorm:"VARCHAR(65)"`
	Created    time.Time `xorm:"TIMESTAMP"`
	ExpiresIn  int       `xorm:"not null INT(6)"`
	Gcode      string    `xorm:"VARCHAR(63)"`
	Gender     string    `xorm:"VARCHAR(8)"`
	Nickname   string    `xorm:"VARCHAR(130)"`
	Openid     string    `xorm:"not null unique(appidopenid) VARCHAR(65)"`
	Province   string    `xorm:"VARCHAR(65)"`
	SessionId  string    `xorm:"default '' comment('sessionid,用于session失效后重复使用') VARCHAR(64)"`
	SessionKey string    `xorm:"VARCHAR(65)"`
	Tcode      string    `xorm:"comment('跟踪') VARCHAR(63)"`
	Tfid       int       `xorm:"comment('跟踪') INT(11)"`
	Tsid       int       `xorm:"comment('跟踪') INT(11)"`
	Tuid       int       `xorm:"comment('跟踪') INT(11)"`
	Uid        int       `xorm:"default 0 index INT(11)"`
	Unionid    string    `xorm:"index VARCHAR(65)"`
	Updated    time.Time `xorm:"default 'CURRENT_TIMESTAMP' comment('更新日期') TIMESTAMP"`
}
