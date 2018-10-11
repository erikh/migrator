package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	. "testing"

	"github.com/jackc/pgx"
	. "gopkg.in/check.v1"
)

var inContainer = os.Getenv("IN_CONTAINER")

type migratorSuite struct{}

var _ = Suite(&migratorSuite{})

var migrateMap = map[string]map[string]struct{}{
	"one": {
		"schema_migrations": {},
		"foo":               {},
		"bar":               {},
		"quux":              {},
		"another":           {},
	},
}

func mkDB(db string) error {
	conn, err := getDB("template1")
	if err != nil {
		return err
	}

	_, err = conn.Exec(fmt.Sprintf("create database %s", db))
	return err
}

func getDB(db string) (*pgx.Conn, error) {
	return pgx.Connect(mkDBParams(db))
}

func mkDBParams(db string) pgx.ConnConfig {
	return pgx.ConnConfig{
		Host:     "localhost",
		Database: db,
		User:     "root",
	}
}

func TestMigrator(t *T) {
	TestingT(t)
}

func clearDBs() error {
	if inContainer == "" {
		panic("stopping tests. You should only run these in the test container; they may drop your databases otherwise!")
	}

	conn, err := getDB("template1")
	if err != nil {
		return err
	}
	defer conn.Close()
	rows, err := conn.Query("select datname from pg_catalog.pg_database")
	if err != nil {
		return err
	}

	var dbname string

	clearDBs := []string{}

	for rows.Next() {
		if err := rows.Scan(&dbname); err != nil {
			return err
		}
		if !(strings.HasPrefix(dbname, "template") || dbname == "postgres") {
			clearDBs = append(clearDBs, dbname)
		}
	}
	rows.Close()

	for _, dbname := range clearDBs {
		_, err := conn.Exec(fmt.Sprintf("drop database %s", dbname))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ms *migratorSuite) SetUpTest(c *C) {
	c.Assert(clearDBs(), IsNil)
}

func (ms *migratorSuite) TearDownTest(c *C) {
	c.Assert(clearDBs(), IsNil)
}

func (ms *migratorSuite) TestBasic(c *C) {
	c.Assert(mkDB("one"), IsNil)

	conn, err := getDB("one")
	c.Assert(err, IsNil)
	defer conn.Close()

	c.Assert(doMigrate(conn, "testdata/one", true), IsNil)
	c.Assert(doMigrate(conn, "testdata/one", true), IsNil)

	rows, err := conn.Query("select tablename from pg_catalog.pg_tables where schemaname='public'")
	c.Assert(err, IsNil)

	var table string

	for rows.Next() {
		c.Assert(rows.Scan(&table), IsNil)
		_, ok := migrateMap["one"][table]
		c.Assert(ok, Equals, true)
	}
	rows.Close()

	num, err := getLatestApplied(conn)
	c.Assert(err, IsNil)
	c.Assert(num, Equals, 3) // N+1 so we know what to pick up next

	dir, err := ioutil.TempDir("", "")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	f, err := os.Create(path.Join(dir, fmt.Sprintf("%d.sql", num)))
	c.Assert(err, IsNil)
	_, err = io.WriteString(f, "create table basic (id int);")
	c.Assert(err, IsNil)

	c.Assert(doMigrate(conn, dir, true), IsNil)

	row := conn.QueryRow("select 1 from pg_catalog.pg_tables where schemaname='public' and tablename='basic';")
	var tmp int
	// this is actually an err check, we don't need the value, but won't get error until scan time
	c.Assert(row.Scan(&tmp), IsNil)
}
