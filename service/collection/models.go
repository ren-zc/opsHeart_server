package collection

import "github.com/jinzhu/gorm"

// Agent data, id, agent id, key，value
// agent id, key: composite indexes
type AgentFact struct {
	gorm.Model
	UUID  string `json:"uuid" gorm:"Size:50;index:agent_uuid;UNIQUE_INDEX:agent_key"`
	Key   string `json:"key" gorm:"UNIQUE_INDEX:agent_key;Size:50"`
	Value string `json:"value" gorm:"Size:200"`
	//AgentID int    `json:"agent_id" gorm:"index:agent_id;index:agent_key"`
}

type tagStatus int

const (
	TAGAVALIABLED tagStatus = 0
	TAGDISABLED   tagStatus = 1
)

// settings for agent:
// id, agent id, key, value
// agent id, key: composite indexes
type AgentTag struct {
	gorm.Model
	UUID   string    `json:"uuid" gorm:"Size:50;index:agent_uuid;UNIQUE_INDEX:agent_key"`
	Key    string    `json:"key" gorm:"UNIQUE_INDEX:agent_key;Size:50"`
	Value  string    `json:"value" gorm:"Size:200"`
	Status tagStatus `json:"status" gorm:"default:0"`
	//AgentID int    `json:"agent_id" gorm:"index:agent_id;index:agent_key"`
}

type collType int

const (
	COLLENTRANCE collType = 1
	COLLGROUP    collType = 2
	COLLIST      collType = 3
	COLLFACT     collType = 4
	COLLSETTING  collType = 5
)

type unionType int

const (
	UUNDEFINED unionType = 0
	UAND       unionType = 1
	UOR        unionType = 2
)

type collStatus int

const (
	AVALIABELD collStatus = 0
	DISABLED   collStatus = 1
)

// type:
//   list: give a agent list saved in `value`
//   fact: filter agent by fact
//   tag: filter agent by user tags
//   starts: ip or uuid start with string of `value`
//   ends: ip or uuid end with string of `value`
//   keyword: ip or uuid contain string of `value`
//   equal: ip or uuid is string of `value`
//      when type is starts, ends, keyword or equal, key must be ip or uuid.
//   group: just a container which hold child collection
//   entrance: a root level of a collection
// Define how agents are selected.
type AgentCollection struct {
	// id
	// name(unique index)
	// type: list, fact, tag, starts, ends, keyword, equal, group, entrance
	// key
	// value
	// parent id
	// union type: `and`, `or`, `undefined`, available for group or entrance
	//     group default: `and`
	//     entrance default: `or`
	// status
	// desc
	gorm.Model
	Name      string     `json:"name" gorm:"Size:50;unique_index"`
	ColTYPE   collType   `json:"col_type" gorm:"default:1"`
	Key       string     `json:"key" gorm:"Size:50"`
	Value     string     `json:"value" gorm:"Type:MEDIUMTEXT"`
	ParentID  uint       `json:"parent_id"`
	UnionType unionType  `json:"union_type"`
	Status    collStatus `json:"status"`
	Desc      string     `json:"desc"`
}
