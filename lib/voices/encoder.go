package voices

import (
	"Kotone-DiVE/lib/config"
	"bytes"
	"github.com/pion/opus/pkg/oggreader"
	"github.com/u2takey/ffmpeg-go"
	"io"
)

func Encode(orgData []byte, voiceInfo VoiceInfo) ([]byte, error) {
	switch voiceInfo.Container {
	//ToDo: Write pure-go
	default:
		// https://datatracker.ietf.org/doc/html/rfc6716#section-2.1.1
		outBuf := bytes.NewBuffer(nil)
		stream := ffmpeg_go.Input("pipe:0").
			Output("pipe:1", ffmpeg_go.KwArgs{"c:v": "libopus", "ac": "2", "ar": "48000", "f": "opus", "frame_duration": "20", "application": "voip", "b:a": "40k"}).
			WithInput(bytes.NewBuffer(orgData)).
			WithOutput(outBuf)
		if config.CurrentConfig.Debug {
			stream.ErrorToStdOut()
		}
		err := stream.Run()
		if err != nil {
			return nil, err
		}
		out, err := io.ReadAll(outBuf)
		if err != nil {
			return nil, err
		}
		return out, nil
	}
}

func SplitToFrame(data []byte) ([][]byte, error) {
	var chunks [][]byte
	reader, _, err := oggreader.NewWith(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	for {
		page, _, err := reader.ParseNextPage()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		for _, v := range page {
			chunks = append(chunks, v)
		}
	}
	return chunks, nil
}
