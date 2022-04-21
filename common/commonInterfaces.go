package common

type LoginCache interface {
	SetCache(userID string,token string,dbName string)(error)
	GetUserID(token string)(string,error)
	GetAppDB(token string)(string,error)
	RemoveUser(userID string)
}

type AppCache interface {
	GetAppDB(appID string)(string,error)
}