package internal

import (
	"fmt"
	"os"
	"time"
)

const SourceRepository = "git@github.com:oktopriima/marvel.git"

type createInternal struct {
	AppName          string
	GitProvider      string
	Organization     string
	TargetRoot       string
	TargetDir        string
	TempDir          string
	DBDriver         string
	AppService       string
	PackageName      string
	IncludeMigration bool
	IncludeSeeder    bool
}

type CreateInternal interface {
	Process() error
}

func NewCreateInternal(
	appName, gitProvider, organization, dbDriver, appService string, includeMigration, includeSeeder bool,
) CreateInternal {
	packageName := fmt.Sprintf("%s/%s/%s", gitProvider, organization, appName)
	targetRoot := fmt.Sprintf("%s/src/%s/%s", os.Getenv("GOPATH"), gitProvider, organization)
	targetDir := fmt.Sprintf("%s/%s", targetRoot, appName)
	return &createInternal{
		AppName:          appName,
		GitProvider:      gitProvider,
		Organization:     organization,
		DBDriver:         dbDriver,
		AppService:       appService,
		TargetRoot:       targetRoot,
		TargetDir:        targetDir,
		TempDir:          fmt.Sprintf("temp-%d", time.Now().UnixNano()),
		PackageName:      packageName,
		IncludeMigration: includeMigration,
		IncludeSeeder:    includeSeeder,
	}
}

func (c *createInternal) Process() error {
	defer func(path string) {
		fmt.Printf("cleaning up temp directory : %s\n", path)
		_ = os.RemoveAll(path)
	}(c.TempDir)

	// clone repository
	if err := c.clone(); err != nil {
		c.writeErrorMessage(err)
		return err
	}

	// setting up database
	if err := c.databaseSetup(); err != nil {
		c.writeErrorMessage(err)
		return err
	}

	// setting up application
	if err := c.applicationSetting(); err != nil {
		c.writeErrorMessage(err)
		return err
	}

	// rename package
	if err := c.renamePackage(); err != nil {
		c.writeErrorMessage(err)
		return err
	}

	if err := c.move(); err != nil {
		c.writeErrorMessage(err)
		return err
	}

	c.writeSuccessMessage()
	return nil
}

func (c *createInternal) writeSuccessMessage() {
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("Successfully setting up application.\n")
	fmt.Printf("Your application store at %s\n", c.TargetDir)
	fmt.Printf("Dont forget to run go mod tidy and go mod vendor \n")
	fmt.Printf("For environment variable, cp example.yaml env.yaml and edit it.\n")
	fmt.Printf("Happy coding!\n")
}

func (c *createInternal) writeErrorMessage(err error) {
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("There was an error creating your application: %s\n", err.Error())
	fmt.Printf("Please try again or contact the developer through the github.\n")
	fmt.Printf("Thanks for your time!\n")
}
