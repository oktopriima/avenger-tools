package internal

import (
	"fmt"
	"os"
)

func (c *createInternal) databaseSetup() error {
	fmt.Println("")
	fmt.Println("Setting up database...")
	var err error
	if c.DBDriver == "mysql" {
		err = c.removePostgresMigration()
		if err != nil {
			return err
		}

		if !c.IncludeMigration {
			err = c.removeMysqlMigration()
			if err != nil {
				return err
			}
		}
	}

	if c.DBDriver == "postgres" {
		err = c.removeMysqlMigration()
		if err != nil {
			return err
		}

		if !c.IncludeMigration {
			err = c.removePostgresMigration()
			if err != nil {
				return err
			}
		}
	}

	if !c.IncludeSeeder {
		err = c.removeSeederFolder()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *createInternal) removePostgresMigration() error {
	migrationFolder := fmt.Sprintf("%s/src/cmd/migrate/postgres", c.TempDir)
	return os.RemoveAll(migrationFolder)
}

func (c *createInternal) removeMysqlMigration() error {
	migrationFolder := fmt.Sprintf("%s/src/cmd/migrate/mysql", c.TempDir)
	return os.RemoveAll(migrationFolder)
}

func (c *createInternal) removeSeederFolder() error {
	seederFolder := fmt.Sprintf("%s/src/cmd/seeder", c.TempDir)
	return os.RemoveAll(seederFolder)
}
