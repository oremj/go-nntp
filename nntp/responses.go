package nntp

type ListItem struct {
	Name string
	High int64
	Low int64
	Status string
}

type OverviewItem struct {
	MsgId string
	Subject string
	Date string
	Bytes string
}
