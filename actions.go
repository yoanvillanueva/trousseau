package trousseau

import (
	"github.com/codegangsta/cli"
	"log"
	"fmt"
	"strings"
	"time"
)

func hasEnoughArgs(args []string, expected int) bool {
	switch expected {
	case -1:
		if len(args) > 0 {
			return true
		} else {
			return false
		}
	default:
		if len(args) == expected {
			return true
		} else {
			return false
		}
	}
}

func CreateAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 1) {
		log.Fatal("Not enough argument supplied to configure command")
	}

	recipients := strings.Split(c.Args()[0], ",")

	meta := Meta{
		CreatedAt:      	time.Now().String(),
		LastModifiedAt: 	time.Now().String(),
		Recipients:     	recipients,
		TrousseauVersion:	TROUSSEAU_VERSION,
	}

	// Create and write empty store file
	CreateStoreFile(gStorePath, &meta)

	fmt.Println("trousseau created")
}

func PushAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 0) {
		log.Fatal("Not enough arguments supplied to push command")
	}

	environment := NewEnvironment()
	err := environment.OverrideWith(map[string]string{
		"S3Bucket": c.String("s3-bucket"),
		"S3Filename": c.String("s3-remote-filename"),
		"SshPrivateKey": c.String("ssh"),
	})
	if err != nil {
		log.Fatal(err)
	}

	switch c.String("remote-storage") {
	case "s3":
		err = uploadUsingS3(environment)
		if err != nil {
			log.Fatal(err)
		}
	case "scp":
		err = uploadUsingScp(environment)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func PullAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 0) {
		log.Fatal("Not enough arguments supplied to push command")
	}

	environment := NewEnvironment()
	err := environment.OverrideWith(map[string]string{
		"S3Bucket": c.String("s3-bucket"),
		"S3Filename": c.String("s3-remote-filename"),
		"SshPrivateKey": c.String("ssh-private-key"),
	})
	if err != nil {
		log.Fatal(err)
	}

	switch c.String("remote-storage") {
	case "s3":
		err = DownloadUsingS3(environment)
		if err != nil {
			log.Fatal(err)
		}
	case "scp":
		err = DownloadUsingScp(environment)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExportAction(c *cli.Context) {
	fmt.Println("trousseau exported")
}

func ImportAction(c *cli.Context) {
	fmt.Println("trousseau imported")
}

func AddRecipientAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 1) {
		log.Fatal("Not enough argument supplied to add-recipient command")
	}

	recipient := c.Args()[0]

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Decrypt()
	if err != nil {
		log.Fatal(err)
	}

	err = store.DataStore.Meta.AddRecipient(recipient)

	err = store.Encrypt()
	if err != nil {
		log.Fatal(err)
	}

	err = store.WriteToFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s added to trousseau recipients", recipient)
}

func RemoveRecipientAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 1) {
		log.Fatal("Not enough argument supplied to remove-recipient command")
	}

	recipient := c.Args()[0]

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Decrypt()
	if err != nil {
		log.Fatal(err)
	}

	err = store.DataStore.Meta.RemoveRecipient(recipient)

	err = store.Encrypt()
	if err != nil {
		log.Fatal(err)
	}

	err = store.WriteToFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s removed from trousseau recipients", recipient)
}

func GetAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 1) {
		log.Fatal("Not enough argument supplied to get command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	value, err := store.Get(c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s key's value: %s\n", c.Args()[0], value)
}

func SetAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 2) {
		log.Fatal("Not enough argument supplied to set command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Set(c.Args()[0], c.Args()[1])
	if err != nil {
		log.Fatal(err)
	}

	err = store.WriteToFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("key-value pair set: %s:%s\n", c.Args()[0], c.Args()[1])
}

func DelAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 1) {
		log.Fatal("Not enough argument supplied to del command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Del(c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	err = store.WriteToFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s key deleted\n", c.Args()[0])
}

func KeysAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 0) {
		log.Fatal("Not enough argument supplied to keys command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	keys, err := store.Keys()
	if err != nil {
		log.Fatal(err)
	} else {
		for _, k := range keys {
			fmt.Println(k)
		}
	}
}

func ShowAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 0) {
		log.Fatal("Not enough argument supplied to show command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	pairs, err := store.Items()
	if err != nil {
		log.Fatal(err)
	} else {
		for _, pair := range pairs {
			fmt.Printf("%s: %s\n", pair.Key, pair.Value)
		}
	}
}

func MetaAction(c *cli.Context) {
	if !hasEnoughArgs(c.Args(), 0) {
		log.Fatal("Not enough argument supplied to show command")
	}

	store, err := NewEncryptedStoreFromFile(gStorePath)
	if err != nil {
		log.Fatal(err)
	}

	pairs, err := store.Meta()
	if err != nil {
		log.Fatal(err)
	}

	for _, pair := range pairs {
		fmt.Printf("%s: %s\n", pair.Key, pair.Value)
	}
}
