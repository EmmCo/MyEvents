package dblayer

import (
	"awesomeProject/peresistence"
	"awesomeProject/peresistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB  DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (peresistence.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}