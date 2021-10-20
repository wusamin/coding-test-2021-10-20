package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

// Question2 is: 設問2の処理
func Question2(filepath string, tryCount int64) {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var failureIpAddrs map[string]struct {
		FailureListIndex int64
		TryCount         int64
		StartTime        string
	} = map[string]struct {
		FailureListIndex int64
		TryCount         int64
		StartTime        string
	}{}

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
			pingResult := splitted[2]
			if pingResult == "-" {
				switch {
				case val.TryCount+int64(1) == tryCount:
					// 規定回数pingが通らない場合、故障リストにIPアドレスを追加する
					failureList = append(failureList, &FailureServerDatum{
						ipAddr,
						val.StartTime,
						"",
					})
					failureIpAddrs[ipAddr] = struct {
						FailureListIndex int64
						TryCount         int64
						StartTime        string
					}{
						int64(len(failureList) - 1),
						val.TryCount + 1,
						val.StartTime,
					}
				case tryCount < val.TryCount:
					// pingが規定回数以上通ってない場合は何もしない
				case val.TryCount < tryCount:
					failureIpAddrs[ipAddr] = struct {
						FailureListIndex int64
						TryCount         int64
						StartTime        string
					}{
						0,
						val.TryCount + 1,
						val.StartTime,
					}
				}
				continue
			}

			if splitted[2] != "-" {
				switch {
				case tryCount <= val.TryCount:
					failureList[val.FailureListIndex].EndFailureTime = splitted[0]
					delete(failureIpAddrs, ipAddr)
				case val.TryCount < tryCount:
					// 規定回数未満でpingが通った場合は故障mapから削除
					delete(failureIpAddrs, ipAddr)
				}

			}
			// 故障しているIPは復帰するまで処理不要
			continue
		}

		if splitted[2] == "-" {
			// 故障（pingが通らない）の場合は故障マップにIPアドレスを追加する
			failureIpAddrs[ipAddr] = struct {
				FailureListIndex int64
				TryCount         int64
				StartTime        string
			}{
				0,
				1,
				splitted[0],
			}
		}
	}

	// ファイル読込中にエラーが起きた場合のハンドリング
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	ExportFailureList(failureList)
}

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
