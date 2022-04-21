package data

import (
	"log"
	"github.com/gin-gonic/gin"
	"crv/frame/common"
	"net/http"
)

type commonRep struct {
	ModelID string `json:"modelID"`
	ViewID *string `json:"viewID"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	Fields *[]field `json:"fields"`
	Sorter *[]sorter `json:"sorter"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	Pagination *pagination `json:"pagination"`
}

type DataController struct {
	DataRepository DataRepository
}

func (controller *DataController) query(c *gin.Context) {
	log.Println("start data query")
	//获取用户账号
	//userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	var rep commonRep
	var errorcode int
	var result *queryResult
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorcode=common.ResultWrongRequest
    } else {
		query:=&Query{
			ModelID:rep.ModelID,
			ViewID:rep.ViewID,
			Pagination:rep.Pagination,
			Filter:rep.Filter,
			Fields:rep.Fields,
			AppDB:appDB,
			Sorter:rep.Sorter,
		}
		result,errorcode=query.Execute(controller.DataRepository)
	}
	rsp:=common.CreateResponse(errorcode,result)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end data query")
}

func (controller *DataController) save(c *gin.Context) {
	log.Println("start data save")
	//获取用户账号
	userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	var rep commonRep
	var errorcode int
	var result *saveResult
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end data save with error")
		return
    }
		
	if rep.List==nil{
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end data save with error")
		return
	}

	save:=&Save{
		ModelID:rep.ModelID,
		AppDB:appDB,
		UserID:userID,
		List:rep.List,
	}
	result,errorcode=save.Execute(controller.DataRepository)
	rsp:=common.CreateResponse(errorcode,result)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end data save success")
}

func (controller *DataController) delete(c *gin.Context) {
	log.Println("start data delete")
	//获取用户账号
	userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	var rep commonRep
	var errorcode int
	var result *map[string]interface {} = nil
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end data delete with error")
		return
    }

	if rep.SelectedRowKeys == nil {
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusOK, rsp.Rsp)
		log.Println("end data delete with error")
		return
	}

	delete:=&Delete{
		ModelID:rep.ModelID,
		AppDB:appDB,
		UserID:userID,
		SelectedRowKeys:rep.SelectedRowKeys,
	}
	result,errorcode=delete.Execute(controller.DataRepository)
	rsp:=common.CreateResponse(errorcode,result)
	c.IndentedJSON(http.StatusOK, rsp.Rsp)
	log.Println("end data save")
}

func (controller *DataController) update(c *gin.Context) {

}

func (controller *DataController)download(c *gin.Context) {
	log.Println("start data download")
	//获取用户账号
	userID:= c.MustGet("userID").(string)
	appDB:= c.MustGet("appDB").(string)
	var rep commonRep
	var errorcode int
	var result *map[string]interface {} = nil
	if err := c.BindJSON(&rep); err != nil {
		log.Println(err)
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusInternalServerError, rsp.Rsp)
		log.Println("end data download with error")
		return
    }

	if rep.List == nil || len(*(rep.List))<=0 {
		errorcode=common.ResultWrongRequest
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusInternalServerError, rsp.Rsp)
		log.Println("end data download with error")
		return
	}

	download:=&Download{
		ModelID:rep.ModelID,
		AppDB:appDB,
		UserID:userID,
		List:rep.List,
	}

	c.Header("Content-Type", "application/octet-stream")
    c.Header("Content-Disposition", "attachment; filename=downloadFile")
    c.Header("Content-Transfer-Encoding", "binary")

	errorCode:=download.Execute(c.Writer)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(errorcode,result)
		c.IndentedJSON(http.StatusInternalServerError, rsp.Rsp)
	}
	log.Println("end data download")
}

func (controller *DataController) Bind(router *gin.Engine) {
	log.Println("Bind DataController")
	router.POST("/data/query", controller.query)
	router.POST("/data/save", controller.save)
	router.POST("/data/delete", controller.delete)
	router.POST("/data/update", controller.update)
	router.POST("/data/download", controller.download)
}