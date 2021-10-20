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

	var failureList []struct {
		IpAddress        string
		StartFailureTime string
		EndFailureTime   string
	}

	for fileScanner.Scan() {
		splitted := strings.Split(fileScanner.Text(), ",")

		if len(splitted) != 3 {
			// ログの形式が違う場合のワーニング
			log.Println("")
			continue
		}

		ipAddr := splitted[1]

		// 故障マップ内にIPアドレスがあるか確認
		if val, ok := failureIpAddrs[ipAddr]; ok {
			// ある場合は
			failureList[val].EndFailureTime = splitted[0]
			delete(failureIpAddrs, ipAddr)
		}

		// 故障（pingが通らない）の場合は故障マップにIPアドレスを追加する
		if splitted[2] == "-" {
			failureList = append(failureList, struct {
				IpAddress        string
				StartFailureTime string
				EndFailureTime   string
			}{
				ipAddr,
				splitted[0],
				"",
			})
			failureIpAddrs[ipAddr] = int64(len(failureList) - 1)
		}
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	for _, v := range failureList {
		var endTime string
		if v.EndFailureTime == "" {
			endTime = "故障中"
		} else {
			// Pingが返ってきた寸前までが故障期間
			formatted, err := time.Parse("20060102150405", v.EndFailureTime)
			if err != nil {
				// TODO 日付フォーマット時のエラーのハンドリング
				log.Println(err)
			}

			endTime = formatted.Add(-time.Second).Format("20060102150405")
		}

		log.Printf("IPアドレス: %s, 故障期間: %s - %s", v.IpAddress, v.StartFailureTime, endTime)
	}
}
