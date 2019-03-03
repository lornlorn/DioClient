package handler

import (
	"app/models"
	"app/scheduler"
	"app/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	seelog "github.com/cihub/seelog"
)

// ExecuteHandler func(res http.ResponseWriter, req *http.Request)
func ExecuteHandler(res http.ResponseWriter, req *http.Request) {

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

	nowTime := time.Now()
	timeFormat := "2006-01-02 15:04:05" // 时间格式化模板
	logFileFormat := "20060102150405"   // 日志文件命名时间格式化
	// dateFormat := "20060102"            // 时间格式化模板
	// nowDataStr := nowTime.Format(dateFormat)
	nowTimeStr := nowTime.Format(timeFormat)
	var sysLog = models.NewLog{
		RunTime: nowTimeStr,
		// RunEnvs: strings.Replace(strings.Trim(fmt.Sprint(es), "[]"), " ", ",", -1),
		RunEnvs: strings.Trim(fmt.Sprint(es), "[]"),
		RunCmd:  cmd.String(),
		RunArgs: strings.Trim(fmt.Sprint(as), "[]"),
	}

	var ret []byte
	output, err := scheduler.Run(cmd.String(), es, as...)
	if err != nil {
		seelog.Errorf("Command Run Error : %v", err.Error())
		sysLog.RunStatus = "失败"
		sysLog.RunMsg = err.Error()
		ret = utils.GetAjaxRetJSON("9999", err)
	} else {
		seelog.Debugf("执行结果 : %v", string(output))
		logFileNameTime := nowTime.Format(logFileFormat)
		UUID := utils.GetUniqueID()
		logFileDir := utils.GetConfig("app", "logdir")
		logFileName := fmt.Sprintf("%v-%v", logFileNameTime, UUID)
		logFilePath := fmt.Sprintf("%v/%v.log", logFileDir, logFileName)
		seelog.Debugf("执行结果日志 : %v", logFilePath)
		err := utils.WriteFile(logFilePath, string(output))
		if err != nil {
			seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
			sysLog.LogfilePath = fmt.Sprintf("执行结果日志写入失败 : %v", err.Error())
		} else {
			sysLog.LogfilePath = logFilePath
		}
		sysLog.RunStatus = "成功"
		sysLog.RunMsg = string(output)

		ret = utils.GetAjaxRetWithDataJSON("0000", nil, string(output))
	}

	err = sysLog.Save()
	if err != nil {
		seelog.Errorf("数据库写入失败 : %v", err.Error())
	}

	res.Write(ret)
	return
}