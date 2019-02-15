package dao

import (
	"jose/josegin/jgin"
	"jose/josegin/models"
)

type UsersDbDao struct{}

func (dao *UsersDbDao) GetByPk(pk int64) interface{} {
	var user models.Users
	orm := jgin.OrmEngin()
	geted, _ := orm.Id(pk).Get(&user)
	if geted {
		return user
	} else {
		return nil
	}
}

func (dao *UsersDbDao) UpdateByUser(user models.Users) bool {
	engine := jgin.OrmEngin()
	affected, err := engine.Id(user.Uid).Update(user)
	if err != nil {
		return affected > 0
	}
	return false
}

func (dao *UsersDbDao) GetUserByUsername(key string, tp string) interface{} {
	var users []models.Users = make([]models.Users, 0)
	orm := jgin.OrmEngin()
	if tp == "user" {
		orm.Where("email = ?", key).Or("mobile = ?", key).Find(&users)
	} else {
		orm.Where(tp+" = ?", key).Find(&users)
	}

	if len(users) > 0 {
		return users[0]
	} else {
		return nil
	}
}

//func (dao *UsersDbDao)UpdateByPk(id int64,stat int)(int64,error){
//	var user entity.User
//	user.ID=id
//	user.Stat=stat
//	orm := jgin.OrmEngin()
//	r,e:=orm.ID(id).Cols("stat").Update(&user)
//	return r,e
//}

//===================按条件查询相关..

// func (dao *UsersDbDao) buildCond(arg form.UserArg) *xorm.Session {

// 	orm := jgin.OrmEngin()
// 	t := orm.Where("id>0")
// 	if 0 < len(arg.Kword) {
// 		t = t.And("name like ?", "%"+arg.Kword+"%")
// 	}

// 	if !arg.Datefrom.IsZero() {
// 		t = t.And("create_at >= ?", arg.Datefrom)
// 	}
// 	if !arg.Dateto.IsZero() {
// 		t = t.And("create_at <= ?", arg.Dateto)
// 	}
// 	if 0 < len(arg.Mobile) {
// 		t = t.And("mobile = ?", arg.Mobile)
// 	}
// 	return t
// }

// func (dao *UsersDbDao) Count(arg form.UserArg) (n int64) {
// 	var user models.User
// 	t := dao.buildCond(arg)
// 	n, _ = t.Count(&user)
// 	return
// }

// func (dao *UsersDbDao) Get(arg form.UserArg) []models.User {
// 	var users []models.User = make([]models.User, 0)
// 	t := dao.buildCond(arg)
// 	if len(arg.Asc) > 0 {
// 		t = t.Asc(arg.Asc)
// 	}
// 	if len(arg.Desc) > 0 {
// 		t = t.Desc(arg.Desc)
// 	}
// 	t.Limit(arg.GetPageSize(), arg.GetPageFrom()*arg.GetPageSize()).Find(&users)
// 	return users
// }

// func (dao *UsersDbDao) GetOne(arg form.UserArg) interface{} {
// 	users := dao.Get(arg)
// 	if len(users) > 0 {
// 		return users[0]
// 	} else {
// 		return nil
// 	}
// }
