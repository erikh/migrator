package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	white = color.New(color.FgHiWhite)
	blue  = color.New(color.FgHiBlue)
	cyan  = color.New(color.FgHiCyan)
)

func doMigrate(conn *pgx.Conn, dir string, quiet bool) (retErr error) {
	defer func() {
		if retErr != nil {
			retErr = errors.Wrap(retErr, "migrator: during migration")
		}
	}()

	if err := createSchemaMigrations(conn); err != nil {
		return err
	}

	min, err := getLatestApplied(conn)
	if err != nil {
		return err
	}

	max, err := getNewestAvailable(dir)
	if err != nil {
		return err
	}

	// off by one in calc here because we're winding backwards
	if min-1 == max {
		return nil
	}

	for i := min; i <= max; i++ {
		if err := apply(conn, dir, i, quiet); err != nil {
			return err
		}
	}

	return updateMigrationTable(conn, max+1)
}

func createSchemaMigrations(conn *pgx.Conn) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("create table if not exists schema_migrations (id int primary key, applied timestamp default now())")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getLatestApplied(conn *pgx.Conn) (int, error) {
	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	row := tx.QueryRow("select max(id) from schema_migrations")
	row.Scan(&id) // XXX deliberately not checking the error
	return id, nil
}

func getNewestAvailable(dir string) (int, error) {
	d, err := os.Open(dir)
	if err != nil {
		return 0, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return 0, err
	}

	var maxID int

	for _, name := range names {
		id, err := strconv.Atoi(strings.TrimRight(name, path.Ext(name)))
		if err != nil {
			return 0, errors.Wrap(err, "the directory must be free of garbage")
		}

		if id > maxID {
			maxID = id
		}
	}

	return maxID, nil
}

func apply(conn *pgx.Conn, dir string, i int, quiet bool) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sqlfile := path.Join(dir, fmt.Sprintf("%d.sql", i))
	content, err := ioutil.ReadFile(sqlfile)
	if err != nil {
		return err
	}

	if !quiet {
		w, _, err := terminal.GetSize(0)
		if err == nil {
			num := fmt.Sprintf("%d.sql", i)
			white.Println(strings.Repeat("-", w))
			cyan.Print(strings.Repeat(" ", w/2-len(num)/2))
			cyan.Print(num)
			cyan.Println(strings.Repeat(" ", w/2-len(num)/2))
			white.Println(strings.Repeat("-", w))
			blue.Println(strings.TrimSpace(string(content)))
			white.Println(strings.Repeat("-", w))
		} else {
			fmt.Println("-----------")
			fmt.Printf("   %d.sql    ", i)
			fmt.Println("-----------")
			fmt.Println(strings.TrimSpace(string(content)))
			fmt.Println("-----------")
		}
	}

	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	return tx.Commit()
}

func updateMigrationTable(conn *pgx.Conn, val int) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into schema_migrations (id) values ($1)", val)
	if err != nil {
		return err
	}

	return tx.Commit()
}
