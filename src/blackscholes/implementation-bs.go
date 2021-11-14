package blackscholes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	myUtils "github.com/options-go/src/utils"
)

type bsCalcReq struct {
	Option     Option     `json:"option"`
	Underlying Underlying `json:"underlying"`
}

func CalcBlackScholes(res http.ResponseWriter, req *http.Request) {
	if myUtils.VerifyHTTPMethod("POST", req, res) {

		body, read_err := ioutil.ReadAll(req.Body)
		if read_err != nil { //test that the function stops here then change back to != nil
			res.WriteHeader(500)
			res.Write([]byte("read error"))
			return
		}
		var bsReq bsCalcReq
		un_err := json.Unmarshal(body, &bsReq)
		if un_err != nil {
			res.WriteHeader(500)
			res.Write([]byte("json unmarshal error"))
			return
		}

		bs := NewBlackScholes(&bsReq.Option, &bsReq.Underlying, 0.1)
		bs_json, err := json.Marshal(bs)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte("Internal Server Error: could not marshal bs"))
			return
		}

		res.WriteHeader(200)
		res.Write(bs_json)
	}
}
