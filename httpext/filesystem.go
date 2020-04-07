package httpext

type MediaType int

const (
	Media_Unknow MediaType = iota
	Media_Audio
	Media_Image
	Media_Video
	Media_Text
	Media_Binary
)

var (
	MediaTypeToString = []string{"unknow", "audio", "image", "video", "text", "binary"}
	StringToMediaType = map[string]MediaType{"unknow": Media_Unknow, "audio": Media_Audio, "image": Media_Image, "video": Media_Video, "text": Media_Text, "binary": Media_Binary}
)
