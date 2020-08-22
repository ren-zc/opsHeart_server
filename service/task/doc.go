package task

// TASK NOTE

// Task type: ROOT, H_GROUP, V_GROUP, X_GROUP, CMD, SCRIPT, SYNC FILE;
// Status: created, approved, disabled, deleted;

// Type HGROUP: child tasks run one by one;
// Type VGROUP: all child tasks run concomitantly;
//     the seq num of all childes must be 1;
// Type XGROUP: only contain two tasks, they are exclusion,
//     in XGROUP: the firs child will run when condition is success,
//     the second child will run when condition is failed.
//     That is to say only one task to be called in this group;

// Only ROOT need approved, and it is a task entrance;
// ROOT, H_GROUP and V_GROUP has more than one child;
// Type CMD, SCRIPT and SYNC FILE no child;

// Type H_GROUP, V_GROUP, CMD, SCRIPT and SYNC FILE has priority on one level;

// `agent_collect_type`: `list`, `collection_name`, `inherit_num`, `inherit_percent`
// `agent_collect_value`: `["0.0.0.0","1.1.1.1"]` or
// a agent collection name which defined at `AgentCollection` table
// or 100 agents or 20% agents
// `inherit_num`, `inherit_percent`: inherit agents ips from parent.

// step_type: null, num or percent
// null: all agent will be run
// num: run agent by NUM list of steps
// percent: run agent by PERCENT list of steps
// for example
// The total number of agents get from `agent_collect_type` and `agent_collect_value` is 1000
// `step_type`: num; `stages`: [100, 200, 300, 400]; 100 + 200 + 300 + 400 = 1000
// `step_type`: percent; `stages`: [10%, 20%, 30%, 40%]; 10% + 20% + 30% + 40% = 100%

// A task instance will get all agents from `agent_collect_*`;
// If `stage_type` not blank, The task instance will split the number of agents by the `stages`;
// If blank, will not split.

// `continue_by_task`: check the specific task before start;
//     the task must be a brother task
// `continue_by_rst`: the rst of the task of `continue_by_task`
//     default: success
// if continue check failed, the task will be set to failed, and all tasks of later won't be running,
// so if you want tasks of later could be run, you can set `ContinueOnFail` to 1 for this task.

// `ContinueByTask` must be brother node id.

// ContinueOnFail: continue next task even failed.
// ParentXGROUP: parent task is a XGROUP task.

// INSTANCE NOTE

// When instance run, first check if there is a stage need to run, than check next task.
// ParentIsX: this is a instance of task whose parent is a a XGROUP task.
// ContinueOnFail: continue next task even this instance is failed.

// only cmd and script task can have arguments.
