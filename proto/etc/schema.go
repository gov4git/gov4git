package etc

import "github.com/gov4git/gov4git/v2/proto"

var (
	EtcNS      = proto.RootNS.Append("etc")
	SettingsNS = EtcNS.Append("settings.json")
)

type Settings struct {
}

var DefaultSettings = Settings{}
