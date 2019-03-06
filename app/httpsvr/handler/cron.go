package handler

import (
	"app/models"
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
		seelog.Errorf("数据库写入失败 : %v", err.Error())
		ret = utils.GetAjaxRetJSON("9999", err)
	} else {
		seelog.Debug("定时任务新增成功 ...")
		ret = utils.GetAjaxRetJSON("0000", nil)
	}

	res.Write(ret)
	return
}
