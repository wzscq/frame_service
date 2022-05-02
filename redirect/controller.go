package redirect

import (
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"crv/frame/common"
	"net/http"
	"encoding/json"
	"bytes"
	//"io/ioutil"
)

type commonRep struct {
	ModelID string `json:"modelID"`
	ViewID *string `json:"viewID"`
	To *string `json:"to"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	//Fields *[]field `json:"fields"`
	//Sorter *[]sorter `json:"sorter"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	//Pagination *pagination `json:"pagination"`
}

type apiItem struct {
	Url string `json:"url"`
}

type RedirectController struct {
	 
}

func (controller *RedirectController)getApiUrl(appDB,apiId string)(string,int){
	apiConfigFile := "apps/"+appDB+"/external_api.json"
	filePtr, err := os.Open(apiConfigFile)
	if err != nil {
		log.Println("Open file failed [Err:%s]", err.Error())
		return "",common.ResultOpenFileError
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	apiConf:=map[string]apiItem{}
	err = decoder.Decode(&apiConf)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
		return "",common.ResultJsonDecodeError
	}

	api,ok:=apiConf[apiId]
	if !ok {
		return "",common.ResultNoExternalApiUrl
	}
	
	return api.Url,common.ResultSuccess
}

func (controller *RedirectController)redirect(c *gin.Context){
	log.Println("start redirect ")
	appDB:= c.MustGet("appDB").(string)
	var rep commonRep
	errorcode:=common.ResultSuccess
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,nil)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end redirect with error")
		return
    }
		
	if rep.To==nil{
		errorcode=common.ResultNoExternalApiId
		rsp:=common.CreateResponse(errorcode,nil)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end redirect with error")
		return
	}

	//get url
	postUrl,errorCode:=controller.getApiUrl(appDB,*rep.To)
	if errorCode != common.ResultSuccess {
		rsp:=common.CreateResponse(errorcode,nil)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		return 
	}

	rep.To=nil
	postJson,_:=json.Marshal(rep)
	postBody:=bytes.NewBuffer(postJson)

	resp,err:=http.Post(postUrl,"application/json",postBody)
	if err != nil {
		errorcode=common.ResultPostExternalApiError	
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { 
		log.Println(resp)
		errorcode=common.ResultPostExternalApiError	
	}

	rsp:=common.CreateResponse(errorcode,nil)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end redirect success")
}

func (controller *RedirectController) Bind(router *gin.Engine) {
	log.Println("Bind RedirectController")
	router.POST("/redirect", controller.redirect)
}