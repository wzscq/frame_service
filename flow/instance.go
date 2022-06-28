package flow

import (
	"log"
	"crv/frame/common"
)

type flowInstance struct {
	 AppDB string `json:"appDB"`
	 FlowID string `json:"flowID"`
	 InstanceID string `json:"instanceID"`
	 UserID string  `json:"UserID"`
	 FlowConf *flowConf `json:"flowConf,omitempty"`
	 ExecutedNodes []executedNode `json:"executedNodes"`
	 Completed bool `json:"completed"`
	 StartTime string `json:"startTime"`
	 EndTime *string `json:"endTime,omitempty"`
}

func (flow *flowInstance)getCurrentNode()(*executedNode){
	exeNodeCount:=len(flow.ExecutedNodes)
	if exeNodeCount>0 {
		return &(flow.ExecutedNodes[exeNodeCount-1])
	}
	return nil
}

func (flow *flowInstance)getStartNode()(*executedNode){
	for _, nodeItem := range (flow.FlowConf.Nodes) {
		if nodeItem.Type == NODE_START {
			return createExecutedNode(nodeItem.ID)
		}
	}
	return nil
}

func (flow *flowInstance)getNextNode(currentNode *executedNode)(*executedNode){
	for _, edgeItem := range (flow.FlowConf.Edges) {
		if edgeItem.Source == currentNode.ID {
			return createExecutedNode(edgeItem.Target)
		}
	}
	return nil
}

func (flow *flowInstance)addExecutedNode(exeNode *executedNode){
	flow.ExecutedNodes=append(flow.ExecutedNodes,*exeNode)
}

func (flow *flowInstance)updateCurrentNode(exeNode *executedNode){
	exeNodeCount:=len(flow.ExecutedNodes)
	flow.ExecutedNodes[exeNodeCount-1]=*exeNode
}

func (flow *flowInstance)runNode(exeNode *executedNode,req *flowRepRsp)(*flowRepRsp,int){
	//根据节点类型，找到对应的节点，然后执行节点
	executor:=getExecutor(exeNode)
	if executor==nil {
		return nil,common.ResultNoExecutorForNodeType
	}
	return executor.run(flow,exeNode,req)
}

func (flow *flowInstance)push(flowRep* flowRepRsp)(*flowRepRsp,int){
	log.Println("start flowInstance push")
	//每个节点的执行都包含两个步骤，启动和结束，
	//先判断当前正在执行的节点（ExecutedNodes中最后一个节点）是否存在，如果存在则加载这个节点并运行
	//如果ExecutedNodes中没有节点，则从FlowConf中获取第一个节点（一般都应该是start节点）加载运行
	//如果当前正在执行的节点执行完成，则从FlowConf中获取下一个待执行节点
	currentNode:=flow.getCurrentNode()
	if currentNode==nil {
		currentNode=flow.getStartNode()
		flow.addExecutedNode(currentNode)
	}

	//循环执行所有同步的node
	for currentNode!=nil {
		result,errorCode:=flow.runNode(currentNode,flowRep)
		if errorCode!= common.ResultSuccess {
			return nil,errorCode
		}
		//更新节点状态
		flow.updateCurrentNode(currentNode)
		//如果执行完，就拿下一个节点继续执行
		if currentNode.Completed {
			currentNode=flow.getNextNode(currentNode)
			flow.addExecutedNode(currentNode)
			//将直接结果参数转换为下一个节点的请求参数
			flowRep=result
		} else {
			//如果没有执行完，说明这个节点是异步节点，直接将结果返回，待后续触发
			return result,common.ResultSuccess
		}
	}
	
	log.Println("end flowInstance push")
	return nil,common.ResultSuccess
}
