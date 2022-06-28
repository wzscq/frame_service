package flow

import (
    "time"
)

type executedNode struct {
	ID string `json:"id"`
	Completed bool `json:"completed"`
	StartTime string `json:"startTime"`
	EndTime *string `json:"endTime,omitempty"`
	Data *interface{} `json:"data,omitempty"`
}

type nodeExecutor interface {
	run(instance *flowInstance,node *executedNode,req *flowRepRsp)(*flowRepRsp,int)
}

const (
	NODE_START = "start"
	NODE_FORM = "form"
	NODE_SAVE = "save"
	NODE_END = "end"
)

func createExecutedNode(id string)(*executedNode){
	return &executedNode{
		ID:id,
		Completed:false,
		StartTime:time.Now().Format("2006-01-02 15:04:05"),
	}
}

func getExecutor(exeNode *executedNode)(nodeExecutor){
	
	return nil
}
