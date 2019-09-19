package utils

import (
	"encoding/base64"
	"testing"
)

func TestWxEncryptAndWxDecrypt(t *testing.T) {
	data := "<xml><ToUserName><![CDATA[oia2TjjewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType>  <![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Descript  ion><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>"
	key := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	appid := "wxxxxxxxxxxxxxxxxx"

	encrypt, err := WxEncrypt([]byte(data), key, appid)
	if err != nil {
		t.Error(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(encrypt))

	decrypt, err := WxDecrypt(encrypt, key)
	if err != nil {
		t.Error(err)
	}

	if string(decrypt) != data {
		t.Fail()
	}

	log.Info().Msg(string(decrypt))
}
