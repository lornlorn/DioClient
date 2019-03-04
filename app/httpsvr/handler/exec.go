package handler

import (
	"app/api"
	"app/models"
	"app/scheduler"
	"app/utils"
	"fmt"
	"net/http"
	"runtime"
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
	// dateFormat := "20060102"            // 时间格式化模板
	// nowDataStr := nowTime.Format(dateFormat)
	nowTimeStr := nowTime.Format(timeFormat)
	var sysLog = models.NewLog{
		RunTime: nowTimeStr,
		// RunEnvs: strings.Replace(strings.Trim(fmt.Sprint(es), "[]"), " ", ",", -1),
		RunEnvs: strings.Trim(fmt.Sprint(es), "[]"),
		RunCmd:  cmd.String(),
		RunArgs: strings.Trim(fmt.Sprint(as), "[]"),
		ReqSrc:  "http",
	}

	var ret []byte
	output, err := scheduler.Run(cmd.String(), es, as...)
	if err != nil {
		seelog.Errorf("Command Run Error : %v", err.Error())
		sysLog.RunStatus = "失败"
		sysLog.RunMsg = err.Error()
		ret = utils.GetAjaxRetJSON("9999", err)
	} else {
		var result []byte
		if runtime.GOOS == "windows" {
			seelog.Debug("Decode GBK to UTF-8 ...")
			result, err = utils.DecodeGBK2UTF8(output)
			if err != nil {
				seelog.Errorf("Decode GBK to UTF-8 Error : %v", err.Error())
				result = output
			}
		} else {
			result = output
		}
		seelog.Debugf("执行结果 : %v", string(result))
		logFilePath, err := api.WriteRunLog2File(nowTime, cmd.String(), string(result))
		if err != nil {
			sysLog.LogfilePath = fmt.Sprintf("执行结果日志写入失败 : %v", err.Error())
		} else {
			sysLog.LogfilePath = logFilePath
		}

		sysLog.RunStatus = "成功"
		sysLog.RunMsg = string(result)

		ret = utils.GetAjaxRetWithDataJSON("0000", nil, string(result))
	}

	err = sysLog.Save()
	if err != nil {
		seelog.Errorf("数据库写入失败 : %v", err.Error())
	}

	res.Write(ret)
	return
}
