package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

/** Question2 is: 設問2の処理
 * ネットワークの状態によっては、一時的にpingがタイムアウトしても、一定期間するとpingの応答が復活することがあり、
 * そのような場合はサーバの故障とみなさないようにしたい。
 * N回以上連続してタイムアウトした場合にのみ故障とみなすように、設問1のプログラムを拡張せよ。
 * Nはプログラムのパラメータとして与えられるようにすること。
 */
func Question2(filepath string, tryCount int64) {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var failureIpAddrs map[string]*FailureIPAddrDatum = map[string]*FailureIPAddrDatum{}

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
					failureIpAddrs[ipAddr].TryCount = val.TryCount + 1
					failureIpAddrs[ipAddr].FailureListIndex = int64(len(failureList) - 1)
				case tryCount < val.TryCount:
					// pingが規定回数以上通ってない場合は何もしない
				case val.TryCount < tryCount:
					failureIpAddrs[ipAddr].TryCount = val.TryCount + 1
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
			failureIpAddrs[ipAddr] = &FailureIPAddrDatum{
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
