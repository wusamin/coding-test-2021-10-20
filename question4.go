package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/** Question4 is: 設問2の処理
 * ネットワーク経路にあるスイッチに障害が発生した場合、そのスイッチの配下にあるサーバの応答がすべてタイムアウトすると想定される。
 * そこで、あるサブネット内のサーバが全て故障（ping応答がすべてN回以上連続でタイムアウト）している場合は、
 * そのサブネット（のスイッチ）の故障とみなそう。
 * 設問2または3のプログラムを拡張して、各サブネット毎にネットワークの故障期間を出力できるようにせよ。
 */
func Question4(filepath string, tryCount int64) {
	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var failureIpAddrs map[string]*FailureIPAddrDatum = map[string]*FailureIPAddrDatum{}

	var failureList []*FailureServerDatum

	var failureSubnet map[string]*FailureSubnetMapDatum = map[string]*FailureSubnetMapDatum{}

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
					ipSplitted := strings.Split(splitted[1], "/")
					subnetMask, err := strconv.ParseInt(ipSplitted[1], 10, 64)
					if err != nil {
						// サブネットマスクに数値以外が入っている場合はpanicにする
						panic(err)
					}

					subnet := getSubnet(ipSplitted[0], subnetMask)

					// サブネットmapの故障IPmapを確認
					if falilureIP, existsKey := failureSubnet[subnet].FaluireTimeMap[splitted[0]]; existsKey {
						// 既に時刻のキーが存在する場合は
						falilureIP.FaliureServerNum++
					} else {
						failureSubnet[subnet].FaluireTimeMap[splitted[0]] = &FaluireTimeDatum{
							val.StartTime,
							"",
							1,
						}
					}
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
					ipSplitted := strings.Split(splitted[1], "/")
					subnetMask, err := strconv.ParseInt(ipSplitted[1], 10, 64)
					if err != nil {
						// サブネットマスクに数値以外が入っている場合はpanicにする
						panic(err)
					}

					subnet := getSubnet(ipSplitted[0], subnetMask)

					// サブネットmapの故障IPmapを確認
					if falilureIP, existsKey := failureSubnet[subnet].FaluireTimeMap[splitted[0]]; existsKey {
						// 既に時刻のキーが存在する場合は
						falilureIP.FaliureServerNum++
					} else {
						failureSubnet[subnet].FaluireTimeMap[splitted[0]] = &FaluireTimeDatum{
							"",
							val.StartTime,
							1,
						}
					}
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

		ipSplitted := strings.Split(splitted[1], "/")
		subnetMask, err := strconv.ParseInt(ipSplitted[1], 10, 64)
		if err != nil {
			// サブネットマスクに数値以外が入っている場合はpanicにする
			panic(err)
		}

		subnet := getSubnet(ipSplitted[0], subnetMask)

		if v, ok := failureSubnet[subnet]; ok {
			// IPを保持したいだけのため、値は何でもよい
			v.FailureIP[splitted[1]] = true
		} else {
			// サブネットmapにサブネットを追加
			failureSubnet[subnet] = &FailureSubnetMapDatum{
				map[string]bool{splitted[1]: true},
				map[string]*FaluireTimeDatum{},
			}
		}
	}

	// ファイル読込中にエラーが起きた場合のハンドリング
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	ExportFailureList(failureList)
	exportFailureSubnetList(getFailureSubnet(failureSubnet))
}

// getSubnet is: IPアドレスとサブネットマスクより、サブネット部分を取得する
func getSubnet(ipAddr string, subnetMask int64) string {
	s := subnetMask / 8
	ip := strings.Split(ipAddr, ".")

	return strings.Join(ip[0:s], ".")
}

// exportFailureSubnetList is: サブネットの故障リストをログに出力する
func exportFailureSubnetList(failureSubnetList []*FailureSubnetDatum) {
	for _, v := range failureSubnetList {
		var endTime string
		if v.EndFailureTime == "" {
			endTime = "スイッチ故障中"
		} else {
			endTime = v.EndFailureTime
		}
		log.Printf("ネットワーク：%v, 故障期間: %v - %v", v.Subnet, v.StartFailureTime, endTime)
	}
}

// getFailureSubnet is: サブネットmapより、各サブネットの故障期間をリストで返す
func getFailureSubnet(failureSubnet map[string]*FailureSubnetMapDatum) []*FailureSubnetDatum {
	var subnetList []*FailureSubnetDatum

	for subnet := range failureSubnet {

		failureTimes := make([]string, 0, len(failureSubnet[subnet].FaluireTimeMap))
		for k := range failureSubnet[subnet].FaluireTimeMap {
			failureTimes = append(failureTimes, k)
		}
		sort.Strings(failureTimes)

		// サブネット内で記録されたIPの数を取得
		ipNum := int64(len(failureSubnet[subnet].FailureIP))

		var subneMask string

		for k := range failureSubnet[subnet].FailureIP {
			subneMask = strings.Split(k, "/")[1]
			break
		}

		// サブネットが故障状態かどうか
		isBreaking := false

		for _, failureTime := range failureTimes {
			// 故障IPの台数がサブネット内IPリストと一致するか
			if failureSubnet[subnet].FaluireTimeMap[failureTime].FaliureServerNum == ipNum {
				// サブネット故障リストに故障時間を追加
				subnetList = append(subnetList, &FailureSubnetDatum{
					subnet + "/" + subneMask,
					failureSubnet[subnet].FaluireTimeMap[failureTime].FailureStartTime,
					"",
				})
				isBreaking = true
			} else {
				if isBreaking {
					// pingが返ってきた寸前までが故障期間
					formatted, err := time.Parse("20060102150405", failureTime)

					if err != nil {
						log.Printf("ログの日付形式が異常です。ログ: %s", failureTime)
						break
					}

					subnetList[len(subnetList)-1].EndFailureTime = formatted.Add(-time.Second).Format("20060102150405")
				}
			}
		}
	}

	return subnetList

}
