package nntp

type ListItem struct {
	Name   string
	High   int64
	Low    int64
	Status string
}

type OverviewItem struct {
	MsgNum  string
	Subject string
	From    string
	Date    string
	MsgId   string
	Bytes   string
}
