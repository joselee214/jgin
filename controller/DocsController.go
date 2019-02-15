package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jose/josegin/jgin"
	"jose/josegin/libs"
	_"jose/josegin/libs"
	"jose/josegin/models"
	"net/http"
	"strconv"
)

type DocsController struct {
	jgin.Controller
}


func (ctrl *DocsController)Router(router *gin.Engine){
	r := router.Group("/apis/posts")
	r.GET("/getMyReadDocs",ctrl.getMyReadDocs)
}

func (ctrl *DocsController)getMyReadDocs(ctx *gin.Context){
	sessionkey:=ctx.Query("__session_id")
	if sessionkey!=""{
		if sessUid,ok := UserService.RSessionGet(sessionkey,"uid").(float64);ok {
			uid,_ := strconv.ParseInt(strconv.FormatFloat(sessUid, 'f', -1, 64),10,64)
			if user,ok := UserService.GetUserById(uid).(models.Users);ok {
				fmt.Println("=========>",user,libs.Typeof(user))
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": 1, "msg": "", "result": 1})
}


//func (ctrl *OpenController)getMyReadDocs(ctx *gin.Context){
//	d := make([]byte, 4)
//	s := gin.NewLen(4)
//	ss := ""
//	d = []byte(s)
//	for v := range d {
//		d[v] %= 10
//		ss += strconv.FormatInt(int64(d[v]), 32)
//	}
//	session := sessions.Get(ctx)
//	session.Set("___verify",ss)
//	session.Save()
//	gin.NewImage(d, 100, 40).WriteTo(ctx.Writer)
//}
