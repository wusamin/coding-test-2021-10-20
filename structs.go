package main

// OverloadIPAddrDatum is: 過負荷の可能性があるサーバー情報を格納するmapの要素
type OverloadIPAddrDatum struct {
	// 応答時間のslice
	ResponseTime []int64

	// 応答回数
	ResponseNum int64

	// 故障開始時刻（文字列）
	StartTime string

	// 過負荷状態のリスト内でのインデックス
	OverloadListIndex int64
}

// OverloadServerDatum is: 過負荷状態のサーバー情報を格納するsliceの要素
type OverloadServerDatum struct {
	// サーバーのIPアドレス
	IpAddress string

	// 過負荷開始時刻
	StartFailureTime string

	// 過負荷終了時刻
	EndFailureTime string
}

// FailureIPAddrDatum is: 故障状態のサーバー情報を格納するmapの要素
type FailureIPAddrDatum struct {
	// 故障サーバーリスト内でのインデックス
	FailureListIndex int64

	// タイムアウトになった回数
	TryCount int64

	// 故障開始時刻
	StartTime string
}

// FailureServerDatum is: 故障状態のサーバー情報を格納するsliceの要素
type FailureServerDatum struct {
	// サーバーのIPアドレス
	IpAddress string

	// 故障開始時刻
	StartFailureTime string

	// 故障終了時刻
	EndFailureTime string
}

// FailureSubnetDatum is: 故障状態のサブネットの情報を格納するsliceの要素
type FailureSubnetDatum struct {
	// サブネットの属するアドレス
	Subnet string

	// 故障開始時刻
	StartFailureTime string

	// 故障終了時刻
	EndFailureTime string
}

// FailureSubnetMapDatum is: サブネットの情報を格納するmapの要素
type FailureSubnetMapDatum struct {
	// サブネットに属するIPアドレス
	FailureIP map[string]bool

	// サブネット内のサーバーの故障情報
	FaluireTimeMap map[string]*FaluireTimeDatum
}

// FaluireTimeDatum is: サブネット情報を格納するmapのうち、各サーバーの故障情報を格納するsliceの要素
type FaluireTimeDatum struct {
	// 故障開始時刻
	FailureStartTime string

	// 故障終了時刻
	FailureEndTime string

	// その時間に何台故障していたか
	FaliureServerNum int64
}
