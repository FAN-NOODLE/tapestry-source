package validate_params

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateParams(t *testing.T) {

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
	c.Request.Header.Add("Content-Type", gin.MIMEXML) // set fake content-type

	/*
				params := &delivery.ParamsAdd{
					ID:           0,
					EnterpriseID: 0,
					OperatorId:   0,
					OperatorName: "",
					Name:         "",
					Type:         0,
					RouteUrl:     "",
					ParamJson:    "",
					Tag:          "",
				}
				_ = delivery.TransAddParams(params)



			params2 := &delivery.ParamsList{
				EnterpriseID: 2,
				Page:         1,
				Size:         10,
				Search:       delivery.ParamsSearch{},
			}
			_ = delivery.TransListParams(params2)

			assert.Equal(t, ValidateParams(c, params), false)
			assert.Equal(t, ValidateParams(c, params), false)



		params3 := &delivery.ParamsStatus{
			DownID:       12,
			EnterpriseId: 2,
			Status:       2,
			FileName:     "asdas.txt",
			Ext:          "sss",
		}
		_ = delivery.TransStatusParams(params3)

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)

		c2.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
		c2.Request.Header.Add("Content-Type", gin.MIMEXML) // set fake content-type
		//assert.Equal(t, ValidateParams(c2, params2), true)

		w3 := httptest.NewRecorder()
		c3, _ := gin.CreateTestContext(w3)

		c2.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"foo\":\"bar\", \"bar\":\"foo\"}"))
		c2.Request.Header.Add("Content-Type", gin.MIMEXML) // set fake content-type
		assert.Equal(t, ValidateParams(c3, params3), false)

	*/

}
