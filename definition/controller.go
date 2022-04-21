package definition

import (
	"log"
	"encoding/json"
	"os"
	"io"
	"github.com/gin-gonic/gin"
	"crv/frame/common"
	"net/http"
)

type functionItem struct {
	ID string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
	Operation map[string]interface{} `json:"operation"`
	Icon string `json:"icon"`
}

type functionGroup struct {
	ID string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
	Children []functionItem `json:"children"`
}

type fieldConf struct {
	Field string `json:"field"`
    Name string `json:"name"`
    DataType string `json:"dataType"`
	QuickSearch bool `json:"quickSearch"`
	//以下字段是在关联字段的级联查询中需要携带的参数，用于关联表数据的查询
	FieldType *string `json:"fieldType,omitempty"`
	RelatedModelID *string `json:"relatedModelID,omitempty"`
	RelatedField *string `json:"relatedField,omitempty"`
}

type operationConf struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Params map[string]interface{} `json:"params"`
	Input map[string]interface{} `json:"input"`
	Description string `json:"description"`
	SuccessOperation *operationConf `json:"successOperation,omitempty"`
	ErrorOperation *operationConf `json:"errorOperation,omitempty"`
}

type viewFieldConf struct {
	Field string `json:"field"`
	Width int `json:"width"`
}

type buttonConf struct {
	OperationID string `json:"operationID"`
	Name string `json:"name"`
	Prompt *string `json:"prompt,omitempty"`
}

type toolbarConf struct {
	ShowCount int `json:"showCount"`
	Width int `json:"width"`
	Buttons  []buttonConf `json:"buttons"`
}

type viewToolbarConf struct {
	ListToolbar *toolbarConf `json:"listToolbar,omitempty"`
	RowToolbar  *toolbarConf `json:"rowToolbar,omitempty"`
}

type viewConf struct {
	ViewID string `json:"viewID"`
	Name string `json:"name"`
	Description string `json:"description"`
	Fields []viewFieldConf `json:"fields"`
	Filter map[string]interface{} `json:"filter"`
	Toolbar *viewToolbarConf `json:"toolbar,omitempty"`
}

type modelViewConf struct {
	ModelID string `json:"modelID"`
	Fields []fieldConf `json:"fields"`
	Operations []operationConf `json:"operations"`
	Views []viewConf `json:"views"`
}

type formConf struct {
	FormID string `json:"formID"`
	ColCount int `json:"colCount"`
	RowCount int `json:"rowCount"`
	RowHeight int `json:"rowHeight"`
	Header map[string]interface{} `json:"header"`
	Footer map[string]interface{} `json:"footer"`
	Controls []map[string]interface{} `json:"controls"`
}

type modelFormConf struct {
	ModelID string `json:"modelID"`
	Fields []fieldConf `json:"fields"`
	Operations []operationConf `json:"operations"`
	Forms []formConf `json:"forms"`
}

type getModelViewRep struct {
	ModelID string `json:"modelID"`
}

type getModelFormRep struct {
	ModelID string `json:"modelID"`
	FormID string `json:"formID"`
}

type DefinitionController struct {
	 
}

func (controller *DefinitionController)getUserFunction(c *gin.Context){
	log.Println("start definition getUserFunction")
	//获取用户账号
	//userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	//获取用户角色信息
	//根据角色过滤出功能列表
	formFile := "apps/"+appDB+"/functions/function.json"
	filePtr, err := os.Open(formFile)
	var errorCode int
    if err != nil {
        log.Println("Open file failed [Err:%s]", err.Error())
        errorCode=common.ResultOpenFileError
    }
    defer filePtr.Close()

	var functions []functionGroup
	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&functions)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
		errorCode=common.ResultJsonDecodeError
	} else {
		errorCode=common.ResultSuccess
	}
	
	rsp:=common.CreateResponse(errorCode,functions)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end definition getUserFunction")
}

func (controller *DefinitionController)getModelViewConf(c *gin.Context){
	log.Println("start definition getModelViewConf")
	//获取用户账号
	//userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	//获取用户角色信息
	//根据角色过滤出模型配置
	var errorCode int
	var model modelViewConf
	var rep getModelViewRep	
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorCode=common.ResultWrongRequest
	} else {
		modelFile := "apps/"+appDB+"/models/"+rep.ModelID+".json"
		filePtr, err := os.Open(modelFile)
		if err != nil {
			log.Println("Open file failed [Err:%s]", err.Error())
			errorCode=common.ResultOpenFileError
		}
		defer filePtr.Close()
		// 创建json解码器
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&model)
		if err != nil {
			log.Println("json file decode failed [Err:%s]", err.Error())
			errorCode=common.ResultJsonDecodeError
		} else {
			errorCode=common.ResultSuccess
		}
	}
	rsp:=common.CreateResponse(errorCode,model)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end definition getModelViewConf")
}

func (controller *DefinitionController)getModelFormConf(c *gin.Context){
	log.Println("start definition getModelFormConf")
	//获取用户账号
	//userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	//获取用户角色信息
	//根据角色过滤出模型配置
	var errorCode int
	var model modelFormConf
	var rep getModelFormRep	
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorCode=common.ResultWrongRequest
	} else {
		modelFile := "apps/"+appDB+"/models/"+rep.ModelID+".json"
		filePtr, err := os.Open(modelFile)
		if err != nil {
			log.Println("Open file failed [Err:%s]", err.Error())
			errorCode=common.ResultOpenFileError
		}
		defer filePtr.Close()
		// 创建json解码器
		decoder := json.NewDecoder(filePtr)
		err = decoder.Decode(&model)
		if err != nil {
			log.Println("json file decode failed [Err:%s]", err.Error())
			errorCode=common.ResultJsonDecodeError
		} else {
			errorCode=common.ResultSuccess
		}
		//过滤对应的formID
		var fromRes []formConf
		for _, form := range model.Forms {
			if form.FormID == rep.FormID  {
				fromRes = append(fromRes, form)
			}
		}

		if fromRes == nil {
			errorCode=common.ResultModelFormNotFound	
		} else {
			model.Forms=fromRes
		}
	}
	rsp:=common.CreateResponse(errorCode,model)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end definition getModelFormConf")
}

func (controller *DefinitionController)getAppImage(c *gin.Context){
	log.Println("start definition getAppImage")

	appDB:= c.MustGet("appDB").(string)
	image := c.Param("image")
	imageFile := "apps/"+appDB+"/images/"+image	
	
	f,err:=os.Open(imageFile)
	if err != nil {
		log.Println(err)
		log.Println("end definition getAppImage")	
		return
	}

	io.Copy(c.Writer,f)
	if err := f.Close(); err != nil {
		log.Println(err)
	}
	log.Println("end definition getAppImage")
}

func (controller *DefinitionController) Bind(router *gin.Engine) {
	log.Println("Bind DefinitionController")
	router.POST("/definition/getUserFunction", controller.getUserFunction)
	router.POST("/definition/getModelViewConf", controller.getModelViewConf)
	router.POST("/definition/getModelFormConf", controller.getModelFormConf)
	router.GET("/appimages/:appId/:image", controller.getAppImage)
}