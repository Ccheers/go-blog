package alarm

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
	"github.com/go-errors/errors"
	"go-blog/common/mail"
	"regexp"
	"strings"
)

// Define AlarmType to string
// for to check the params is right
type Type string

type MailReceive string

// this are some const params what i defined
// only this can be to input
const (
	AlarmTypeOne   Type = "mail"
	AlarmTypeTwo   Type = "wechat"
	AlarmTypeThree Type = "message"
)

type Param struct {
	Types  Type
	MailTo MailReceive
}

var alarmParam *Param

// Define a closure type to next
type ap func(*Param) (interface{}, error)

// can use this function to set a new value
// but to check it is a right type
func (alarm *Param) SetType(t Type) ap {
	return func(alarm *Param) (interface{}, error) {
		str := strings.Split(string(t), ",")
		if len(str) == 0 {
			return nil, errors.New("you must input a value")
		}
		for _, types := range str {
			s := Type(types)
			_, err := s.IsCurrentType()
			if err != nil {
				return nil, err
			}
		}
		ty := alarm.Types
		alarm.Types = t
		return ty, nil
	}
}

func (alarm *Param) SetMailTo(t MailReceive) ap {
	return func(alarm *Param) (interface{}, error) {
		to := alarm.MailTo
		_, err := t.CheckIsNull()
		if err != nil {
			return nil, err
		}
		_, err = t.MustMailFormat()
		if err != nil {
			return nil, err
		}
		alarm.MailTo = t
		return to, nil
	}
}

// alarm receive account can not null
func (t MailReceive) CheckIsNull() (MailReceive, error) {
	if len(t) == 0 {
		return "", errors.New("value can not be null")
	}
	return t, nil
}

// alarm receive account must be mail format
func (t MailReceive) MustMailFormat() (MailReceive, error) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", string(t)); !m {
		return "", errors.New("value format is not right")
	}
	return t, nil
}

// judge it is a right type what i need
// if is it a wrong type, i must return a panic to above
func (at Type) IsCurrentType() (Type, error) {
	switch at {
	case AlarmTypeOne:
		return at, nil
	case AlarmTypeTwo:
		return at, nil
	case AlarmTypeThree:
		return at, nil
	default:
		return at, errors.New("the alarm type is error")
	}
}

// implementation value
func (alarm *Param) AlarmInit(options ...ap) error {
	q := &Param{}
	for _, option := range options {
		_, err := option(q)
		if err != nil {
			return err
		}
	}
	alarmParam = q
	return nil
}

func Alarm(content string) {
	types := strings.Split(string(alarmParam.Types), ",")
	var err error
	for _, a := range types {
		switch Type(a) {
		case AlarmTypeOne:
			if alarmParam.MailTo == "" {
				break
			}
			err = mail.SendMail(string(alarmParam.MailTo), "报警", content)
			break
		case AlarmTypeTwo:
			break
		case AlarmTypeThree:
			break
		}
		if err != nil {
		}
	}
}
