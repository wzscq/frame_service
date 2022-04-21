package data

import (
	"database/sql"
)

type SaveRelatedModel interface {
	save(pID string,dataRepository DataRepository,tx *sql.Tx,modelID string,fieldValue map[string]interface{})(int)
}

func GetRelatedModelSaver(fieldType string,appDB string,userID string,fieldName string)(SaveRelatedModel){
	if fieldType ==FIELDTYPE_MANY2MANY {
		return &SaveManyToMany{
			AppDB:appDB,
			UserID:userID,
		}
	} else if fieldType == FIELDTYPE_FILE {
		return &SaveFile{
			AppDB:appDB,
			UserID:userID,
			FieldName:fieldName,
		}
	}
	/* else if fieldType ==FIELDTYPE_ONE2MANY {
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
		}
	}*/
	return nil
}