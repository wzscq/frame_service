package flow

import (
	"log"
	"crv/frame/common"
)

type executedNode struct {

}

type flowInstance struct {
	 AppDB string
	 FlowID string
	 InstanceID string
	 UserID string
	 FlowConf *flowConf
	 ExecutedNodes []executedNode
}

func (flow *flowInstance)push(flowRep* flowRep)(*flowResult,int){
	log.Println("start flowInstance push")
	//每个节点的执行都包含两个步骤，启动和结束，
	//先判断当前正在执行的节点（ExecutedNodes中最后一个节点）是否存在，如果存在则加载这个节点并运行
	//如果ExecutedNodes中没有节点，则从FlowConf中获取第一个节点（一般都应该是start节点）加载运行
	
	
	
	log.Println("end flowInstance push")
	return nil,common.ResultSuccess
}