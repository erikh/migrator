package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"path"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const appVersion = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Usage = "Migrate SQL databases with ordered execution"
	app.UsageText = path.Base(os.Args[0]) + " [options] [dir]"
	app.Version = appVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Usage: "Sets the username for the connection",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "Sets the password for the connection",
		},
		cli.StringFlag{
			Name:  "database, d",
			Usage: "Set the database for the migration (default: name of dir provided)",
		},
		cli.StringFlag{
			Name:  "host, t",
			Usage: "Set the host to connect to (can be a unix socket)",
			Value: "localhost",
		},
		cli.IntFlag{
			Name:  "port, o",
			Usage: "Set the port to connect to",
		},
		cli.BoolFlag{
			Name:  "ssl, s",
			Usage: "Require SSL",
		},
		cli.BoolFlag{
			Name:  "skip-verify",
			Usage: "Skip SSL verification",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Do not output fancy text telling you what ran",
		},
	}

	app.Action = migrate

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error running "+os.Args[0]+": "+err.Error())
		os.Exit(1)
	}
}

func migrate(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		cli.ShowAppHelp(ctx)
		return errors.New("please provide a migration directory")
	}

	db := ctx.String("database")
	if db == "" {
		db = path.Base(args[0])
	}

	var tlsConfig *tls.Config
	if ctx.Bool("ssl") {
		tlsConfig = &tls.Config{
			ServerName:         ctx.String("host"),
			InsecureSkipVerify: ctx.Bool("skip-verify"),
		}
	}

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:      ctx.String("host"),
		Port:      uint16(ctx.Int("port")),
		Database:  db,
		User:      ctx.String("username"),
		Password:  ctx.String("password"),
		TLSConfig: tlsConfig,
	})
	if err != nil {
		return errors.Wrap(err, "migrator")
	}

	return doMigrate(conn, args[0], ctx.Bool("quiet"))
}
