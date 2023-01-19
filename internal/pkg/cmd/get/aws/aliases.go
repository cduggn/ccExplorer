package aws

import (
	"github.com/cduggn/ccexplorer/internal/pkg/service/aws"
)

type Alias struct {
	Name  string
	Query aws.CostAndUsageRequestType
}

//func CreateAlias(req aws.CostAndUsageRequestType) {
//	alias := Alias{
//		Name:  req.Alias,
//		Query: req,
//	}
//	fmt.Println(alias)
//}
