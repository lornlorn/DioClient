package handler

import (
	"app/scheduler"
	"app/utils"
	"html/template"
	"net/http"

	seelog "github.com/cihub/seelog"
)

// TestHandler func(res http.ResponseWriter, req *http.Request)
/*
Route Not Found 404 Page
And
Route "/" Direct To "/index"
*/
func TestHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Router Test : %v", req.URL)

	tmpl, err := template.ParseFiles("./views/test/test.html")
	if err != nil {
		seelog.Errorf("template.ParseFiles Error : %v", err)
		return
	}

	tmpl.Execute(res, nil)
}

// TestAjaxHandler func(res http.ResponseWriter, req *http.Request)
func TestAjaxHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Router Test Ajax : %v", req.URL)
	// key := mux.Vars(req)["key"]

	reqBody := utils.ReadRequestBody2JSON(req.Body)
	seelog.Debugf("Request Body : %v", string(reqBody))

	reqURL := req.URL.Query()
	seelog.Debugf("Request Params : %v", reqURL)

	/*
		var retdata []byte
		switch key {

		default:
		}
	*/
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
	output, err := scheduler.Run(cmd.String(), es, as...)
	if err != nil {
		seelog.Errorf("Command Run Error : %v", err.Error())
		res.Write(utils.GetAjaxRetJSON("9999", err))
		return
	}

	seelog.Debug(string(output))
	res.Write(utils.GetAjaxRetWithDataJSON("0000", nil, string(output)))
	return
}

// NotFoundHandler func(res http.ResponseWriter, req *http.Request)
/*
Route Not Found 404 Page
And
Route "/" Direct To "/index"
*/
func NotFoundHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Router 404 : %v", req.URL)

	if req.URL.Path == "/favicon.ico" {
		seelog.Debugf("Request A favicon [%v]", "./assets/img/favicon.ico")
		http.ServeFile(res, req, "./assets/img/favicon.ico")
		return
	}

	tmpl, err := template.ParseFiles("./views/error/404.html")
	if err != nil {
		seelog.Errorf("template.ParseFiles Error : %v", err)
		return
	}

	tmpl.Execute(res, req.URL)

}
