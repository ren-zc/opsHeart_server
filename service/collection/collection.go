package collection

import "opsHeart/db"

func (ac *AgentCollection) GetAllIPs() ([]string, error) {
	return nil, nil
}

func QueryCollByName(name string) (ac *AgentCollection, err error) {
	err = db.DB.Model(&AgentCollection{}).Where("name = ?", name).First(ac).Error
	return
}
