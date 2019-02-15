package form


type WxUserArg struct {
	PagesArg
	Uid string `form:"uid" json:"uid"`
	Appid string `form:"appid" json:"appid"`
	Openid string `form:"openid" json:"openid"`
}
