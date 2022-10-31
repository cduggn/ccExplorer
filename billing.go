package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"log"
)

func ToBillable(result *costexplorer.GetCostAndUsageOutput) *Billable {
	billable := &Billable{
		Services: make(map[string]Service),
	}
	for _, v := range result.ResultsByTime {
		billable.Start = *v.TimePeriod.Start
		billable.End = *v.TimePeriod.End
		for _, g := range v.Groups {
			service := Service{}
			for _, k := range g.Keys {
				service.Name = k
			}
			for key, m := range g.Metrics {
				metrics := Metrics{
					Name:   key,
					Amount: *m.Amount,
					Unit:   *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			billable.Services[service.Name] = service
		}
	}
	return billable
}

func (b Billable) print(s Service) {
	fmt.Println("Service: ", s.Name)
	for _, m := range s.Metrics {
		log.Println(m.Name, m.Amount, m.Unit)
	}

}

func (b Billable) billingSetails(s string) Service {
	return b.Services[s]
}

//floatNumAmount, _ := strconv.ParseFloat(*m.Amount, 64)
//fmt.Printf("\t Metric %v \n", key)
//
//
//fmt.Printf("\t\tAmount (Float) %5.10f, Amount (String) %v Unit %v \n", floatNumAmount, *m.Amount, *m.Unit)
