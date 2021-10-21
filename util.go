package main

import (
	"log"
	"time"
)

func ExportFailureList(failureList []*FailureServerDatum) {
	for _, v := range failureList {
		var endTime string
		if v.EndFailureTime == "" {
			endTime = "故障中"
		} else {
			// pingが返ってきた寸前までが故障期間
			formatted, err := time.Parse("20060102150405", v.EndFailureTime)

			if err != nil {
				log.Printf("ログの日付形式が異常です。ログ: %s", v.EndFailureTime)
			}

			endTime = formatted.Add(-time.Second).Format("20060102150405")
		}

		log.Printf("IPアドレス: %s, 故障期間: %s - %s", v.IpAddress, v.StartFailureTime, endTime)
	}
}

func ExportOverloadList(overloadList []*OverloadServerDatum) {
	for _, v := range overloadList {
		var endTime string
		if v.EndFailureTime == "" {
			endTime = "過負荷継続中"
		} else {
			// pingが返ってきた寸前までが故障期間
			formatted, err := time.Parse("20060102150405", v.EndFailureTime)

			if err != nil {
				log.Printf("ログの日付形式が異常です。ログ: %s", v.EndFailureTime)
			}

			endTime = formatted.Add(-time.Second).Format("20060102150405")
		}

		log.Printf("IPアドレス: %s, 過負荷期間: %s - %s", v.IpAddress, v.StartFailureTime, endTime)
	}
}

func SumArray(slice []int64) int64 {
	var ret int64
	for _, v := range slice {
		ret = ret + v
	}
	return ret
}
