package main

import (
	"cart_service/config"
	"cart_service/migrate/migrations"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"

	pgStorage "shared/storage/pg"

	"github.com/urfave/cli/v2"
)

type migrationsCLIApp struct {
	ctx context.Context

	// lazy init
	dbOnce sync.Once
	db     *bun.DB
	config *config.Config
}

func newMigrationsCLIApp(ctx context.Context, config *config.Config) *migrationsCLIApp {
	app := &migrationsCLIApp{config: config}
	app.ctx = context.WithValue(ctx, struct{}{}, app)
	return app
}

func main() {
	app := &cli.App{
		Name: "bun",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func (app *migrationsCLIApp) DB() *bun.DB {

	app.dbOnce.Do(func() {
		db := pgStorage.NewClientWithDSN(zerolog.Logger{}, app.config.PgSDN, false)
		app.db = db
	})
	return app.db
}

func (app *migrationsCLIApp) CloseDb() {
	err := app.db.Close()
	if err != nil {
		fmt.Println("CloseDb err:", err)
	}
}

func startCLI(c *cli.Context) (context.Context, *migrationsCLIApp, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	app := newMigrationsCLIApp(c.Context, cfg)
	return app.ctx, app, nil
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(ctx, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					ctx, app, err := startCLI(c)
					if err != nil {
						return err
					}
					defer app.CloseDb()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}

					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}
