package hashid

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
	"github.com/speps/go-hashids"
)

const (
	defaultSalt   = "i must add a salt what is only for me"
	defaultLength = 8
)

type Params struct {
	Salt      string
	MinLength int
}

var hashIdParams *Params

func (hd *Params) SetHashIdSalt(salt string) func(*Params) interface{} {
	return func(hd *Params) interface{} {
		hs := hd.Salt
		hd.Salt = salt
		return hs
	}
}

func (hd *Params) SetHashIdLength(minLength int) func(*Params) interface{} {
	return func(hd *Params) interface{} {
		ml := hd.MinLength
		hd.MinLength = minLength
		return ml
	}
}

func (hd *Params) HashIdInit(options ...func(*Params) interface{}) (*hashids.HashID, error) {
	q := &Params{
		Salt:      defaultSalt,
		MinLength: defaultLength,
	}
	for _, option := range options {
		option(q)
	}
	hashIdParams = q
	hds := hashids.NewData()
	hds.Salt = hashIdParams.Salt
	hds.MinLength = hashIdParams.MinLength
	h, err := hashids.NewWithData(hds)
	if err != nil {
		return nil, err
	}
	return h, nil
}
