package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/clearcare/jcapi"
	"github.com/codegangsta/cli"
	"github.com/grahamgreen/goutils"
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

//TODO add flag for apikey but get default from env
//prompt before delete unless the force flag is set

//var APIKey string = os.Getenv("JUMPCLOUD_APIKEY")

func main() {
	app := cli.NewApp()
	app.Name = "ClearCare Jumpcloud"
	app.Usage = "Work w/ the Clouds of Jump"
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
					Name:   "delete",
					Usage:  "Delete system from JumpCLoud",
					Action: DeleteSystem,
				},
			},
		},
	}

	app.Run(os.Args)

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

	systemTagNames := make([]string, len(systemTags)+1)
	for i, tag := range systemTags {
		systemTagNames[i] = tag.Name
	}
	systemTagNames[len(systemTagNames)-1] = tagNameToAdd
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
		fmt.Printf("Removeing tag:%s\n", tagNameToRemove)
	}

	system, err := conf.jc.GetSystemById(conf.systemID, true)
	goutils.Check(err)

	currentTags := system.Tags
	fmt.Println(len(currentTags))

	newTagNames := make([]string, len(currentTags)-1)
	i := 0
	for _, tag := range currentTags {
		if tag.Name != tagNameToRemove {
			fmt.Println(tag.Name)
			newTagNames[i] = tag.Name
			i++
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

func DeleteSystem(c *cli.Context) {
	conf := buildConfig(c)
	//if

}

func CreateTag(c *cli.Context) {
	fmt.Println("my args: " + c.Args().First())
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
