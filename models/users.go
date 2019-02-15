package models

type Users struct {
	Uid               int    `xorm:"not null pk autoincr INT(11)"`
	Avatarurl         string `xorm:"not null TEXT"`
	Created           int    `xorm:"not null default 0 comment('创建日期') index INT(10)"`
	Email             string `xorm:"index VARCHAR(255)"`
	EncryptedPassword string `xorm:"VARCHAR(100)"`
	Identity          int    `xorm:"not null default 1 comment('会员类型 1:普通用户 3:厂家官方号! 4:厂家员工  5:经销商主账号') TINYINT(1)"`
	IsRegUser         int    `xorm:"not null default 0 comment('是否是实际注册用户') TINYINT(1)"`
	LastLogin         int    `xorm:"not null default 0 comment('最后登录日期') INT(10)"`
	Mobile            string `xorm:"index VARCHAR(20)"`
	Name              string `xorm:"default '' VARCHAR(100)"`
	NameSeted         int    `xorm:"default 0 SMALLINT(2)"`
	TraceCode         string `xorm:"not null default '' VARCHAR(100)"`
	Updated           int    `xorm:"not null default 0 comment('更新日期') INT(10)"`
	Username          string `xorm:"index VARCHAR(200)"`
	ValidateStatus    int    `xorm:"not null default 0 comment('验证状态:0 所有的未验证,1:email已验证；2:手机号码已验证；为两个都验证了, 负数为 禁用') TINYINT(3)"`
}
