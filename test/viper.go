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

	ns := config.GetStringSlice("mutate.namespaces")
	fmt.Println(ns)

	require := []string{
		"aaa",
		"bbb",
	}
	if isValueInList(require, "aaa"){
		println(11)
	}

}


func isValueInList(slice []string, val string) (bool) {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}