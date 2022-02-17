package chuanglan

import (
	"bytes"
	"context"
	xerr "github.com/goclub/error"
	xhttp "github.com/goclub/http"
	xjson "github.com/goclub/json"
	xtime "github.com/goclub/time"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type V1 struct {
	client *xhttp.Client
	auth AuthV1
}
type AuthV1 struct {
	Origin string
	Account string
	Password string
}
func NewV1(httpClient *http.Client, auth AuthV1) (V1, error) {
	if auth.Origin == "" {
		return V1{},xerr.New("goclub/chuanglan: NewV1(httpClient, auth) auth.Origin can not be empty string")
	}
	if auth.Account == "" {
		return V1{},xerr.New("goclub/chuanglan: NewV1(httpClient, auth) auth.Account can not be empty string")
	}
	if auth.Password == "" {
		return V1{},xerr.New("goclub/chuanglan: NewV1(httpClient, auth) auth.Account Password not be empty string")
	}
	return V1{
		client: xhttp.NewClient(httpClient),
		auth: auth,
	}, nil
}
// MsgV1SendJson 批量提交接口 https://www.chuanglan.com/document/6110e57909fd9600010209de/6110efa909fd960001020a23
type MsgV1SendJsonRequest struct {
	Msg 		string 		`eg:"【企业】XX项目30分钟内,连续产生2次退款."`
	Phone 		[]string 	`eg:"[]string{\"13411112222\",\"13488882222\"}"`
	SendTime 	time.Time 	`note:"注意传入中国时区的 time.Time"`
	Extend 		string
	Report		bool
	UID 		string
}
type MsgV1SendJsonReply struct {
	Code string `json:"code"` // 返回“0”表示提交成功
	MsgID string `json:"msgId"` // 消息id
	Time string `json:"time"` // 响应时间
	ErrorMsg string `json:"errorMsg"`
}
type HttpDetail struct {
	Request *http.Request
	Response *http.Response
}
func (v V1) MsgV1SendJson(ctx context.Context, data MsgV1SendJsonRequest) (httpDetail HttpDetail, reply MsgV1SendJsonReply, err error) {
	sendData := map[string]interface{}{
		"account": v.auth.Account,
		"password": v.auth.Password,
		"msg": data.Msg,
		"phone": strings.Join(data.Phone, ","),
	}
	if data.SendTime.IsZero() == false {
		sendData["sendtime"] = data.SendTime.In(xtime.LocationChina).Format("20060102150405")
	}
	if data.Report  {
		sendData["report"] = "true"
	}
	if data.Extend != "" {
		sendData["extend"] = data.Extend
	}
	if data.UID != "" {
		sendData["uid"] = data.UID
	}
	var statusCode int
	var bodyClose func () error
	httpDetail.Request, httpDetail.Response, bodyClose, statusCode, err = v.client.Send(ctx, xhttp.POST, v.auth.Origin, "/msg/v1/send/json", xhttp.SendRequest{
		JSON: sendData,
	}) ; if err != nil {
	    return
	}
	defer bodyClose()
	if statusCode != 200 {
		err = xerr.Errorf("goclub/chuanglan: /msg/v1/send/json response status code is %v", statusCode)
		return
	}
	bodyBytes, err := ioutil.ReadAll(httpDetail.Response.Body) ; if err != nil {
	    return
	}
	httpDetail.Response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	err = xjson.Unmarshal(bodyBytes, &reply) ; if err != nil {
	    return
	}
	return
}