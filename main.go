package main

import (
	"flag"
)

func main() {
	var (
		file          = flag.String("file", "", "処理対象のファイル")
		question      = flag.Int("question", -1, "実行する設問番号")
		pingCount     = flag.Int("ping-count", -1, "サーバーが故障状態とみなされるタイムアウトの回数")
		overloadTime  = flag.Float64("overload-time", -1, "サーバーを過負荷状態とみなす場合の平均応答時間")
		overloadCount = flag.Int("overload-count", -1, "サーバーを過負荷状態とみなす場合の平均応答時間の母数")
	)

	flag.Parse()

	if *file == "" {
		return
	}

	switch *question {
	case 1:
		Question1(*file)
	case 2:
		if *pingCount < 0 {
			return
		}
		Question2(*file, int64(*pingCount))
	case 3:
		if *pingCount < 0 ||
			*overloadTime < 0 ||
			*overloadCount < 0 {
			return
		}
		Question3(*file, int64(*pingCount), float64(*overloadTime), int64(*overloadCount))
	}
}
