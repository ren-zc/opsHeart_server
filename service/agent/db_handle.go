package agent

import (
	"errors"
	"opsHeart_server/db"
	"time"
)

func (a *Agent) InsertDat() error {
	rst := db.DB.Create(a)
	return rst.Error
}

func (a *Agent) IsExist() bool {
	var tmpA Agent
	db.DB.First(&tmpA, "uuid = ?", a.UUID)
	return tmpA.UUID != ""
}

func (a *Agent) UpdateDat() error {
	//return db.DB.Save(a).Error
	a.UpdatedAt = time.Now()
	return db.DB.Model(&Agent{}).Omit("uuid").UpdateColumns(a).Error
}

func (a *Agent) QueryByUUID(u string) error {
	err := db.DB.Find(a, "uuid = ?", u).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) ChangeStatus() error {
	var t Agent
	d := db.DB.First(&t, "id = ?", a.ID)
	if t.ID == 0 {
		return errors.New("no record found in db")
	}
	d.Update("status", a.Status)
	if a.Status == ACCEPTED {
		d.Update("token", a.Token)
	}
	return d.Error
}

func GetAllUnreg() ([]Agent, error) {
	var all []Agent
	d := db.DB.Where("status = ?", REGISTER).Find(&all)
	return all, d.Error
}
