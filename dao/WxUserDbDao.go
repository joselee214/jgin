package dao

import (
	"jose/josegin/form"
	"jose/josegin/models"

	"github.com/go-xorm/xorm"
	"jose/josegin/jgin"
)

type WxUserDbDao struct{}

func (dao *WxUserDbDao) GetByPk(pk int64) interface{} {
	var user models.WxUsers
	orm := jgin.OrmEngin()
	geted, _ := orm.Id(pk).Get(&user)
	if geted {
		return user
	} else {
		return nil
	}
}

//===================按条件查询相关..

func (dao *WxUserDbDao) buildCond(arg form.WxUserArg) *xorm.Session {

	orm := jgin.OrmEngin()
	t := orm.Where("id>0")
	//if (0<len(arg.Uid)){
	//	t = t.And("uid like ?","%"+arg.Uid+"%")
	//}
	if 0 < len(arg.Openid) {
		t = t.And("openid = ?", arg.Openid)
	}
	if 0 < len(arg.Appid) {
		t = t.And("appid = ?", arg.Appid)
	}
	return t
}

func (dao *WxUserDbDao) Count(arg form.WxUserArg) (n int64) {
	var user models.WxUsers
	t := dao.buildCond(arg)
	n, _ = t.Count(&user)
	return
}

func (dao *WxUserDbDao) Get(arg form.WxUserArg) []models.WxUsers {
	var users []models.WxUsers = make([]models.WxUsers, 0)
	t := dao.buildCond(arg)
	if len(arg.Asc) > 0 {
		t = t.Asc(arg.Asc)
	}
	if len(arg.Desc) > 0 {
		t = t.Desc(arg.Desc)
	}
	t.Limit(arg.GetPageSize(), arg.GetPageFrom()*arg.GetPageSize()).Find(&users)
	return users
}

func (dao *WxUserDbDao) GetOne(arg form.WxUserArg) interface{} {
	users := dao.Get(arg)
	if len(users) > 0 {
		return users[0]
	} else {
		return nil
	}
}

// 具体业务功能

func (dao *WxUserDbDao) GetByOpenid(openid string) interface{} {
	var users []models.WxUsers = make([]models.WxUsers, 0)
	orm := jgin.OrmEngin()
	orm.Where("openid=?", openid).Find(&users)

	if len(users) > 0 {
		return users[0]
	} else {
		return nil
	}
}

func (dao *WxUserDbDao) GetUidFromUnionid(unionid string) int {
	var users []models.WxUsers = make([]models.WxUsers, 0)
	orm := jgin.OrmEngin()
	orm.Where("unionid=?", unionid).Where("uid>0").Find(&users)

	if len(users) > 0 {
		return users[0].Uid
	} else {
		return 0
	}
}

func (dao *WxUserDbDao) UpdateByWxUser(wxuser models.WxUsers) bool {
	engine := jgin.OrmEngin()
	affected, err := engine.Id(wxuser.Id).Update(wxuser)
	if err != nil {
		return affected > 0
	}
	return false
}

func (dao *WxUserDbDao) InsertByWxUser(wxuser models.WxUsers) bool {
	engine := jgin.OrmEngin()
	ret, err := engine.InsertOne(wxuser)
	if err != nil {
		return ret > 0
	}
	return false
}
