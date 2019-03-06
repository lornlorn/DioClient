package api

import (
	"app/utils"
	"fmt"

	"github.com/cihub/seelog"
)

// WriteRunLog2File func(nowTime time.Time, cmd string) error
// Write Execute Result To Log File
func WriteRunLog2File(uuid string, result string) (string, error) {
	// logFileFormat := "20060102150405" // 日志文件命名时间格式化
	// logFileNameTime := nowTime.Format(logFileFormat)
	var UUID string
	if uuid == "" || len(uuid) != 32 {
		UUID = utils.GetUniqueID()
	} else {
		UUID = uuid
	}

	// fileBaseName := filepath.Base(cmd)
	logFileDir := utils.GetConfig("app", "logdir")
	logFilePath := fmt.Sprintf("%v/%v.log", logFileDir, UUID)
	seelog.Debugf("写入执行结果日志 : %v", logFilePath)
	err := utils.WriteFile(logFilePath, result)
	if err != nil {
		seelog.Errorf("执行结果日志写入失败 : %v", err.Error())
		return "", err
	}
	return logFilePath, nil
}
