package collection

import "opsHeart_server/db"

func (af *AgentFact) IsExist() bool {
	var t AgentFact
	db.DB.Model(&AgentFact{}).First(&t, "uuid = ? and key = ?", af.UUID, af.Key)
	return t.ID != 0
}

func (af *AgentFact) Create() error {
	return db.DB.Model(&AgentFact{}).Create(af).Error
}

func (af *AgentFact) Update() error {
	return db.DB.Model(&AgentFact{}).Omit("uuid", "key").Updates(af).Error
}

func (af *AgentFact) DeleteUUIDAll() error {
	return db.DB.Model(&AgentFact{}).
		Where("uuid = ?", af.UUID).
		Delete(AgentFact{}).Error
}

func (af *AgentFact) DeleteAKey() error {
	return db.DB.Model(&AgentFact{}).
		Where("uuid = ? and key = ?", af.UUID, af.Key).
		Delete(AgentFact{}).Error
}

func (af *AgentFact) QueryValue() (string, error) {
	var t AgentFact
	d := db.DB.Model(&AgentFact{}).
		Where("uuid = ? and key = ?", af.UUID, af.Key).
		First(&t)
	return t.Value, d.Error
}

func (af *AgentFact) QueyAllKeyValueByUUID() ([]AgentFact, error) {
	var afs []AgentFact
	d := db.DB.Model(&AgentFact{}).Where("uuid = ?", af.UUID).Find(&afs)
	return afs, d.Error
}
