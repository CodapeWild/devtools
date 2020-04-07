package filesys

const def_tab_file = "tab_file"

type MFile struct {
	FId   string // column: f_id
	DId   string // column: d_id
	IsDir bool   // column: is_dir
	Count int    // column: count
	Path  string // column: path
}
