package form


type UserArg struct {
	PageArg
	Mobile string `form:"mobile" json:"mobile"`
}
