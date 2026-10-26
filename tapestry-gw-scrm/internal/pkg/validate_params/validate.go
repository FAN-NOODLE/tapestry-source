package validate_params

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	scrm_errs "gw-scrm/internal/pkg/errs"
	response "gw-scrm/internal/pkg/resoponse"
	"reflect"
)

var validate *validator.Validate

// ValidateStruct 结构体参数验证
func ValidateStruct(obj interface{}) (bool, string) {

	validate = validator.New()

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	err := validate.Struct(obj)

	if err == nil {
		return true, ""
	}

	for _, e := range err.(validator.ValidationErrors) {
		tag := e.Tag()
		field := e.Field()

		st, _ := v.Type().FieldByName(field)
		msg := st.Tag.Get(tag)

		return false, msg
	}

	return false, err.Error()
}

// ValidateParams 接收参数并验证
func ValidateParams(c *gin.Context, params interface{}) bool {

	err := c.BindJSON(params)

	if err != nil {
		if err.Error() == "EOF" {
			response.Notice(c, errorx.NewErrorMsg(scrm_errs.DefaultErrorCode, "参数JSON格式且不能为空"))
			return false
		}

		response.Notice(c, errorx.NewErrorMsg(scrm_errs.DefaultErrorCode, err.Error()))
		return false
	}

	ok, msg := ValidateStruct(params)
	if ok {
		return true
	}

	response.Notice(c, errorx.NewErrorMsg(scrm_errs.DefaultErrorCode, msg))
	return false
}

//
//func ValidateVariable() {
//
//	myEmail := "joeybloggs.gmail.com"
//
//	errs := validate.Var(myEmail, "required,email")
//
//	if errs != nil {
//		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
//		return
//	}
//
//	// email ok, move on
//}
