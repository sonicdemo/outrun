package bgtasks

import (
	"log"

	"github.com/RunnersRevival/outrun/consts"
	"github.com/RunnersRevival/outrun/db/dbaccess"
)

func TouchAnalyticsDB() {
	err := dbaccess.Set(consts.DBBucketAnalytics, "touch", []byte{})
	if err != nil {
		log.Println("[ERR] Unable to touch " + consts.DBBucketAnalytics + ": " + err.Error())
	}
}
