package chuanglan_test

import (
	"context"
	"github.com/goclub/chuanglan"
	xhttp "github.com/goclub/http"
	xjson "github.com/goclub/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func Test_Client_MsgV1SendJson(t *testing.T) {
	client, err := chuanglan.NewV1(&http.Client{}, chuanglan.AuthV1{
		Origin:   testConfig.Origin,
		Account:  testConfig.Account,
		Password: testConfig.Password,
	}) ; assert.NoError(t, err)
	ctx := context.Background()
	httpDetail, reply, err := client.MsgV1SendJson(ctx, chuanglan.MsgV1SendJsonRequest{
		Msg:      testConfig.MsgV1SendJsonMessage,
		Phone:    testConfig.Tels,
		Report:   true,
	}) ; assert.NoError(t, err)
	xjson.Print("reply", reply)
	if reply.Code != "0" {
		log.Print(xhttp.DumpRequestResponseString(httpDetail.Request, httpDetail.Response, true))
	}
}
