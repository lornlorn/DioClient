package handler

import (
	"app/models"
	"app/scheduler"
	"app/utils"
	"fmt"
	"net/http"
	"strings"

	seelog "github.com/cihub/seelog"
)

// CronAddHandler func(res http.ResponseWriter, req *http.Request)
func CronAddHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Route CronAdd : %v", req.URL)
	// key := mux.Vars(req)["key"]

	reqBody := utils.ReadRequestBody2JSON(req.Body)
	seelog.Debugf("Request Body : %v", string(reqBody))

	reqURL := req.URL.Query()
	seelog.Debugf("Request Params : %v", reqURL)

	cronName := utils.GetJSONResultFromRequestBody(reqBody, "data.CronName")
	cronSpec := utils.GetJSONResultFromRequestBody(reqBody, "data.CronSpec")
	cronEnvs := utils.ReadJSONData2Array(reqBody, "data.CronEnvs")
	cronCmd := utils.GetJSONResultFromRequestBody(reqBody, "data.CronCmd")
	cronArgs := utils.ReadJSONData2Array(reqBody, "data.CronArgs")
	cronStatus := utils.GetJSONResultFromRequestBody(reqBody, "data.CronStatus")
	cronDesc := utils.GetJSONResultFromRequestBody(reqBody, "data.CronDesc")

	var es []string
	for _, env := range cronEnvs {
		es = append(es, env.String())
	}
	var as []string
	for _, arg := range cronArgs {
		as = append(as, arg.String())
	}

	seelog.Debugf("[%v][%v][%v]", es, cronCmd, as)

	var cron models.NewCron

	cron = models.NewCron{
		CronName:   cronName.String(),
		CronSpec:   cronSpec.String(),
		CronEnvs:   strings.Trim(fmt.Sprint(es), "[]"),
		CronCmd:    cronCmd.String(),
		CronArgs:   strings.Trim(fmt.Sprint(as), "[]"),
		CronStatus: cronStatus.String(),
		CronDesc:   cronDesc.String(),
		CronUuid:   utils.GetUniqueID(),
	}

	var ret []byte

	err := cron.Save()
	if err != nil {
		seelog.Errorf("Set Cron Job -> Write DB Fail : %v", err.Error())
		ret = utils.GetAjaxRetJSON("9999", err)
		res.Write(ret)
		return
	}

	err = scheduler.AddCronJob(cron)
	if err != nil {
		seelog.Errorf("Set Cron Job -> Register Job Fail : %v", err.Error())
		ret = utils.GetAjaxRetJSON("9999", err)
		cron.Delete()
		res.Write(ret)
		return
	}

	seelog.Debug("Set Cron Job Success ...")
	ret = utils.GetAjaxRetJSON("0000", nil)

	res.Write(ret)
	return
}

// CronDeleteHandler func(res http.ResponseWriter, req *http.Request)
func CronDeleteHandler(res http.ResponseWriter, req *http.Request) {

	seelog.Infof("Route CronAdd : %v", req.URL)
	// key := mux.Vars(req)["key"]

	reqBody := utils.ReadRequestBody2JSON(req.Body)
	seelog.Debugf("Request Body : %v", string(reqBody))

	reqURL := req.URL.Query()
	seelog.Debugf("Request Params : %v", reqURL)

	cronName := utils.GetJSONResultFromRequestBody(reqBody, "data.CronName")

	seelog.Debugf("Delete Request CronName : [%v]", cronName.String())

	var cron models.NewCron

	cron = models.NewCron{
		CronName: cronName.String(),
	}

	var ret []byte

	err := cron.Delete()
	if err != nil {
		seelog.Errorf("Unset Cron Job -> Write DB Fail : %v", err.Error())
		ret = utils.GetAjaxRetJSON("9999", err)
		res.Write(ret)
		return
	}

	scheduler.DelCronJob(cron)

	seelog.Debug("Job Stop And Unregister Success ...")
	ret = utils.GetAjaxRetJSON("0000", nil)

	res.Write(ret)
	return
}
