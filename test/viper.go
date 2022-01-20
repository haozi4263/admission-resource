package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type initContainers struct {
	initContainers  Containers
}

type Containers struct {
	Command []string
	EnvFrom configMapRef
	Image string
	ImagePullPolicy string
	Name string
}
type configMapRef struct {
	Name string
}

func isValueInList(slice []string, val string) (bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main()  {
	config := viper.New()
	config.AddConfigPath("/Users/zhanghao/code/golang/src/crd/admission-resource/test/")
	config.SetConfigName("mutate")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	var init initContainers
	//initContainers := config.Get("initContainers")
	err := config.Unmarshal(&init)
	fmt.Println(1, err)
	fmt.Println(init.initContainers)

	config.Set("annotations.baymax.io/register-plugin-extra.cmdb.tags.groupName","api")
	config.Set("labels.baymax.io/groupname", "api")

	allAnnoKey:= config.AllKeys()
	annotations := map[string]string{}
	labels := map[string]string{}
	for _, key := range allAnnoKey{
		if strings.HasPrefix(key, "annotations"){
			keys := strings.Replace(key, "annotations.","",-1)
			if strings.HasSuffix(keys,".groupname"){
				config.Set(keys,"api")
			}
			annotations[keys] = config.GetString(key)
		}
		if strings.HasPrefix(key, "labels"){
			if strings.HasSuffix(key,"groupname"){
				config.Set(key,"xxxx")
			}
			keys := strings.Replace(key, "labels.","",-1)
			labels[keys] = config.GetString(key)
		}
	}



	//fmt.Println(config.GetBool("mutate.labels"))
	//fmt.Println(config.GetBool("mutate.annotations"))


	require := []string{
		"aaa",
		"bbb",
	}
	if isValueInList(require, "aaa"){
		println(11)
	}






}


