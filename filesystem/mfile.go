package filesystem

import (
	"database/sql"
	"devtools/comerr"
	"fmt"
	"os"
	"time"
)

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
	MediaTypeStrings = []string{"unknow", "audio", "image", "video", "text", "binary"}
	MediaTypeConsts  = map[string]MediaType{"unknow": Media_Unknow, "audio": Media_Audio, "image": Media_Image, "video": Media_Video, "text": Media_Text, "binary": Media_Binary}
)

const (
	File_NotFound = iota + 1
	File_Normal
	File_Hidden
	File_Forbidden
)

type MFile struct {
	Code         string      // column: code
	DirCode      string      // column: dir_code
	IsDirectory  bool        // column: is_dir
	Path         string      // column: path
	OriginalName string      // column: origin_name
	Contains     int         // column: contains
	FileMode     os.FileMode // column: mode
	FileSize     int64       // column: size
	Media        MediaType   // column: media
	Span         int64       // column: span
	Created      int64       // column: created
	Updated      int64       // column: updated
	State        int         // column: state
}

func (this *MFile) Name() string {
	return this.Code
}

func (this *MFile) Size() int64 {
	return this.FileSize
}

func (this *MFile) Mode() os.FileMode {
	return this.FileMode
}

func (this *MFile) ModTime() time.Time {
	return time.Unix(this.Updated, 0)
}

func (this *MFile) IsDir() bool {
	return this.IsDirectory
}

func (this *MFile) Sys() interface{} {
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("create table if not exists '%s' (code text primary key, dir_code text, is_dir integer, path text unique, origin_name text, contains integer, mode integer, size integer, media integer, span integer, created integer, updated integer, state integer)\n", def_tab_mfile))
	if err != nil {
		return err
	}

	if _, err = db.Exec(fmt.Sprintf("create index if not exists '%s_path_index' on '%s' (path)\n", def_tab_mfile, def_tab_mfile)); err != nil {
		return err
	}
	if _, err = db.Exec(fmt.Sprintf("create index if not exists '%s_contains_index' on '%s' (contains asc)\n", def_tab_mfile, def_tab_mfile)); err != nil {
		return err
	}
	if _, err = db.Exec(fmt.Sprintf("create index if not exists '%s_created_index' on '%s' (created desc)\n", def_tab_mfile, def_tab_mfile)); err != nil {
		return err
	}
	_, err = db.Exec(fmt.Sprintf("create index if not exists '%s_updated_index' on '%s' (updated desc)\n", def_tab_mfile, def_tab_mfile))

	return err
}

func findMFile(db *sql.DB, where string) (*MFile, error) {
	row := db.QueryRow(fmt.Sprintf("select * from '%s' where %s\n", def_tab_mfile, where))
	m := &MFile{}

	return m, row.Scan(&m.Code, &m.DirCode, &m.IsDirectory, &m.Path, &m.OriginalName, &m.Contains, &m.FileMode, &m.FileSize, &m.Media, &m.Span, &m.Created, &m.Updated, &m.State)
}

func findMFiles(db *sql.DB, where string) ([]os.FileInfo, error) {
	rows, err := db.Query(fmt.Sprintf("select * from '%s' where %s\n", def_tab_mfile, where))
	if err != nil {
		return nil, err
	}

	var ms []os.FileInfo
	for rows.Next() {
		m := &MFile{}
		if err = rows.Scan(&m.Code, &m.DirCode, &m.IsDirectory, &m.Path, &m.OriginalName, &m.Contains, &m.FileMode, &m.FileSize, &m.Media, &m.Span, &m.Created, &m.Updated, &m.State); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func findDirCodes(db *sql.DB, where string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("select code from '%s' where %s\n", def_tab_mfile, where))
	if err != nil {
		return nil, err
	}

	var cs []string
	for rows.Next() {
		s := ""
		if err = rows.Scan(&s); err != nil {
			return nil, err
		}
		cs = append(cs, s)
	}

	return cs, nil
}

func updateContains(db *sql.DB, dir string, v int) error {
	m, err := findMFile(db, "code='"+dir+"'")
	if err != nil {
		return err
	}
	if !m.IsDirectory {
		return comerr.ParamInvalid
	}

	v = m.Contains + v
	if v < 0 {
		v = 0
	}

	stmt, err := db.Prepare(fmt.Sprintf("update '%s' set contains=%d where code='%s'", def_tab_mfile, v, dir))
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Stmt(stmt).Exec(); err != nil {
		return err
	}

	return tx.Commit()
}

func insertMFile(db *sql.DB, mfile *MFile) error {
	stmt, err := db.Prepare(fmt.Sprintf("insert into '%s' values(?,?,?,?,?,?,?,?,?,?,?,?,?)", def_tab_mfile))
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Stmt(stmt).Exec(mfile.Code, mfile.DirCode, mfile.IsDirectory, mfile.Path, mfile.OriginalName, mfile.Contains, mfile.FileMode, mfile.FileSize, mfile.Media, mfile.Span, time.Now().Unix(), 0, mfile.State); err != nil {
		return err
	}
	if !mfile.IsDirectory {
		if err = updateContains(db, mfile.DirCode, 1); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func deleteMFile(db *sql.DB, code string) error {
	m, err := findMFile(db, "code='"+code+"'")
	if err != nil {
		return err
	}

	if m.IsDirectory {
		_, err = db.Exec(fmt.Sprintf("delete from '%s' where code='%s' and dir_code='%s'", def_tab_mfile, m.Code, m.Code))
	} else {
		if _, err = db.Exec(fmt.Sprintf("delete from '%s' where code='%s'", def_tab_mfile, m.Code)); err == nil {
			err = updateContains(db, m.DirCode, -1)
		}
	}

	return err
}
