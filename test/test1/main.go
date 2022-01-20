package main

import (
	"fmt"
	"github.com/spf13/viper"
)


type initContainers struct {
	InitContainers  Containers
}

type Containers struct {
	Command []string
	EnvFrom Configmap
	Image string
	ImagePullPolicy string
	Name string
}

type Configmap struct {
	ConfigMapRef
}
type ConfigMapRef struct {
	Name string
}



func main()  {
	config := viper.New()
	config.AddConfigPath("/Users/zhanghao/code/golang/src/crd/admission-resource/test/test1/")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	var cf initContainers
	if err := config.Unmarshal(&cf); err != nil{
		fmt.Println(err)
	}
	fmt.Println(cf.InitContainers)
}


