package qqcaptcha

//
//
//
//        ***************************     ***************************         *********      ************************
//      *****************************    ******************************      *********      *************************
//     *****************************     *******************************     *********     *************************
//    *********                         *********                *******    *********     *********
//    ********                          *********               ********    *********     ********
//   ********     ******************   *********  *********************    *********     *********
//   ********     *****************    *********  ********************     *********     ********
//  ********      ****************    *********     ****************      *********     *********
//  ********                          *********      ********             *********     ********
// *********                         *********         ******            *********     *********
// ******************************    *********          *******          *********     *************************
//  ****************************    *********            *******        *********      *************************
//    **************************    *********              ******       *********         *********************
//
//

import (
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

type QQCaptcha struct {
	Aid          string
	AppSecretKey string
	Ticket       string
	Randstr      string
	UserIP       string
	Url          string
}

type qct func(qc *QQCaptcha) interface{}

func (qc *QQCaptcha) SetAid(aid string) qct {
	return func(qc *QQCaptcha) interface{} {
		a := qc.Aid
		qc.Aid = aid
		return a
	}
}

func (qc *QQCaptcha) SetSecretKey(sk string) qct {
	return func(qc *QQCaptcha) interface{} {
		a := qc.AppSecretKey
		qc.AppSecretKey = sk
		return a
	}
}

var qqCaptcha *QQCaptcha

func (qc *QQCaptcha) QQCaptchaInit(options ...qct) error {
	q := &QQCaptcha{}
	for _, option := range options {
		option(q)
	}
	qqCaptcha = q
	return nil
}

type Response struct {
	Response  int `json:"response"`
	EvilLevel int `json:"evil_level"`
}

func QqCaptchaVerify(ticket string, randStr string, userIP string) (*http.Response, []error) {
	const QCapUrl = "https://ssl.captcha.qq.com/ticket/verify"
	resp := new(Response)
	res, _, err := gorequest.New().
		Get(QCapUrl).
		Param("aid", qqCaptcha.Aid).
		Param("AppSecretKey", qqCaptcha.AppSecretKey).
		Param("Ticket", ticket).
		Param("Randstr", randStr).
		Param("UserIP", userIP).
		Timeout(time.Minute * 1).Type(gorequest.TypeUrlencoded).EndStruct(resp)
	return res, err
}
