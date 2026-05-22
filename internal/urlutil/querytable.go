package urlutil

import (
	"video_cache/bootstrap"
)

const (
	TQY = "yt_m3u8_qy"
	TYK = "yt_m3u8_yk"
	TQQ = "yt_m3u8_qq"
	TMG = "yt_m3u8_mg"
)

func GetTable(url string) string {
	return url
}
func GetPanTable(table string, env *bootstrap.Env) string {
	switch table {
	case TQQ:
		return env.PanQQ
	case TQY:
		return env.PanQY
	case TYK:
		return env.PanYK
	case TMG:
		return env.PanMG
	default:
		return ""
	}
}
