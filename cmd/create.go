package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/oktopriima/avenger-tools/internal"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project from a template",
	RunE: func(cmd *cobra.Command, args []string) error {
		return create()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func create() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is your application name?")
	appName, _ := reader.ReadString('\n')
	appName = strings.TrimSpace(appName)
	if appName == "" {
		return errors.New("no application name provided")
	}

	fmt.Println("What is your git provider? (github.com | gitlab.com | bitbucket.org)")
	gitProvider, _ := reader.ReadString('\n')
	gitProvider = strings.TrimSpace(gitProvider)
	if gitProvider == "" {
		return errors.New("no git provider provided")
	}

	fmt.Println("What is your organization?")
	org, _ := reader.ReadString('\n')
	org = strings.TrimSpace(org)
	if org == "" {
		return errors.New("no organization provided")
	}

	fmt.Println("Which database driver would you like to use? (mysql)")
	fmt.Println("mysql | postgres")
	dbDriver, _ := reader.ReadString('\n')
	dbDriver = strings.TrimSpace(dbDriver)
	if dbDriver == "" {
		dbDriver = "mysql"
	}

	fmt.Println("Do you want to include the migration? (yes)")
	fmt.Println("yes | no")
	migration, _ := reader.ReadString('\n')
	migration = strings.TrimSpace(migration)
	if migration == "" {
		migration = "yes"
	}

	var includeMigration, includeSeeder bool
	if dbDriver == "mysql" {
		fmt.Println("3. Are you want to include the seeder? (yes)")
		fmt.Println("yes | no")

		seeder, _ := reader.ReadString('\n')
		seeder = strings.TrimSpace(seeder)
		if seeder == "" {
			seeder = "yes"
		}

		includeSeeder = seeder == "yes"
	}

	fmt.Println("What kind of application would you like to build? (default: api)")
	fmt.Println("api (API with http service)| kafka (KAFKA Listener) | pubsub | (Google PUBSUB Listener)")
	app, _ := reader.ReadString('\n')
	app = strings.TrimSpace(app)
	if app == "" {
		app = "api"
	}

	includeMigration = migration == "yes"

	tools := internal.NewCreateInternal(
		appName,
		gitProvider,
		org,
		dbDriver,
		app,
		includeMigration,
		includeSeeder,
	)

	if err := tools.Process(); err != nil {
		return err
	}
	return nil
}
