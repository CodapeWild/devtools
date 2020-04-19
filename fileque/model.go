package fileque

import (
	"database/sql"
	"fmt"
	"sync"
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
		if _, err = db.Exec(fmt.Sprintf("create index if not exists '%s_%s_index' on '%s'(%s)\n", def_tab_file, v, def_tab_file, v)); err != nil {
			return err
		}
	}

	return nil
}

func addFile(db *sql.DB, m *MFile) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	rslt, err := tx.Exec(fmt.Sprintf("insert into '%s' values(?,?,?,?,?)\n", def_tab_file), m.FId, m.DId, m.IsDir, m.Capacity, m.Path)
	if err != nil {
		return err
	}

	if _, err = rslt.LastInsertId(); err != nil {
		return err
	}

	return tx.Commit()
}

func findFiles(db *sql.DB, where string) ([]*MFile, error) {
	rows, err := db.Query(fmt.Sprintf("select * from '%s' where %s\n", def_tab_file, where))
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
	stmt, err := db.Prepare(fmt.Sprintf("update '%s' set capacity=%d where f_id='%s'\n", def_tab_file, capacity, fid))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()

	return err
}

func deleteFile(db *sql.DB, where string) error {
	_, err := db.Exec(fmt.Sprintf("delete from '%s' where %s\n", def_tab_file, where))

	return err
}

func write(db *sql.DB, mu *sync.Mutex, i, count int) {
	// mu.Lock()
	// defer mu.Unlock()
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("begin. Exec error=%s", err)
		return
	}
	defer tx.Commit()

	result, err := tx.Exec(`INSERT INTO users (user_name) VALUES ("Bobby");`)
	if err != nil {
		fmt.Printf("user insert. Exec error=%s", err)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		fmt.Printf("user writer. LastInsertId error=%s", err)
	}

	fmt.Printf("+")
}
