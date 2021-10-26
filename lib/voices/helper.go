package voices

import (
	"errors"
)

var Voices []string

func init() {
	Voices = []string{Watson, Gtts}
}

func VerifyVoice(source *string, voice *string, voiceerror string) error {
	switch *source {
	case Watson:
		return WatsonVerify(voice)
	case Gtts:
		return GttsVerify(voice)
	default:
		return errors.New(voiceerror)
	}
}
