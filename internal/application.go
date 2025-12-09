package internal

import (
	"fmt"
	"os"
)

func (c *createInternal) applicationSetting() error {
	fmt.Println("")
	fmt.Printf("Setting up application %s service...\n", c.AppService)

	var err error

	// remove git folder
	err = os.RemoveAll(fmt.Sprintf("%s/.git", c.TempDir))
	if err != nil {
		return err
	}

	if c.AppService == "api" {
		err = c.removeKafkaService()
		if err != nil {
			return err
		}

		err = c.removePubsubService()
		if err != nil {
			return err
		}
	}

	if c.AppService == "kafka" {
		err = c.removeApiService()
		if err != nil {
			return err
		}
		err = c.removePubsubService()
		if err != nil {
			return err
		}
	}

	if c.AppService == "pubsub" {
		err = c.removeApiService()
		if err != nil {
			return err
		}
		err = c.removeKafkaService()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *createInternal) removeApiService() error {
	var err error
	httpFolder := fmt.Sprintf("%s/src/cmd/http", c.TempDir)
	httpDomainFolder := fmt.Sprintf("%s/src/app/domain/http", c.TempDir)

	err = os.RemoveAll(httpFolder)
	if err != nil {
		return err
	}
	err = os.RemoveAll(httpDomainFolder)
	if err != nil {
		return err
	}
	return nil
}

func (c *createInternal) removeKafkaService() error {
	var err error
	kafkaFolder := fmt.Sprintf("%s/src/cmd/kafka", c.TempDir)
	kafkaDomainFolder := fmt.Sprintf("%s/src/app/domain/kafka", c.TempDir)

	err = os.RemoveAll(kafkaFolder)
	if err != nil {
		return err
	}
	err = os.RemoveAll(kafkaDomainFolder)
	if err != nil {
		return err
	}
	return nil
}

func (c *createInternal) removePubsubService() error {
	var err error
	pubsubFolder := fmt.Sprintf("%s/src/cmd/pubsub", c.TempDir)
	pubsubDomainFolder := fmt.Sprintf("%s/src/app/domain/pubsub", c.TempDir)

	err = os.RemoveAll(pubsubFolder)
	if err != nil {
		return err
	}
	err = os.RemoveAll(pubsubDomainFolder)
	if err != nil {
		return err
	}
	return nil
}
