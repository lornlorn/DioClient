package scheduler

import (
	"app/api"
	"app/models"
	"app/utils"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cihub/seelog"
)

// Execute func(src string, command string, envs []string, args ...string) ([]byte, error)
func Execute(src string, uuid string, command string, envs []string, args ...string) ([]byte, error) {
	nowTime := time.Now()
	timeFormat := "2006-01-02 15:04:05" // 时间格式化模板
	nowTimeStr := nowTime.Format(timeFormat)
	var sysLog = models.NewLog{
		RunTime: nowTimeStr,
		// RunEnvs: strings.Replace(strings.Trim(fmt.Sprint(es), "[]"), " ", ",", -1),
		RunEnvs: strings.Trim(fmt.Sprint(envs), "[]"),
		RunCmd:  command,
		RunArgs: strings.Trim(fmt.Sprint(args), "[]"),
		ReqSrc:  src,
	}
	output, err := Run(command, envs, args...)
	if err != nil {
		sysLog.RunStatus = "失败"
		sysLog.RunMsg = err.Error()
		return nil, err
	}
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

	logFilePath, err := api.WriteRunLog2File(uuid, string(result))
	if err != nil {
		seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
		sysLog.LogfilePath = fmt.Sprintf("执行结果日志写入失败 : %v", err.Error())
	} else {
		sysLog.LogfilePath = logFilePath
	}

	sysLog.RunStatus = "成功"
	sysLog.RunMsg = string(result)

	err = sysLog.Save()
	if err != nil {
		seelog.Errorf("数据库写入失败 : %v", err.Error())
	}

	return result, nil

}

/*
Run func(command string, args ...string) ([]byte, error)
*/
func Run(command string, envs []string, args ...string) ([]byte, error) {
	// 执行
	cmd := exec.Command(command, args...)
	cmd.Env = envs

	output, err := cmd.StdoutPipe()
	// output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	var out = make([]byte, 0, 1024)
	for {
		tmp := make([]byte, 128)
		n, err := output.Read(tmp)
		out = append(out, tmp[:n]...)
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return nil, err
	}

	return out, nil
}
