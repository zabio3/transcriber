package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

	"github.com/zabio3/transcriber/transcribe"
)

var speechClient *speech.Client

func NewSpeechClient(ctx context.Context) (*speech.Client, error) {
	if speechClient != nil {
		return speechClient, nil
	}

	var err error
	speechClient, err = speech.NewClient(ctx)
	return speechClient, err
}

func RecognizeSpeech(ctx context.Context, in []byte, si *transcribe.SignalInfo, lang string) ([]*transcribe.Transcript, error) {
	if _, err := NewSpeechClient(ctx); err != nil {
		return nil, fmt.Errorf("failed to create cloud speech client (err: %s)", err)
	}

	enc, err := getAudioEnc(si.Type)
	if err != nil {
		return nil, err
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        enc,
			SampleRateHertz: si.Rate,
			LanguageCode:    lang,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: in,
			},
		},
	}

	resp, err := speechClient.Recognize(ctx, req)
	if err != nil {
		return nil, err
	}

	var trs []*transcribe.Transcript
	for _, rst := range resp.Results {
		for _, alt := range rst.Alternatives {
			trs = append(trs, &transcribe.Transcript{
				Content:    alt.Transcript,
				Confidence: alt.Confidence,
			})
		}
	}

	return trs, nil
}

func getAudioEnc(audioType string) (speechpb.RecognitionConfig_AudioEncoding, error) {
	switch audioType {
	case "wav":
		return speechpb.RecognitionConfig_LINEAR16, nil
	case "flac":
		return speechpb.RecognitionConfig_FLAC, nil
	default:
		return speechpb.RecognitionConfig_ENCODING_UNSPECIFIED, fmt.Errorf("unknown extension: " + audioType)
	}
}
