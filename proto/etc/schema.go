package etc

import "github.com/gov4git/gov4git/proto"

var (
	EtcNS      = proto.RootNS.Append("etc")
	SettingsNS = EtcNS.Append("settings.json")
)

type Settings struct {
}

var DefaultSettings = Settings{}
