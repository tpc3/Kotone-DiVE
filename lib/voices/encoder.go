package voices

import (
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"bytes"
	"github.com/pion/opus/pkg/oggreader"
	"github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"os/exec"
)

var opusencPath string

func init() {
	path, err := exec.LookPath("opusenc")
	if err != nil {
		log.Print("WARN: Opusenc doesn't exists on PATH.")
	} else {
		opusencPath = path
	}
	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		log.Print("WARN: ffmpeg doesn't exists on PATH.")
	}
}

func Encode(orgData []byte, voiceInfo VoiceInfo) ([]byte, error) {
	switch voiceInfo.Container {
	//ToDo: Write pure-go
	case "ogg":
		if voiceInfo.Format == "opus" {
			return orgData, nil
		}
		fallthrough
	case "wav":
		if opusencPath != "" {
			cmd := exec.Command(opusencPath, "--bitrate", "40", "--speech", "--downmix-stereo", "-", "-")
			cmd.Stdin = bytes.NewReader(orgData)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			err := cmd.Run()
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
		fallthrough
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
