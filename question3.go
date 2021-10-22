package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

/** Question2 is: 設問2の処理
 * サーバが返すpingの応答時間が長くなる場合、サーバが過負荷状態になっていると考えられる。
 * そこで、直近m回の平均応答時間がtミリ秒を超えた場合は、サーバが過負荷状態になっているとみなそう。
 * 設問2のプログラムを拡張して、各サーバの過負荷状態となっている期間を出力できるようにせよ。
 * mとtはプログラムのパラメータとして与えられるようにすること。
 */
func Question3(filepath string, tryCount int64, overloadTime float64, overloadCount int64) {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var failureIpAddrs map[string]*FailureIPAddrDatum = map[string]*FailureIPAddrDatum{}

	var failureList []*FailureServerDatum

	var overloadIPAddrs map[string]*OverloadIPAddrDatum = map[string]*OverloadIPAddrDatum{}

	var overloadList []*OverloadServerDatum

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
			// 故障の場合は以降の処理は行わない
			continue
		}

		resTime, err := strconv.ParseInt(splitted[2], 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		// 過負荷mapに存在するか
		if val, ok := overloadIPAddrs[ipAddr]; ok {
			val.ResponseTime = append(val.ResponseTime, resTime)
			switch {
			case val.ResponseNum+1 < overloadCount:
				// まだ過負荷でない
				val.ResponseNum++
			case overloadCount <= val.ResponseNum+1:
				// 過負荷か確認
				if overloadTime < float64(SumArray(val.ResponseTime))/float64(overloadCount) {
					if val.OverloadListIndex == -1 {
						overloadList = append(overloadList, &OverloadServerDatum{
							ipAddr,
							val.StartTime,
							"",
						})
						overloadIPAddrs[ipAddr].OverloadListIndex = int64(len(overloadList) - 1)
					}

				} else {
					if 0 <= val.OverloadListIndex {
						overloadList[val.OverloadListIndex].EndFailureTime = splitted[0]
						delete(overloadIPAddrs, ipAddr)
					}
				}

				val.ResponseTime = append(val.ResponseTime[2:], resTime)
			}
		} else {
			// 過負荷mapに存在しない場合は入れる
			overloadIPAddrs[ipAddr] = &OverloadIPAddrDatum{
				[]int64{resTime},
				1,
				splitted[0],
				-1,
			}
		}
	}

	// ファイル読込中にエラーが起きた場合のハンドリング
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	ExportFailureList(failureList)
	ExportOverloadList(overloadList)
}
