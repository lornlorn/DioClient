package api

import (
	"app/utils"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cihub/seelog"
)

// WriteRunLog2File func(nowTime time.Time, cmd string) error
// Write Execute Result To Log File
func WriteRunLog2File(nowTime time.Time, cmd string, result string) (string, error) {
	logFileFormat := "20060102150405" // 日志文件命名时间格式化
	logFileNameTime := nowTime.Format(logFileFormat)
	UUID := utils.GetUniqueID()
	fileBaseName := filepath.Base(cmd)
	logFileDir := utils.GetConfig("app", "logdir")
	logFileName := fmt.Sprintf("%v_%v-%v", fileBaseName, logFileNameTime, UUID)
	logFilePath := fmt.Sprintf("%v/%v.log", logFileDir, logFileName)
	seelog.Debugf("写入执行结果日志 : %v", logFilePath)
	err := utils.WriteFile(logFilePath, result)
	if err != nil {
		seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
		return "", err
	}
	return logFilePath, nil
}
