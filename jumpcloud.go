package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/clearcare/jcapi"
	"github.com/grahamgreen/goutils"
)

const (
	UrlBase     string = "https://console.jumpcloud.com/api"
	authUrlBase string = "https://auth.jumpcloud.com"
)

var APIKey string = os.Getenv("JUMPCLOUD_APIKEY")

func main() {
	var agentConfig string
	flag.StringVar(&agentConfig, "agentConfig", "/opt/jc/jcagent.conf", "jc agent config file")

	flag.Parse()

	conf, err := ioutil.ReadFile(agentConfig)
	goutils.Check(err)

	var dat map[string]interface{}

	if err := json.Unmarshal(conf, &dat); err != nil {
		panic(err)
	}
	systemKey := dat["systemKey"].(string)
	//get agent config from command line or default location
	//get id from confi
	//get system by id
	//add tag

	//cmd options
	//  system add-tags
	//  user add
	//  user add-tags

	jc := jcapi.NewJCAPI(APIKey, UrlBase)
	system, err := jc.GetSystemById(systemKey, true)
	goutils.Check(err)
	tagNameToAdd := []string{"prod"}
	system.TagList = tagNameToAdd
	updatedSystemID, err := jc.UpdateSystem(system)
	goutils.Check(err)

	system, err = jc.GetSystemById(updatedSystemID, true)
	goutils.Check(err)
	fmt.Println(system)

}
