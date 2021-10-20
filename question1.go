package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

// Question1 is: 設問1の処理
func Question1(filepath string) {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var failureIpAddrs map[string]int64 = map[string]int64{}

	var failureList []*FailureServerDatum

	for fileScanner.Scan() {
		logLine := fileScanner.Text()
		splitted := strings.Split(logLine, ",")

		if len(splitted) != 3 {
			// ログの形式が違う場合のワーニング
			log.Printf("ログの形式が異常です。ログ: %s", logLine)
			continue
		}

		ipAddr := splitted[1]

		// 故障マップ内にIPアドレスがあるか確認
		if val, ok := failureIpAddrs[ipAddr]; ok {
			// ある場合
			if splitted[2] != "-" {
				failureList[val].EndFailureTime = splitted[0]
				delete(failureIpAddrs, ipAddr)
			}
			// 故障しているIPは復帰するまで処理不要
			continue
		}

		if splitted[2] == "-" {
			// 故障（pingが通らない）の場合は故障マップにIPアドレスを追加する
			failureList = append(failureList, &FailureServerDatum{
				ipAddr,
				splitted[0],
				"",
			})
			failureIpAddrs[ipAddr] = int64(len(failureList) - 1)
		}
	}

	// ファイル読込中にエラーが起きた場合のハンドリング
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	exportFailureList(failureList)
}

func exportFailureList(failureList []*FailureServerDatum) {
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

type FailureServerDatum struct {
	IpAddress        string
	StartFailureTime string
	EndFailureTime   string
}
