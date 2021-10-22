package main

// OverloadIPAddrDatum is: 過負荷の可能性があるサーバー情報を格納するmapの要素
type OverloadIPAddrDatum struct {
	ResponseTime      []int64
	ResponseNum       int64
	StartTime         string
	OverloadListIndex int64
}

// OverloadServerDatum is: 過負荷状態のサーバー情報を格納するsliceの要素
type OverloadServerDatum struct {
	IpAddress        string
	StartFailureTime string
	EndFailureTime   string
}

// FailureIPAddrDatum is: 故障状態のサーバー情報を格納するmapの要素
type FailureIPAddrDatum struct {
	FailureListIndex int64
	TryCount         int64
	StartTime        string
}

// FailureServerDatum is: 故障状態のサーバー情報を格納するsliceの要素
type FailureServerDatum struct {
	IpAddress        string
	StartFailureTime string
	EndFailureTime   string
}

type FailureSubnetDatum struct {
	Subnet           string
	StartFailureTime string
	EndFailureTime   string
}

type FailureSubnetMapDatum struct {
	FailureIP      map[string]bool
	FaluireTimeMap map[string]*FaluireTimeDatum
}

type FaluireTimeDatum struct {
	FailureStartTime string
	FailureEndTime   string
	FaliureServerNum int64
}
