package scheduler

import (
	"app/models"
	"strings"

	"github.com/cihub/seelog"
	"github.com/rfyiamcool/cronlib"
)

// InitCron func()
func InitCron() error {
	cron := cronlib.New()

	crons, err := models.GetCrons()
	if err != nil {
		seelog.Errorf("Get Cron List Error : %v", err.Error())
		return err
	}

	// seelog.Debug(crons)
	for i, v := range crons {
		seelog.Debugf("Cron Job %v : %v", i, v)
		cronName := v.CronName
		cronSpec := v.CronSpec
		cronEnvs := strings.Split(v.CronEnvs, " ")
		cronCmd := v.CronCmd
		cronArgs := strings.Split(v.CronArgs, " ")
		cronUuid := v.CronUuid

		job, err := cronlib.NewJobModel(
			cronSpec,
			func() {
				Execute("cron", cronUuid, cronCmd, cronEnvs, cronArgs...)
			},
			cronlib.AsyncMode(),
		)
		if err != nil {
			seelog.Errorf("Cron Set Fail : [%v]", cronName)
			return err
		}

		err = cron.Register(cronName, job)
		if err != nil {
			seelog.Errorf("Cron Register Error : %v", err.Error())
			return err
		}
	}

	cron.Start()
	// cron.Join()
	return nil
}
