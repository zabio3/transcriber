package transcribe

import (
	"fmt"

	"github.com/krig/go-sox"
)

func GetSampleRate(path string) (*SignalInfo, error) {
	if !sox.Init() {
		return nil, fmt.Errorf("failed to initialize sox")
	}
	defer sox.Quit()

	in := sox.OpenRead(path)
	if in == nil {
		return nil, fmt.Errorf("failed to open input file")
	}
	defer in.Release()

	return &SignalInfo{
		Type:     in.Type(),
		Rate:     int32(in.Signal().Rate()),
		Channels: in.Signal().Channels(),
	}, nil
}
