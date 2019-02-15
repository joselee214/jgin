package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"jose/josegin/jgin"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"

	"jose/josegin/dao"
	"jose/josegin/models"
)

type UserService struct{}

var usersDbDao dao.UsersDbDao

//根据userId 获取用户编号
// func (service *UserService) FindOne(userId int64) interface{} {
// 	return userDbDao.GetByPk(userId)
// }

// func (service *UserService) Count(arg form.UserArg) (n int64) {
// 	n = userDbDao.Count(arg)
// 	return
// }

// func (service *UserService) Query(arg form.UserArg) []models.User {
// 	users := userDbDao.Get(arg)
// 	return users
// }

// func (service *UserService) GetOne(arg form.UserArg) interface{} {
// 	user := userDbDao.GetOne(arg)
// 	return user
// }

// func (service *UserService) UpdateStat(id int64, stat int) (int64, error) {
// 	var user models.User
// 	user.ID = id
// 	user.Stat = stat
// 	orm := jgin.OrmEngin()
// 	r, e := orm.ID(id).Cols("stat").Update(&user)
// 	return r, e
// }

// //登录服务,通过手机号/邮箱/用户名登录
// func (service *UserService) Login(ctx *gin.Context, kword string, passwd string) (u models.User, err error) {
// 	ismobile := jgin.IsMobile(kword)
// 	isemail := jgin.IsEmail(kword)
// 	var user models.User
// 	orm := jgin.OrmEngin()

// 	if ismobile {
// 		_, err = orm.Where("mobile = ?", kword).Get(&user)
// 	} else if isemail {
// 		_, err = orm.Where("email = ?", kword).Get(&user)
// 	} else {
// 		_, err = orm.Where("account = ?", kword).Get(&user)
// 	}
// 	if err != nil {
// 		return
// 	}
// 	if user.ID == 0 {
// 		err = errors.New("该用户不存在")
// 		return
// 	}
// 	if jgin.Md5encode(passwd) != user.Passwd {
// 		err = errors.New("密码不正确,请重试")
// 		return
// 	}
// 	u = user
// 	jgin.SetSession(ctx, "roleid", u.RoleId)
// 	return
// }

func (service *UserService) UpdateByUser(user models.Users) bool {
	return usersDbDao.UpdateByUser(user)
}

func (service *UserService) LoginUid(ctx *gin.Context, uid int) {
	var loginlog models.UserLoginLog

	loginlog.Uid = uid
	loginlog.Created = int(time.Now().Unix())
	loginlog.Ip = ctx.ClientIP()

	orm := jgin.OrmEngin()
	ret, err := orm.InsertOne(loginlog)

	fmt.Println("===============inser login log OK ===", ret, err)

	if user,ok := usersDbDao.GetByPk(int64(uid)).(models.Users);ok {
		user.LastLogin = int(time.Now().Unix())
		service.UpdateByUser(user)
	}

	// jgin.SetSession(ctx, "uid", uid)
	service.SetSession(ctx, "uid", uid)
	return
}

func (service *UserService) SetSession(ctx *gin.Context, key string, value interface{}) {
	jgin.SetSession(ctx, key, value)
}



func (service *UserService) GetUserById(uid int64) interface{} {
	user := usersDbDao.GetByPk(uid)
	return user
}

func (service *UserService) GetRedis(rtype string) redis.Conn {
	//配置读取
	cfg := new(jgin.Config)
	cfg.Parse("config/app.properties")
	jgin.SetCfg(cfg)

	//连接 Redis
	for k, ds := range cfg.Redis {
		db, _ := strconv.Atoi(ds["db"])
		maxIdle, _ := strconv.Atoi(ds["maxIdle"])
		maxActive, _ := strconv.Atoi(ds["maxActive"])
		idleTimeoutSecond, _ := strconv.Atoi(ds["idleTimeoutSecond"])
		jgin.ConnPools.AddRedis(k, ds["host"], ds["passwd"], db, maxIdle, maxActive, idleTimeoutSecond)
	}
	redisconn := jgin.ConnPools.GetRedisA(rtype)
	//defer redisconn.Close()
	return redisconn
}

//注册服务,注册后自动登录
func (service *UserService) Register(uinfo map[string]interface{}) (p models.Users, err error) {
	var user models.Users

	if _, ok := uinfo["email"].(string); ok {
		user.Email = uinfo["email"].(string)
	}
	if _, ok := uinfo["mobile"].(string); ok {
		user.Mobile = uinfo["mobile"].(string)
	}
	if _, ok := uinfo["username"].(string); ok {
		user.Username = uinfo["username"].(string)
	}
	if _, ok := uinfo["name"].(string); ok {
		user.Name = uinfo["name"].(string)
	}
	if _, ok := uinfo["validate_status"].(string); ok {
		user.ValidateStatus, _ = strconv.Atoi(uinfo["validate_status"].(string))
	}
	if _, ok := uinfo["trace_code"].(string); ok {
		user.TraceCode = uinfo["trace_code"].(string)
	}
	if _, ok := uinfo["pw"].(string); ok {
		user.EncryptedPassword = jgin.Md5encode(uinfo["pw"].(string))
	} else {
		pwd := fmt.Sprintf("%d", time.Now().Unix())
		user.EncryptedPassword = jgin.Md5encode(pwd)
	}

	cfg := jgin.GetCfg()
	site_domain := cfg.App["protocal"] + "://" + cfg.App["domain"]
	user.Avatarurl = site_domain + "/res/images/tzs108.png"
	if _, ok := uinfo["avatarUrl"].(string); ok {
		user.Avatarurl = uinfo["avatarUrl"].(string)
	}

	user.Identity = 1

	c := usersDbDao.GetUserByUsername(user.Username, "username")
	if c == nil {
		engine := jgin.OrmEngin()
		user.Created = int(time.Now().Unix())
		user.LastLogin = int(time.Now().Unix())
		_, errcc := engine.InsertOne(user)

		if errcc == nil {
			if tempp,ok := usersDbDao.GetUserByUsername(user.Username, "username").(models.Users);ok {
				p = tempp
			}
			return
		}
	}

	err = errors.New("注册用户失败")
	return
}

// func (service *UserService) xxxRegister(ctx *gin.Context, user *models.User) (p *models.User, err error) {

// 	isemail := jgin.IsEmail(user.Email)
// 	if !isemail {
// 		err = errors.New("email格式不正确")
// 		return
// 	}
// 	if len(user.Passwd) < 6 {
// 		err = errors.New("注册失败,太短了")
// 		return
// 	}
// 	var u models.User

// 	orm := jgin.OrmEngin()
// 	t := orm.Where("id>0")

// 	t.Where("email=?", user.Email)

// 	t.Get(&u)
// 	if u.ID > 0 {
// 		err = errors.New("该账户已存在")
// 		return
// 	}
// 	user.Passwd = jgin.Md5encode(user.Passwd)
// 	user.Stat = 1
// 	//user.CreateAt = jgin.JsonDateTime(time.Now())
// 	user.ID, err = orm.InsertOne(user)
// 	jgin.SetSession(ctx, "user", user)
// 	p = user
// 	return
// }

//===================== 业务相关

func (service *UserService) InitGcode(ctx *gin.Context) string {
	gcode := ctx.Query("__gcode")
	if gcode == "" {
		gcodecookie, err := ctx.Cookie("gcode")
		if err != nil {
			//无cookie
			gcode = fmt.Sprintf("%d%04v_%04v", time.Now().Unix(), rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000), rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
			ctx.SetCookie("gcode", gcode, 86400000, "/", "upimgs.cn", false, true)
		} else {
			gcode = gcodecookie
		}
	}
	return gcode
}



//============================== store session in redis
func (service *UserService) GetSessionId(userkey string) string {
	return jgin.Md5encode( userkey+"7yes" )
}

func (service *UserService) RSessionSet(sessionid string,key string,value interface{}) bool {
	sessionmap,ok := service.GetRSessionBySessionId(sessionid).(map[string]interface{})
	if !ok {
		sessionmap = make(map[string]interface{})
	}
	sessionmap[key] = value
	return service.SetRSessionBySessionId(sessionid,sessionmap)
}
func (service *UserService) RSessionGet(sessionid string,key string) interface{} {
	if sessionmap,ok := service.GetRSessionBySessionId(sessionid).(map[string]interface{});ok {
		if _, ok := sessionmap[key]; ok {
			return sessionmap[key]
		}
	}
	return nil
}
func (service *UserService) RSessionDel(sessionid string,key string) bool {
	if sessionmap,ok := service.GetRSessionBySessionId(sessionid).(map[string]interface{});ok {
		if _, ok := sessionmap[key]; ok {
			delete(sessionmap,key)
			return service.SetRSessionBySessionId(sessionid,sessionmap)
		}
	}
	return false
}
func (service *UserService) RSessionClear(sessionid string) bool {
	redisconn := service.GetRedis("default")
	_, err := redisconn.Do("DEL", sessionid)
	return nil==err
}

//通过sessionid获取session的map
func (service *UserService) GetRSessionBySessionId(sessionid string) interface{} {
	redisconn := service.GetRedis("default")
	storeStr, err := redis.String(redisconn.Do("GET", sessionid))
	if err==nil {
		var sessionmap map[string]interface{}
		err1 := json.Unmarshal([]byte(storeStr), &sessionmap)
		if err1==nil {
			return sessionmap
		}
	}
	return nil
}
func (service *UserService) SetRSessionBySessionId(sessionid string,smap map[string]interface{}) bool {
	redisconn := service.GetRedis("default")
	storeStr,_ := json.Marshal(smap)
	cfg := new(jgin.Config)
	cfg.Parse("config/app.properties")
	sessionmaxLifetime, _ := strconv.Atoi(cfg.Session["timelive"])
	_, err := redisconn.Do("SETEX", sessionid, sessionmaxLifetime, string(storeStr))
	//fmt.Println("zzzzz=>>>>",sessionid, sessionmaxLifetime,string(storeStr),cfg.Session["timelive"],err)
	return nil==err
}
//==============================End  store session in redis
