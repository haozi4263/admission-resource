package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type cmdb struct {
	appName string `json:"appName"`
	clusterName string `json:"clusterName"`
	cpu string `json:"cpu"`
	disk_model_name string `json:"disk_model_name"`
	disk_size string `json:"disk_size"`
	env string `json:"env"`
	mem string `json:"mem"`
	serverType string `json:"serverType"`
	GroupName `json:"groupName"`
	treeAppId int `json:"treeAppId"`
	userID string `json:"userID"`
	userName string `json:"userName"`
	zone string `json:"zone"`
}
type GroupName struct {
	GroupName string
}


func main()  {
	config := viper.New()


	config.AddConfigPath("/Users/zhanghao/code/golang/src/crd/admission-resource/test/")
	config.SetConfigName("mutate")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	config.Set("annotations.baymax.io/register-plugin-extra.cmdb.tags.groupName","api")
	config.Set("labels.baymax.io/groupname", "api")

	allAnnoKey:= config.AllKeys()
	anno := map[string]string{}
	labels := map[string]string{}
	for _, key := range allAnnoKey{
		if strings.HasPrefix(key, "annotations"){
			anno[key] = config.GetString(key)
		}
		if strings.HasPrefix(key, "labels"){
			labels[key] = config.GetString(key)
		}
	}


	//annotaions := map[string]string{}
	for key, val := range anno{
		//fmt.Println(key)
		fmt.Printf("%s:%s\n",key, val)
	//	split := strings.Split(key, ".")
	//	if len(split) > 3 {
	//		k := strings.Join(split[:len(split)-1], ".")
	//		fmt.Println(k)
	//		annotaions[key] = val
	//	}
	//fmt.Println(annotaions)

	}





}
