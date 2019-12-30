package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

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

func StreamRecognizeSpeech(ctx context.Context, in []byte, lang string) error {
	if _, err := NewSpeechClient(ctx); err != nil {
		return fmt.Errorf("failed to create cloud speech client (err: %s)", err)
	}
	stream, err := speechClient.StreamingRecognize(ctx)
	if err != nil {
		return err
	}

	req := &speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    lang,
				},
				InterimResults: true,
			},
		},
	}

	if err := stream.Send(req); err != nil {
		return err
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				if err := stream.Send(&speechpb.StreamingRecognizeRequest{
					StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
						AudioContent: buf[:n],
					}}); err != nil {
					panic(err)
					return
				}
			}
			if err == io.EOF {
				if err := stream.CloseSend(); err != nil {
					panic(err)
				}
				return
			}

			if err != nil {
				panic(err)
			}
		}
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		if err := resp.Error; err != nil {
			// Workaround while the API doesn't give a more informative error.
			if err.Code == 3 || err.Code == 11 {
				log.Print("WARNING: Speech recognition request exceeded limit of 60 seconds.")
			}
			panic(err)
		}

		for _, result := range resp.Results {
			if result.IsFinal {
				for _, alt := range result.Alternatives {
					fmt.Printf("\033[2K\033[G%+v(%v)\n", alt.Transcript, alt.Confidence)
					fmt.Println("============================================================")
				}
			} else {
				fmt.Printf("\033[2K\033[G%+v", result.Alternatives[0].Transcript)
			}
		}
	}

	return nil
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
