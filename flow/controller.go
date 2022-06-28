package flow

import (
	"log"
	"github.com/gin-gonic/gin"
	"crv/frame/common"
	"net/http"
)

//流中的请求和应答采用相同的结构，便于不同节点间传递数据
type flowRepRsp struct {
	ModelID *string `json:"modelID"`
	ViewID *string `json:"viewID"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	Fields *interface{} `json:"fields"`
	Sorter *interface{} `json:"sorter"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	Pagination *interface{} `json:"pagination"`
	FlowID *string `json:"flowID"`
	FlowInstanceID *string `json:"flowInstanceID"`
	Total *int `json:"total,omitempty"`
	Operation *[]map[string]interface{} `json:"operation,omitempty"`
	Value *string `json:"value,omitempty"`
}

type FlowController struct {
	InstanceRepository FlowInstanceRepository
}

func (controller *FlowController)start(c *gin.Context){
	log.Println("start FlowController start")
	//获取相关参数
	//userRoles:= c.MustGet("userRoles").(string)
	userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)

	var rep flowRepRsp
	var result *flowRepRsp
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.ResultWrongRequest,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end FlowController start")
		return
	} 

	//启动流必须指定流的ID
	if rep.FlowID == nil || len(*(rep.FlowID))==0 {
		rsp:=common.CreateResponse(common.ResultStartFlowWithoutID,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end FlowController start")
		return
	}

	//创建流
	flowInstance,errorCode:=createInstance(appDB,*(rep.FlowID),userID)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(errorCode,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end FlowController start")
		return
	}
	//执行流
	result,errorCode=flowInstance.push(&rep)

	//如果流中存在待执行的节点，则保存流实例到缓存
	if !flowInstance.Completed {
		err:=controller.InstanceRepository.saveInstance(flowInstance)
		if err!=nil {
			rsp:=common.CreateResponse(common.ResultCacheFlowInstanceError,result)
			c.IndentedJSON(http.StatusOK, rsp.Rsp)
			log.Println("end FlowController start")
			return
		}
	}
	
	rsp:=common.CreateResponse(errorCode,result)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end FlowController start")
}

func (controller *FlowController)terminate(c *gin.Context){

}

func (controller *FlowController)push(c *gin.Context){

}

func (controller *FlowController) Bind(router *gin.Engine) {
	log.Println("Bind FlowController")
	router.POST("/flow/start", controller.start)
	router.POST("/flow/terminate", controller.terminate)
	router.POST("/flow/push", controller.push)
}