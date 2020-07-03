package agent

import (
	"opsHeart/service/collection"
	"time"
)

type regStatus int

const (
	DENIED   regStatus = -1
	BLANK    regStatus = 0
	REGISTER regStatus = 1
	ACCEPTED regStatus = 2
)

type runTimeStatus int

const (
	WORKING runTimeStatus = 1
	OFFLINE runTimeStatus = -1
)

type Agent struct {
	ID           uint          `json:"id" gorm:"primary_key"`
	UUID         string        `json:"uuid" gorm:"Size:50;unique_index"`
	Hostname     string        `json:"hostname" gorm:"Size:100"`
	RemoteAddr   string        `json:"remote_addr" gorm:"Size:50"`
	OsType       string        `json:"os_type" gorm:"Size:50"`
	OsVersion    string        `json:"os_version" gorm:"Size:200"`
	OsArch       string        `json:"os_arch" gorm:"Size:20"`
	AgentVersion string        `json:"agent_version" gorm:"Size:20"`
	Token        string        `json:"token" gorm:"Size:50"`
	Status       regStatus     `json:"status" gorm:"type:tinyint(2);default:0"`
	HbsStatus    runTimeStatus `gorm:"type:tinyint(2);default:0"`
	HbsTime      time.Time     `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AgentFacts   []collection.AgentFact `json:"agent_facts" gorm:"foreignkey:UUID;association_foreignkey:UUID"`
	AgentTags    []collection.AgentTag  `json:"agent_tags" gorm:"foreignkey:UUID;association_foreignkey:UUID"`
}
