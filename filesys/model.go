package filesys

const def_tab_file = "tab_file"

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

type MFile struct {
	FId     string    // column: f_id
	DId     string    // column: d_id
	Name    string    // column: name
	Path    string    // column: path
	IsDir   bool      // column: is_dir
	Count   int       // column: count
	Size    int64     // column: size
	Span    int64     // column: span
	Media   MediaType // column: media
	Created int64     // column: created
}
