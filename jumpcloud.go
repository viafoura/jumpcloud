package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/grahamgreen/goutils"
	"github.com/thejumpcloud/jcapi"
	"github.com/urfave/cli"
)

const (
	UrlBase     string = "https://console.jumpcloud.com/api"
	authUrlBase string = "https://auth.jumpcloud.com"
)

type config struct {
	verbose  bool
	force    bool
	systemID string
	jc       jcapi.JCAPI
}

func main() {
	app := cli.NewApp()
	app.Version = "0.1.1"
	app.Name = "ClearCare JumpCloud"
	app.Usage = "jumpcloud"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "/opt/jc/jcagent.conf",
			Usage: "Specify an alternate agentConfig Default: /opt/jc/jcagent.conf",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Be verbose",
		},
		cli.StringFlag{
			Name:   "apikey, k",
			Usage:  "Your jumpcloud api key",
			EnvVar: "JUMPCLOUD_APIKEY",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Don't ask before deleting",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "tag",
			Usage: "Tag operations",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "create a new tag",
					Action: CreateTag,
				},
			},
		},
		{
			Name:  "system",
			Usage: "System operations",
			Subcommands: []cli.Command{
				{
					Name:   "addTag",
					Usage:  "add a new tag to the system",
					Action: AddTagToSystem,
				},
				{
					Name:   "removeTag",
					Usage:  "remove tag from system",
					Action: RemoveTagFromSystem,
				},
				{
					Name: "updateConfig",
					Usage: `Change System Properties. Valid properties are
	    displayName - string
	    allowSshRootLogin - bool
	    allowSshPasswordAuthentication - bool - required for MFA
	    allowMultiFactorAuthentication - bool
	    allowPublicKeyAuthentication - bool`,
					Action: UpdateSystemConfig,
				},
				{
					Name:   "delete",
					Usage:  "Delete system from JumpCLoud",
					Action: DeleteSystem,
				},
				{
					Name:   "getInfo",
					Usage:  "Geat all the systme settings",
					Action: GetSystem,
				},
			},
		},
	}

	app.Run(os.Args)

}

func GetSystem(c *cli.Context) {
	conf := buildConfig(c)
	system, err := conf.jc.GetSystemById(conf.systemID, true)
	goutils.Check(err)
	out, _ := json.MarshalIndent(system, "", " ")
	os.Stdout.Write(out)
}

func AddTagToSystem(c *cli.Context) {
	tagNameToAdd := c.Args().First()
	goutils.NotEmpty(tagNameToAdd)

	conf := buildConfig(c)

	system, err := conf.jc.GetSystemById(conf.systemID, true)
	goutils.Check(err)

	tagToAdd, err := conf.jc.GetTagByName(tagNameToAdd)
	goutils.Check(err)
	if conf.verbose {
		fmt.Printf("Adding tag Name:%s ID:%s\n", tagToAdd.Name, tagToAdd.Id)
	}

	systemTags := system.Tags

	systemTagNames := make([]string, 0)
	for _, tag := range systemTags {
		if tag.Name != tagNameToAdd {
			systemTagNames = append(systemTagNames, tag.Name)
		}
	}
	systemTagNames = append(systemTagNames, tagNameToAdd)
	if conf.verbose {
		fmt.Printf("Proposed Tag List %v\n", systemTagNames)
	}
	system.TagList = systemTagNames
	updatedSystemID, err := conf.jc.UpdateSystem(system)
	goutils.Check(err)

	system, err = conf.jc.GetSystemById(updatedSystemID, true)
	if conf.verbose {
		fmt.Printf("Tags After Add: %v\n", system.Tags)
	}
}

func RemoveTagFromSystem(c *cli.Context) {
	tagNameToRemove := c.Args().First()
	goutils.NotEmpty(tagNameToRemove)

	conf := buildConfig(c)

	if conf.verbose {
		fmt.Printf("Removing tag:%s\n", tagNameToRemove)
	}

	system, err := conf.jc.GetSystemById(conf.systemID, true)
	goutils.Check(err)

	currentTags := system.Tags

	newTagNames := make([]string, 0)
	for _, tag := range currentTags {
		if tag.Name != tagNameToRemove {
			newTagNames = append(newTagNames, tag.Name)
		}
	}
	if conf.verbose {
		fmt.Printf("Proposed Tag List %v\n", newTagNames)
	}

	system.TagList = newTagNames
	updatedSystemID, err := conf.jc.UpdateSystem(system)
	goutils.Check(err)

	system, err = conf.jc.GetSystemById(updatedSystemID, true)
	if conf.verbose {
		fmt.Printf("Tags After Remove: %v\n", system.Tags)
	}

}

func UpdateSystemConfig(c *cli.Context) {
	configToChange := c.Args().First()
	goutils.NotEmpty(configToChange)
	newConfigValue := c.Args().Get(1)
	goutils.NotEmpty(newConfigValue)

	conf := buildConfig(c)

	if conf.verbose {
		fmt.Printf("Updating %s to %s\n", configToChange, newConfigValue)
	}

	system, err := conf.jc.GetSystemById(conf.systemID, false)
	fmt.Println(system.ToString())
	switch {
	case configToChange == "displayName":
		system.DisplayName = newConfigValue
	case configToChange == "allowSshRootLogin":
		configBool, err := strconv.ParseBool(newConfigValue)
		goutils.Check(err)
		system.AllowSshRootLogin = configBool
	case configToChange == "allowSshPasswordAuthentication":
		configBool, err := strconv.ParseBool(newConfigValue)
		goutils.Check(err)
		system.AllowSshPasswordAuthentication = configBool
	case configToChange == "allowMultiFactorAuthentication":
		configBool, err := strconv.ParseBool(newConfigValue)
		goutils.Check(err)
		system.AllowMultiFactorAuthentication = configBool
	case configToChange == "allowPublicKeyAuthentication":
		configBool, err := strconv.ParseBool(newConfigValue)
		goutils.Check(err)
		system.AllowPublicKeyAuth = configBool
	default:
		log.Fatal("Not a valid config parameter")

	}

	updatedSystemID, err := conf.jc.UpdateSystem(system)
	goutils.Check(err)

	system, err = conf.jc.GetSystemById(updatedSystemID, false)
	goutils.Check(err)

	//out, _ := json.MarshalIndent(system, "", " ")
	//goutils.Check(err2)
	//os.Stdout.Write(out)

	if conf.verbose {
		fmt.Printf("%s now updated to %s\n", configToChange, newConfigValue)
	}

}

func DeleteSystem(c *cli.Context) {
	conf := buildConfig(c)
	if !(conf.force) {
		fmt.Println("You sure you want to do this?(Type \"Yes\" to delete):")
		askForConfirmation()
	}
	system, err := conf.jc.GetSystemById(conf.systemID, false)
	goutils.Check(err)

	if conf.verbose {
		fmt.Println("Deleting: " + system.ToString())
	}
	err = conf.jc.DeleteSystem(system)
	goutils.Check(err)

	if conf.verbose {
		fmt.Println("Deleted: " + conf.systemID)
	}
}

func CreateTag(c *cli.Context) {
	tagNameToAdd := c.Args().First()
	goutils.NotEmpty(tagNameToAdd)
	conf := buildConfig(c)

	if conf.verbose {
		fmt.Printf("Creating tag: %s\n", tagNameToAdd)
	}
	tag := jcapi.JCTag{
		Name:        tagNameToAdd,
		Systems:     make([]string, 0),
		SystemUsers: make([]string, 0),
	}
	tagId, err := conf.jc.AddUpdateTag(jcapi.Insert, tag)
	goutils.Check(err)

	if conf.verbose {
		fmt.Printf("New tag id: %s\n", tagId)
	}
}

func buildConfig(c *cli.Context) config {
	cfg := config{}
	cfg.verbose = c.GlobalBool("verbose")
	cfg.force = c.GlobalBool("force")
	config_file := c.GlobalString("config")
	APIKey := c.GlobalString("apikey")
	var dat map[string]interface{}
	conf, err := ioutil.ReadFile(config_file)
	goutils.Check(err)
	if err := json.Unmarshal(conf, &dat); err != nil {
		panic(err)
	}
	cfg.systemID = dat["systemKey"].(string)
	cfg.jc = jcapi.NewJCAPI(APIKey, UrlBase)

	return cfg
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponse := "Yes"
	if response == okayResponse {
		return true
	} else {
		fmt.Println("You must type \"Yes\" in order to proceed:")
		return askForConfirmation()
	}
}
