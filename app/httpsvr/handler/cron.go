package handler

import (
	"app/utils"
	"net/http"

	seelog "github.com/cihub/seelog"
)

// CronHandler func(res http.ResponseWriter, req *http.Request)
func CronHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Route Execute : %v", req.URL)
	// key := mux.Vars(req)["key"]

	reqBody := utils.ReadRequestBody2JSON(req.Body)
	seelog.Debugf("Request Body : %v", string(reqBody))

	reqURL := req.URL.Query()
	seelog.Debugf("Request Params : %v", reqURL)

	envs := utils.ReadJSONData2Array(reqBody, "data.envs")
	cmd := utils.GetJSONResultFromRequestBody(reqBody, "data.cmd")
	args := utils.ReadJSONData2Array(reqBody, "data.args")

	var es []string
	for _, env := range envs {
		es = append(es, env.String())
	}
	var as []string
	for _, arg := range args {
		as = append(as, arg.String())
	}

	seelog.Debugf("[%v][%v][%v]", es, cmd, as)

	res.Write(nil)
	return
}
