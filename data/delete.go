package data

import (
	"crv/frame/common"
	"log"
	"database/sql"
)

type Delete struct {
	ModelID string `json:"modelID"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"` 
	AppDB string `json:"appDB"`
	UserID string `json:"userID"`
}

func (delete *Delete)idListToString(idList *[]string)(string,int){
	strIDs:=""
	if len(*idList)<=0 {
		return "",common.ResultNoIDWhenDelete
	}

	for _, strID := range *idList {
		strIDs=strIDs+"'"+strID+"',"
	}
	//去掉末尾的逗号
	strIDs=strIDs[0:len(strIDs)-1]
	return strIDs,common.ResultSuccess
}

func (delete *Delete) delete(dataRepository DataRepository,tx *sql.Tx,modelID string,idList *[]string)(*map[string]interface {},int) {
	//获取所有待删数据ID列表字符串，类似：'id1','id2'
	strIDs,errorCode:=delete.idListToString(idList)
	if errorCode!=common.ResultSuccess {
		return nil,errorCode
	}

	sql:="delete from "+delete.AppDB+"."+modelID+" where id in ("+strIDs+")"
	_,rowCount,err:=dataRepository.execWithTx(sql,tx)
	if err != nil {
		return nil,common.ResultSQLError
	}
	result := map[string]interface{}{}		
	result["count"]=rowCount
	result["modelID"]=modelID

	//还要删掉和当前模型相关联的中间表的数据
	dr:=DeleteReleated{
		ModelID:modelID,
		AppDB:delete.AppDB,
		UserID:delete.UserID,
		IdList:idList,
	}
	errorCode=dr.Execute(dataRepository,tx)

	return &result,errorCode
}

func (delete *Delete) Execute(dataRepository DataRepository)(*map[string]interface {},int) {
	//开启事务
	tx,err:= dataRepository.begin()
	if err != nil {
		log.Println(err)
		return nil,common.ResultSQLError
	}
	//执行保存动作
	result,errorcode:=delete.delete(dataRepository,tx,delete.ModelID,delete.SelectedRowKeys)
	if errorcode == common.ResultSuccess {
		//提交事务
		if err := tx.Commit(); err != nil {
			log.Println(err)
			errorcode=common.ResultSQLError
		}
	} else {
		tx.Rollback()
	}
	return result,errorcode
}

