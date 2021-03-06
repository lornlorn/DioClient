package models

import (
    "app/utils"
    "errors"

    "github.com/cihub/seelog"
)

/*
SysCron struct map to table sys_cron
*/
type SysCron struct {
    CronId     int    `xorm:"INTEGER NOT NULL UNIQUE PK"`
    CronName   string `xorm:"VARCHAR(256) NOT NULL UNIQUE"`
    CronSpec   string `xorm:"VARCHAR(128) NOT NULL"`
    CronEnvs   string `xorm:"VARCHAR(512)"`
    CronCmd    string `xorm:"VARCHAR(512)   NOT NULL"`
    CronArgs   string `xorm:"VARCHAR(512)"`
    CronStatus string `xorm:"VARCHAR(16)   NOT NULL"`
    CronDesc   string `xorm:"VARCHAR(1024)"`
    CronUuid   string `xorm:"VARCHAR(32) NOT NULL UNIQUE"`
}

/*
NewCron struct map to table sys_cron without column Id
*/
type NewCron struct {
    CronName   string `xorm:"VARCHAR(256) NOT NULL UNIQUE"`
    CronSpec   string `xorm:"VARCHAR(128) NOT NULL"`
    CronEnvs   string `xorm:"VARCHAR(512)"`
    CronCmd    string `xorm:"VARCHAR(512)   NOT NULL"`
    CronArgs   string `xorm:"VARCHAR(512)"`
    CronStatus string `xorm:"VARCHAR(16)   NOT NULL"`
    CronDesc   string `xorm:"VARCHAR(1024)"`
    CronUuid   string `xorm:"VARCHAR(32) NOT NULL UNIQUE"`
}

/*
TableName xorm mapper
NewComponent struct map to table tb_component
*/
func (cron NewCron) TableName() string {
    return "sys_cron"
}

// Save insert method
func (cron NewCron) Save() error {
    affected, err := utils.Engine.Insert(cron)
    if err != nil {
        return err
    }
    seelog.Debugf("%v insert : %v", affected, cron)

    return nil
}

// Delete method
func (cron NewCron) Delete() error {
    affected, err := utils.Engine.Where("cron_name=? and cron_uuid=?", cron.CronName, cron.CronUuid).Delete(cron)
    if err != nil {
        return err
    }
    seelog.Debugf("%v delete : %v", affected, cron)

    if affected == 0 {
        return errors.New("0 updated, nothing changed")
    }

    return nil
}

// UpdateByUUID method
func (cron NewCron) UpdateByUUID() error {
    affected, err := utils.Engine.Where("cron_name=? and cron_uuid=?", cron.CronName, cron.CronUuid).Update(cron)
    if err != nil {
        // seelog.Errorf("utils.Engine.ID.Update Error : %v", err)
        return err
    }
    seelog.Debugf("%v update : %v", affected, cron)

    if affected == 0 {
        return errors.New("0 updated, nothing changed")
    }

    return nil
}

/*
GetCrons func() ([]SysCron, error)
*/
func GetCrons() ([]SysCron, error) {
    crons := make([]SysCron, 0)

    if err := utils.Engine.Find(&crons); err != nil {
        // seelog.Errorf("utils.Engine.Find Error : %v", err)
        return nil, err
    }

    return crons, nil
}

/*
DelCronByKeys func(cronID int, cronName string, cronUUID string) error
*/
func DelCronByKeys(cronID int, cronName string, cronUUID string) error {
    cron := new(SysCron)
    cron.CronId = cronID
    cron.CronName = cronName
    cron.CronUuid = cronUUID

    affected, err := utils.Engine.Delete(cron)
    if err != nil {
        // seelog.Errorf("utils.Engine.Delete Error : %v", err)
        return err
    }
    seelog.Debugf("%v delete : %v", affected, cron)

    return nil
}
