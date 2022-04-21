package definition

import (
	"log"
	"encoding/json"
	"os"
	"crv/frame/common"
)

type ModelConf struct {
	ModelID string `json:"modelID"`
	Fields []fieldConf `json:"fields"`
}

func GetModelConf(appDB string,modelID string)(*ModelConf,int){
	modelFile := "apps/"+appDB+"/models/"+modelID+".json"
	filePtr, err := os.Open(modelFile)
	if err != nil {
		log.Println("Open file failed [Err:%s]", err.Error())
		return nil,common.ResultOpenFileError
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	modelConf:=&ModelConf{}
	err = decoder.Decode(modelConf)
	if err != nil {
		log.Println("json file decode failed [Err:%s]", err.Error())
		return nil,common.ResultJsonDecodeError
	}
	
	return modelConf,common.ResultSuccess
}

