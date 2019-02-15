package controller

import (
	"encoding/json"
	"fmt"
	"jose/josegin/dao"
	"jose/josegin/libs"
	"jose/josegin/models"
	"jose/josegin/service"
	"net/http"
	"strconv"
	"time"

	"jose/josegin/jgin"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/mikemintang/go-curl"
)

//用户管理控制器
type ApisController struct {
	jgin.Controller
}

//用户服务层
//var userService service.UserService

//路由注册
func (ctrl *ApisController) Router(router *gin.Engine) {

	r := router.Group("/apis/wxsession")
	r.GET("", ctrl.wxsession)
	r.GET("test", ctrl.test)
	r.GET("/updateUserInfo", ctrl.updateUserInfo)

	//r.GET("posts/getMyReadDocs",ctrl.getMyReadDocs)

}

var WxUserDbDao dao.WxUserDbDao
var UserService service.UserService


func (ctrl *ApisController) test(ctx *gin.Context) {

	cfg := jgin.GetCfg()

	jgin.SetSession(ctx, "uidxxx", 1)

	var cdata map[string]interface{}
	cdata = make(map[string]interface{})
	cdata["t"] = time.Now().Local().Format("2006-01-02 15:04:05")
	cdata["tt"] = time.Now().Local().String()
	cdata["ttdd"] = time.Now().UTC().String()
	cdata["ttt"] = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	cdata["session_id"], _ = ctx.Cookie(cfg.Session["name"])
	cdata["sess_uid"] = jgin.GetSession(ctx, "uid")
	cdata["sess_uidxxx"] = jgin.GetSession(ctx, "uidxxx")

	ctx.JSON(http.StatusOK, gin.H{"data": cdata, "msg": "", "result": 1})

}

//授权信息
func (ctrl *ApisController) updateUserInfo(ctx *gin.Context) {
	openid := ctx.Query("openid")

	encryptedData := ctx.Query("encryptedData")
	iv := ctx.Query("iv")

	var cdata map[string]string
	cdata = make(map[string]string)

	// meconfig := ctx.Query("meconfig")
	// realname := ctx.Query("realname")

	ret := gin.H{"status": 0, "msg": "updateUserInfo"}
	cfg := jgin.GetCfg()

	if openid != "" {
		cdata["openid"] = openid
		wxuser := getWxUser(ctx, cdata, true)

		//fmt.Println("===============3=", iv)

		if (wxuser.Id != 0) && (encryptedData != "") && (iv != "") {
			var uinfo map[string]string
			uinfo = make(map[string]string)

			ret["status"] = 1
			ret["reguid"] = wxuser.Uid
			// uinfo["nickName"] = ctx.Query("nickName")
			// uinfo["avatarUrl"] = ctx.Query("avatarUrl")
			// uinfo["gender"] = ctx.Query("gender")
			// uinfo["city"] = ctx.Query("city")
			// uinfo["province"] = ctx.Query("province")
			// uinfo["country"] = ctx.Query("country")

			uinfo["openid"] = openid
			uinfo = expUserInfo(encryptedData, iv, wxuser.SessionKey, uinfo)
			//fmt.Println("==============4==", uinfo)
			wxuserN := getWxUser(ctx, uinfo, true)
			//fmt.Println("==============4==", wxuserN)

			isReg := 0

			//uid
			if wxuserN.Uid == 0 {
				var ruinfo map[string]interface{}
				ruinfo = make(map[string]interface{})

				ruinfo["name"] = wxuserN.Nickname
				ruinfo["username"] = "wx_" + wxuserN.Openid
				ruinfo["avatarUrl"] = wxuserN.Avatarurl
				ruinfo["pw"] = fmt.Sprintf("%d", time.Now().Unix())
				ruinfo["validate_status"] = "0"
				user, err := UserService.Register(ruinfo)

				// fmt.Println("==============555=", user, err)

				if err == nil {
					isReg = 1
					ret["reguid"] = user.Uid
					wxuserN.Uid = user.Uid
					//更新 wxuser uid//
					wxuserN.Updated = time.Now().Local()
					WxUserDbDao.UpdateByWxUser(wxuserN)
				}
			}

			if  user,ok := UserService.GetUserById(int64(wxuser.Uid)).(models.Users);ok {
				//session //
				UserService.LoginUid(ctx, wxuser.Uid)

				//处理RSession
				sessionid:=ctx.Query("__session_id")
				if sessionid!="" {
					UserService.RSessionSet(sessionid,"uid",wxuser.Uid)
				}

				//updateUinfo
				if (isReg == 0) && (user.Avatarurl != wxuserN.Avatarurl) {
					user.Avatarurl = wxuserN.Avatarurl
					UserService.UpdateByUser(user)
				}
				//set updateUinfo cache
				redisconn := UserService.GetRedis("default")
				rrkey := openid + "_" + cfg.App["xcxappid"]
				redis.String(redisconn.Do("SET", rrkey, 1, 3*86400))
				//fmt.Println("========>===", check)

			}

			ret["isReg"] = isReg
		}

		//realname
		meconfig := ctx.Query("meconfig")
		realname := ctx.Query("realname")
		realname = libs.StripTags(realname)

		if (meconfig != "") && (realname != "") {
			if user_up,ok := UserService.GetUserById(int64(wxuser.Uid)).(models.Users);ok {
				user_up.Name = realname
				user_up.NameSeted = 1
				UserService.UpdateByUser(user_up)
				ret["status"] = 1
				ret["realname"] = realname
			}
		}

	}
	////////////////////
	ctx.JSON(http.StatusOK, ret)
}

/*
	解密获取 unionId
*/
func expUserInfo(encryptedData string, iv string, sessionKey string, uinfo map[string]string) map[string]string {
	// ret["zzz"] = "zzz"
	plainText := libs.WxDecryptData(encryptedData, iv, sessionKey)

	var expinfo map[string]interface{}
	err := json.Unmarshal(plainText, &expinfo)

	if err == nil {

		if _, ok := expinfo["nickName"].(string); ok {
			uinfo["nickName"] = expinfo["nickName"].(string)
			uinfo["nickName"] = libs.StripTags(uinfo["nickName"])
		}
		if _, ok := expinfo["avatarUrl"].(string); ok {
			uinfo["avatarUrl"] = expinfo["avatarUrl"].(string)
			uinfo["avatarUrl"] = libs.StripTags(uinfo["avatarUrl"])
		}
		if _, ok := expinfo["gender"].(float64); ok {
			uinfo["gender"] = strconv.FormatFloat(expinfo["gender"].(float64), 'f', -1, 64)
		}
		if _, ok := expinfo["city"].(string); ok {
			uinfo["city"] = expinfo["city"].(string)
			uinfo["city"] = libs.StripTags(uinfo["city"])
		}
		if _, ok := expinfo["province"].(string); ok {
			uinfo["province"] = expinfo["province"].(string)
			uinfo["province"] = libs.StripTags(uinfo["province"])
		}
		if _, ok := expinfo["country"].(string); ok {
			uinfo["country"] = expinfo["country"].(string)
			uinfo["country"] = libs.StripTags(uinfo["country"])
		}
		if _, ok := expinfo["unionId"].(string); ok {
			uinfo["unionId"] = expinfo["unionId"].(string)
		}

	}
	return uinfo
}

//入口
func (ctrl *ApisController) wxsession(ctx *gin.Context) {
	_jscode := ctx.Query("_jscode")
	cfg := jgin.GetCfg()

	if _jscode != "" {
		url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + cfg.App["xcxappid"] + "&secret=" + cfg.App["xcxsecret"] + "&js_code=" + _jscode + "&grant_type=authorization_code"

		req := curl.NewRequest()
		resp, err := req.SetUrl(url).Get()

		if err != nil {
			fmt.Println(err)
			jgin.ResultFail(ctx, "error request interface")
		} else {
			if resp.IsOk() {
				//fmt.Println("resp.Body",resp.Body)

				var wxloginresp map[string]interface{}
				err1 := json.Unmarshal([]byte(resp.Body), &wxloginresp)

				if err1 == nil {
					if session_key, ok := wxloginresp["session_key"].(string); ok {
						if openid, ok2 := wxloginresp["openid"].(string); ok2 {

							fmt.Println("==============根据openid查询：oSjTH5UdrWdp_MpiffgRdy9K1wes", wxloginresp)

							var cdata map[string]string
							cdata = make(map[string]string)
							cdata["openid"] = openid           //"oSjTH5UdrWdp_MpiffgRdy9K1wes"  //wxloginresp["openid"].(string)
							cdata["session_key"] = session_key //"Q+Yj/6zDxt/wWtwhBR3hpw==" //wxloginresp["session_key"]
							wxuser := getWxUser(ctx, cdata, false)

							fmt.Println("=========== getUser =====", wxuser)

							wxloginresp["unionid"] = wxuser.Unionid
							wxloginresp["gcode"] = wxuser.Gcode

							site_domain := cfg.App["protocal"] + "://" + cfg.App["domain"]
							wxloginresp["site_domain"] = site_domain
							wxloginresp["can_sign_up"] = 1
							wxloginresp["can_session_login"] = 0

							sessionkey := wxuser.Gcode + cdata["openid"] + string(time.Now().Unix())
							wxloginresp["sessionid"] = UserService.GetSessionId( sessionkey )

							wxloginresp["show_web_version"] = "也可用电脑浏览器打开WEB版" + site_domain + "进行文件操作！"

							redisconn := UserService.GetRedis("default")
							rrkey := cdata["openid"] + "_" + cfg.App["xcxappid"]
							check, _ := redis.String(redisconn.Do("GET", rrkey))

							fmt.Println("===============99999===", check, wxuser)

							if (wxuser.Nickname == "") || (check == "") {
								wxloginresp["updateUinfo"] = 4000
							}

							wxloginresp["expires_in"] = wxuser.ExpiresIn
							wxloginresp["introrealname"] = "由于无法获取群用户备注名,为方便阅读统计,故需要确认您的真实姓名"

							if wxuser.Uid != 0 {
								ret := UserService.GetUserById(int64(wxuser.Uid))
								fmt.Println("===============9999911===", ret)
								if user,ok := ret.(models.Users);ok{
									UserService.LoginUid(ctx, wxuser.Uid)
									sessionid := UserService.GetSessionId(sessionkey)
									UserService.RSessionSet(sessionid,"uid",wxuser.Uid)

									wxloginresp["uid"] = user.Uid
									wxloginresp["realname"] = user.Name
									wxloginresp["realname_seted"] = user.NameSeted

									if user.Name == "" {
										wxloginresp["introrealname"] = ""
									}
								}
							}

							wxloginresp["intro"] = "「通知说」免费文档工具，分享至微信群，可以统计文档的已读/未读人数；可以分享各种文档，在线报名/结果下载。电脑打开web版查看 " + site_domain

							ctx.JSON(http.StatusOK, gin.H{"data": wxloginresp, "msg": "", "result": 1})

						} else {
							ctx.JSON(http.StatusOK, gin.H{"data": resp.Body, "jscode": _jscode, "result": 0, "f": 4})
						}
					} else {
						// jgin.ResultFail(ctx, resp.Body)
						ctx.JSON(http.StatusOK, gin.H{"data": resp.Body, "jscode": _jscode, "result": 0, "f": 4})
					}
				} else {
					jgin.ResultFail(ctx, resp.Body)
				}
			} else {
				jgin.ResultFail(ctx, "error request interface")
				ctx.JSON(http.StatusOK, gin.H{"msg": "error request interface", "result": 0, "f": 5})
			}
		}

	} else {
		// jgin.ResultFail(ctx, "empty jscode")
		ctx.JSON(http.StatusOK, gin.H{"msg": "empty jscode", "result": 0, "f": 5})
	}

	//ret := userService.FindOne(2)
	//res := userService.FindOne(userId)

	//rett := make(map[string]interface{})
	//rett["ret"] = ret
	//rett["res"] = res

	//rett := output{ret:ret,res:res,zzz:11}
}

//查找/更新wxuser
func getWxUser(ctx *gin.Context, udata map[string]string, onlycheck bool) models.WxUsers {

	var wxuser models.WxUsers

	ret := WxUserDbDao.GetByOpenid(udata["openid"])

	if onlycheck {
		if ret == nil {
			return wxuser
		}
	}

	gcode := UserService.InitGcode(ctx)
	//fmt.Println("wxuser :: wxuser :: ", gcode)

	update := false
	// insert := false
	uid := 0

	if _, ok := udata["unionId"]; ok {
		uid = WxUserDbDao.GetUidFromUnionid(udata["unionId"])
	}

	if xwxuser,ok := ret.(models.WxUsers);ok {
		wxuser = xwxuser
		//修改更新用户信息
		if _, ok := udata["avatarUrl"]; ok {
			if wxuser.Avatarurl != udata["avatarUrl"] {
				wxuser.Avatarurl = libs.StripTags(udata["avatarUrl"])
				update = true
			}
		}
		if _, ok := udata["nickName"]; ok {
			if wxuser.Nickname != udata["nickName"] {
				wxuser.Nickname = libs.StripTags(udata["nickName"])
				update = true
			}
		}
		if _, ok := udata["gender"]; ok {
			if wxuser.Gender != udata["gender"] {
				wxuser.Gender = libs.StripTags(udata["gender"])
				update = true
			}
		}
		if _, ok := udata["city"]; ok {
			if wxuser.City != udata["city"] {
				wxuser.City = libs.StripTags(udata["city"])
				update = true
			}
		}
		if _, ok := udata["province"]; ok {
			if wxuser.Province != udata["province"] {
				wxuser.Province = libs.StripTags(udata["province"])
				update = true
			}
		}
		if _, ok := udata["country"]; ok {
			if wxuser.Country != udata["country"] {
				wxuser.Country = libs.StripTags(udata["country"])
				update = true
			}
		}
		if _, ok := udata["unionId"]; ok {
			if wxuser.Unionid != udata["unionId"] {
				wxuser.Unionid = udata["unionId"]
				update = true
			}
		}

	} else {
		// insert = true
		//新增wxuser
		cfg := jgin.GetCfg()
		wxuser.ExpiresIn = 300
		wxuser.Gcode = gcode
		wxuser.Appid = cfg.App["xcxappid"]
		wxuser.Updated = time.Now().Local()
		wxuser.Created = time.Now().Local()
		wxuser.Openid = udata["openid"]
		if _, ok := udata["unionId"]; ok {
			wxuser.Unionid = udata["unionId"]
		}
		if _, ok := udata["session_key"]; ok {
			wxuser.SessionKey = udata["session_key"]
		}

		WxUserDbDao.InsertByWxUser(wxuser)
		retinsert := WxUserDbDao.GetByOpenid(udata["openid"])
		if zwxuser,ok := retinsert.(models.WxUsers);ok{
			wxuser = zwxuser
		}
	}

	if (wxuser.Uid != uid) && (uid > 0) {
		wxuser.Uid = uid
		update = true
	}
	if gcode != wxuser.Gcode {
		update = true
		wxuser.Gcode = gcode
	}
	if _, ok := udata["session_key"]; ok {
		if wxuser.SessionKey != udata["session_key"] {
			wxuser.SessionKey = udata["session_key"]
			update = true
		}
	}
	if update {
		wxuser.Updated = time.Now().Local()
		WxUserDbDao.UpdateByWxUser(wxuser)
	}
	return wxuser
}
