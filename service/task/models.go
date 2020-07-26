package task

import (
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

type Status int

const (
	StatusAvaliabled Status = 1
	StatusInvilid    Status = 0
)

type Rst int

const (
	RstSuccess Rst = 1
	RstFailed  Rst = 0
)

type AgentCollectType int

const (
	//CollInherit        AgentCollectType = 2
	CollInheritPercent AgentCollectType = 0
	CollInheritNum     AgentCollectType = 1
	CollName           AgentCollectType = 2
	CollList           AgentCollectType = 3
)

type StType int

const (
	StageNumber  StType = 1
	StagePercent StType = 0
)

type TType int

const (
	TASKROOT   TType = 0
	HGROUP     TType = 1
	VGROUP     TType = 2
	XGROUP     TType = 3
	TASKCMD    TType = 4
	TASKSCRIPT TType = 5
	SYNCFILE   TType = 6
)

type splitColl uint

const (
	DoSplit    = 1
	DoNotSplit = 0
)

// args type
type ArgsType uint

const (
	COMMONSTR ArgsType = 0 // common string argument
	AGENTFACT ArgsType = 1 // a agent fact key which map to a value
	AGENTTAG  ArgsType = 2 // a agent tag key which map to a value
	//ROOTARGS  ArgsType = 3 // inherit arguments from task root
)

type TaskArg struct {
	gorm.Model
	TaskID   uint     `json:"task_id" gorm:"index"`
	TaskName string   `json:"task_name" gorm:"index;UNIQUE_INDEX:task_name_arg"`
	ArgName  string   `json:"arg_name" gorm:"UNIQUE_INDEX:task_name_arg"`
	ArgType  ArgsType `json:"arg_type"`
	ArgValue string   `json:"arg_value"` // default value, it can be reset when start run
}

type Task struct {
	gorm.Model
	Name            string           `json:"name" gorm:"size:100;index"`
	TkType          TType            `json:"tk_type" gorm:"default:0"`
	ParentTaskID    uint             `json:"parent_task_id" gorm:"default:0;index"`
	SeqNum          uint             `json:"seq_num" gorm:"default:1"`
	ContinueByTask  uint             `json:"continue_by_task" gorm:"default:0"`
	ContinueRst     StageStatus      `json:"status" gorm:"default:5"`
	CollectionType  AgentCollectType `json:"collection_type" gorm:"default:0"`
	CollectionValue string           `json:"collection_value" gorm:"type:MEDIUMTEXT;default:100"`
	SplitParent     splitColl        `json:"child_split_coll" gorm:"default:0"`
	StageType       StType           `json:"stage_type" gorm:"default:0"` // default: percent
	Stages          string           `json:"stages" gorm:"size:500"`
	CreateBy        string           `json:"create_by" gorm:"size:50"`
	Desc            string           `json:"desc" gorm:"size:500"`
	Status          Status           `json:"status" gorm:"default:1"`
	Comments        string           `json:"comments" gorm:"size:500"`
	TaskArgs        []TaskArg        `json:"task_args"`
	ContinueOnFail  uint             `json:"continue_on_fail" gorm:"default:0"` // 0: no, 1: yes.
	ParentXGROUP    uint             `gorm:"-"`                                 // 0: no, 1: yes.
	ParentVGROUP    uint             `gorm:"-"`                                 // 0: no, 1: yes.
	ChildesNum      int              `gorm:"-"`
}

// task type cmd
type TaskCmd struct {
	gorm.Model
	TaskID  uint     `json:"task_id" gorm:"not null;index"`
	Cmd     string   `json:"cmd" gorm:"not null;size:50"`
	Opt     string   `json:"opt" gorm:"size:1000"`
	Timeout uint     `json:"timeout"`
	Args    []InsArg `gorm:"-"`
}

// task type script
type TaskScript struct {
	gorm.Model
	TaskID  uint     `json:"task_id" gorm:"not null;index"`
	Shell   string   `json:"shell" gorm:"size:100"`
	Name    string   `json:"name" gorm:"size:100"`
	RunAs   string   `json:"run_as" gorm:"size:50"`
	Timeout uint     `json:"timeout"`
	Args    []InsArg `gorm:"-"`
}

// task type sync file
type TaskSyncFile struct {
	gorm.Model
	TaskID uint        `json:"task_id" gorm:"not null;index"`
	Src    string      `json:"src"`
	Dst    string      `json:"dst"`
	User   string      `json:"user"`
	Group  string      `json:"group"`
	Perm   os.FileMode `json:"perm"`
}

type StageStatus int

const (
	STAGEREADY       StageStatus = -1
	STAGENEEDCONFIRM StageStatus = 0
	STAGERUNNING     StageStatus = 1
	STAGEPAUSED      StageStatus = 2
	STAGESTOPPED     StageStatus = 3
	STAGEFAILED      StageStatus = 4
	STAGESUCCESS     StageStatus = 5
)

type TaskInstance struct {
	gorm.Model
	Name            string           `json:"name" gorm:"index;size:50;index:name_parent_task"`
	ParentInsID     uint             `json:"parent_ins_id" gorm:"index;default:0;index:parent_task;index:name_parent_task"`
	TaskID          uint             `json:"task_id" gorm:"index;index:parent_task;index:name_parent_task"`
	Stage           int              `json:"stage" gorm:"default:100"`
	StageSeq        uint             `json:"stage_seq"`
	StageAgents     string           `json:"stage_agents" gorm:"size:50"`
	StartAt         time.Time        `json:"start_at"`
	EndAt           time.Time        `json:"end_at"`
	RunBy           string           `json:"run_by" gorm:"size:50"`
	Status          StageStatus      `json:"status" gorm:"default:0"`
	InsMsg          string           `json:"ins_msg" gorm:"size:500"`
	ParentIsX       uint             `json:"parent_is_x" gorm:"default:0"`      // 0: no, 1: yes.
	ContinueOnFail  uint             `json:"continue_on_fail" gorm:"default:0"` // 0: no, 1: yes.
	ParentIsV       uint             `json:"parent_is_v" gorm:"default:0"`      // 0: no, 1: yes.
	BrotherNum      int              `json:"brother_num" gorm:"default:0"`
	CallbackVGROUP  bool             `gorm:"-"`
	TaskStageAgents []TaskStageAgent `json:"task_stage_agents" gorm:"foreignkey:stageName;association_foreignkey:stageAgents"`
	TaskLogs        []TaskLog
}

type InsArg struct {
	gorm.Model
	InsName  string   `json:"ins_name" gorm:"index:ins_task_arg"`
	TaskID   uint     `json:"task_id" gorm:"index:ins_task_arg"`
	ArgName  string   `json:"arg_name"`
	ArgType  ArgsType `json:"arg_type"`
	ArgValue string   `json:"arg_value"`
}

type ChildUsed int

const (
	IsUsed  = 1
	NotUsed = 0
)

// task instance stage ips
type TaskStageAgent struct {
	gorm.Model
	StageName string    `json:"stage_name" gorm:"size:50"`
	IP        string    `json:"ip" gorm:"size:20"`
	ChildUse  ChildUsed `json:"child_use" gorm:"default:0"` // default 0
}

// task run logs
type TaskLog struct {
	// id, ins id (index), task id (index), agent, status，duration, msg
	gorm.Model
	TaskInstanceID uint
}

type basicTask interface {
	QueryByTaskID() error
	start(agents *[]string) error
}
