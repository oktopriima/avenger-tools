package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: avenger-tools <command>")
		return
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		fmt.Println("GOPATH environment variable not set")
		return
	}

	switch args[0] {
	case "create":
		if err := create(); err != nil {
			fmt.Println(err)
		}

	case "help":
		if err := help(); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Usage: avenger-tools create <command> [<args>]")
	}

	_ = os.RemoveAll("temp")
	return
}

func help() error {
	b, err := os.ReadFile("avenger-tools-help")
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func create() error {
	repoUrl := "git@github.com:oktopriima/marvel.git"
	tempFolder := "temp"

	_, err := git.PlainClone(tempFolder, false, &git.CloneOptions{
		URL:      repoUrl,
		Progress: nil,
	})

	if err != nil {
		log.Println(err)
		return errors.New("failed to clone repository")
	}

	appConfiguration(tempFolder)

	fmt.Println("What is you application name ?")
	reader := bufio.NewReader(os.Stdin)
	inputApp, _, _ := reader.ReadLine()
	appName := string(inputApp)
	if appName == "" {
		return errors.New("no application name provided")
	}

	fmt.Println("What is your git provider ? (github | gitlab | bitbucket)")
	inputGitProvider, _, _ := reader.ReadLine()
	gitProvider := string(inputGitProvider)
	if gitProvider == "" {
		return errors.New("no git user provided")
	}
	gitProvider = gitProvider + ".com"

	fmt.Println("What is your organization ? ")
	inputOrg, _, _ := reader.ReadLine()
	org := string(inputOrg)
	if org == "" {
		return errors.New("no organization provided")
	}

	if err = replacePackageName(tempFolder, fmt.Sprintf("%s/%s/%s", gitProvider, org, appName)); err != nil {
		return err
	}

	targetDir := fmt.Sprintf("%s/src/%s/%s/%s", os.Getenv("GOPATH"), gitProvider, org, appName)
	err = exec.Command("mkdir", "-p", fmt.Sprintf("%s/src/%s/%s", os.Getenv("GOPATH"), gitProvider, org)).Run()
	if err != nil {
		return err
	}

	err = os.Rename(tempFolder, targetDir)
	if err != nil {
		return err
	}

	return nil
}

func appConfiguration(tempFolder string) {
	reader := bufio.NewReader(os.Stdin)
	var driver, migration, seeder, app string

	fmt.Println("------------- DATABASE SECTION ---------------------")
	fmt.Println("1. Which database driver would you like to use? (mysql)")
	fmt.Println("mysql | postgres")

	inputDriver, _, _ := reader.ReadLine()
	driver = string(inputDriver)
	if driver == "" {
		driver = "mysql"
	}

	fmt.Println("-----------------------------------------------------")
	fmt.Println("2. Are you want to include the migration? (yes)")
	fmt.Println("yes | no")
	inputMigration, _, _ := reader.ReadLine()
	migration = string(inputMigration)
	if migration == "" {
		migration = "yes"
	}

	if migration == "no" {
		migrationFolder := fmt.Sprintf("%s/src/cmd/migrate", tempFolder)
		_ = os.RemoveAll(migrationFolder)
	}

	if migration == "yes" {
		if driver == "postgres" {
			// remove mysql migration folder
			mysqlMigrationFolder := fmt.Sprintf("%s/src/cmd/migrate/mysql", tempFolder)
			_ = os.RemoveAll(mysqlMigrationFolder)
		}

		if driver == "mysql" {
			// remove postgres migration folder
			postgresMigrationFolder := fmt.Sprintf("%s/src/cmd/migrate/postgres", tempFolder)
			_ = os.RemoveAll(postgresMigrationFolder)
		}
	}

	fmt.Println("Your project init with")
	fmt.Printf("Database driver: %s\n", driver)
	fmt.Printf("With migration: %s\n", migration)

	if migration == "yes" {
		fmt.Printf("You can check the migration on src/cmd/migrate/%s/\n", driver)
	}

	if driver == "mysql" {
		fmt.Println("-----------------------------------------------------")
		fmt.Println("3. Are you want to include the seeder? (yes)")
		fmt.Println("yes | no")

		inputSeeder, _, _ := reader.ReadLine()
		seeder = string(inputSeeder)
		if seeder == "" {
			seeder = "yes"
		}

		if seeder == "no" {
			// remove seeder folder
			seederFolder := fmt.Sprintf("%s/src/cmd/seeder", tempFolder)
			_ = os.RemoveAll(seederFolder)
		}

		if seeder == "yes" {
			fmt.Printf("You can check the seeder on src/cmd/seeder/\n")
		}
	}

	fmt.Println("------------- APPLICATION SECTION ---------------------")
	fmt.Println("1. What kind of application would you like to build? (api)")
	fmt.Println("api (API with http service)| kafka (KAFKA Listener) | pubsub | (Google PUBSUB Listener)")

	inputApp, _, _ := reader.ReadLine()
	app = string(inputApp)
	if app == "" {
		app = "no"
	}

	kafkaFolder := fmt.Sprintf("%s/src/cmd/kafka", tempFolder)
	kafkaDomainFolder := fmt.Sprintf("%s/src/app/domain/kafka", tempFolder)
	httpFolder := fmt.Sprintf("%s/src/cmd/http", tempFolder)
	httpDomainFolder := fmt.Sprintf("%s/src/app/domain/http", tempFolder)
	pubsubFolder := fmt.Sprintf("%s/src/cmd/pubsub", tempFolder)
	pubsubDomainFolder := fmt.Sprintf("%s/src/app/domain/pubsub", tempFolder)
	if app == "api" {
		// remove kafka
		_ = os.RemoveAll(kafkaFolder)
		_ = os.RemoveAll(kafkaDomainFolder)
		_ = os.RemoveAll(pubsubFolder)
		_ = os.RemoveAll(pubsubDomainFolder)
	}

	if app == "kafka" {
		_ = os.RemoveAll(httpFolder)
		_ = os.RemoveAll(httpDomainFolder)
		_ = os.RemoveAll(pubsubFolder)
		_ = os.RemoveAll(pubsubDomainFolder)
	}

	if app == "pubsub" {
		_ = os.RemoveAll(kafkaFolder)
		_ = os.RemoveAll(kafkaDomainFolder)
		_ = os.RemoveAll(httpFolder)
		_ = os.RemoveAll(httpDomainFolder)
	}

	return
}

func replacePackageName(tempFolder, targetString string) error {
	return filepath.Walk(tempFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".so") || strings.Contains(path, "node_modules") {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Convert to string
		content := string(data)

		oldStr := "github.com/oktopriima/marvel"
		if strings.Contains(content, oldStr) {
			content = strings.ReplaceAll(content, oldStr, targetString)

			// Write back updated file
			err = ioutil.WriteFile(path, []byte(content), info.Mode())
			if err != nil {
				return err
			}
		}

		return nil
	})
}
