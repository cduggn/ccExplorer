package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
)

var (
	conn *storage.CostDataStorage
)

func init() {
	newConnection()
}

func newConnection() {
	conn = &storage.CostDataStorage{}
	err := conn.New("./cloudcost.db")
	if err != nil {
		logger.Error(err.Error())
	}
}

func isEmpty(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}
