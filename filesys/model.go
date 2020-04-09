package filesys

import (
	"database/sql"
	"fmt"
)

const def_tab_file = "tab_file"

type MFile struct {
	FId      string // column: f_id
	DId      string // column: d_id
	IsDir    bool   // column: is_dir
	Capacity int    // column: capacity
	Path     string // column: path
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf("create table if not exists '%s'(f_id text primary key, d_id text, is_dir integer, capacity integer, path text)\n", def_tab_file))
	if err != nil {
		return err
	}
	for _, v := range []string{"d_id", "path", "capacity"} {
		if _, err = db.Exec(fmt.Sprintf("create index if not exists '%s_%s_index' on '%s'(%s)", def_tab_file, v, def_tab_file, v)); err != nil {
			return err
		}
	}

	return nil
}

func addFile(db *sql.DB, m *MFile) error {
	stmt, err := db.Prepare(fmt.Sprintf("insert into '%s' values(?,?,?,?,?)\n", def_tab_file))
	if err != nil {
		return err
	}

	_, err = stmt.Exec(m.FId, m.DId, m.IsDir, m.Capacity, m.Path)

	return err
}

func findFiles(db *sql.DB, where string) ([]*MFile, error) {
	stmt, err := db.Prepare(fmt.Sprintf("select * from '%s' where %s\n", def_tab_file, where))
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var ms []*MFile
	for rows.Next() {
		m := &MFile{}
		if err = rows.Scan(&m.FId, &m.DId, &m.IsDir, &m.Capacity, &m.Path); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	return ms, nil
}

func updateDirCap(db *sql.DB, fid string, capacity int) error {
	_, err := db.Exec(fmt.Sprintf("update '%s' set capacity=%d where f_id='%s'\n", def_tab_file, capacity, fid))

	return err
}

func deleteFile(db *sql.DB, where string) error {
	_, err := db.Exec(fmt.Sprintf("delete from '%s' where %s", def_tab_file, where))

	return err
}
