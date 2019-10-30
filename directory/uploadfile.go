package directory

type MediaType byte

const (
	Media_Audio MediaType = iota + 1
	Media_Image
	Media_Video
	Media_Bytes
)

const (
	Audio = "audio"
	Image = "image"
	Video = "video"
	Bytes = "bytes"
)

var (
	MediaTypeToString = map[MediaType]string{
		Media_Audio: Audio,
		Media_Image: Image,
		Media_Video: Video,
		Media_Bytes: Bytes,
	}
	MediaStringToType = map[string]MediaType{
		MediaTypeToString[Media_Audio]: Media_Audio,
		MediaTypeToString[Media_Image]: Media_Image,
		MediaTypeToString[Media_Video]: Media_Video,
		MediaTypeToString[Media_Bytes]: Media_Bytes,
	}
)

/*
	Url:   domain/pattern/resid
	Media: mime media type
	Cover: res data(video, audio resource) cover url
	Span:  duration in seconds
	Size:  size in bytes
*/
type UploadFileData struct {
	Url   string    `json:"url" bson:"url"`
	Media MediaType `json:"media" bson:"media"`
	Cover string    `json:"cover" bson:"cover"`
	Span  int64     `json:"span" bson:"span"`
	Size  int64     `json:"size" bson:"size"`
}
