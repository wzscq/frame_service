package data

const (
	FIELDTYPE_MANY2MANY = "MANY_TO_MANY"
	FIELDTYPE_MANY2ONE = "MANY_TO_ONE"
	FIELDTYPE_ONE2MANY = "ONE_TO_MANY"
	FIELDTYPE_FILE = "FILE"
)

type QueryRelatedModel interface {
	query(dataRepository DataRepository,parentList *queryResult,refField *field)(int)
}

func getRelatedModelID(modelID string,relatedModelID string)(string){
	if modelID >= relatedModelID {
		return relatedModelID+"_"+modelID
	}
	return modelID+"_"+relatedModelID
}

func GetRelatedModelQuerier(fieldType string,appDB string,modelID string)(QueryRelatedModel){
	if fieldType ==FIELDTYPE_MANY2MANY {
		return &QueryManyToMany{
			AppDB:appDB,
			ModelID:modelID,
		}
	} else if fieldType ==FIELDTYPE_ONE2MANY {
		return &QueryOneToMany{
			AppDB:appDB,
		}
	} else if fieldType ==FIELDTYPE_MANY2ONE {
		return &QueryManyToOne{
			AppDB:appDB,
		}
	} else if fieldType == FIELDTYPE_FILE {
		return &QueryFile{
			AppDB:appDB,
			ModelID:modelID,
		}
	}
	return nil
}

func GetFieldValues(res *queryResult,fieldName string)([]string){
	var valList []string
	for _,row:=range res.List {
		sVal:=row[fieldName].(string)
		valList=append(valList,sVal)
	}
	return valList
}