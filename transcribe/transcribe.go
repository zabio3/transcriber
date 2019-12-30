package transcribe

type SignalInfo struct {
	Type     string
	Rate     int32
	Channels uint
}

type Results struct {
	Transcripts []*Transcript
}

type Transcript struct {
	Content    string
	Confidence float32
}
