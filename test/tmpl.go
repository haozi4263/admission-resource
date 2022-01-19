package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
)

func main()  {
	data, err := ioutil.ReadFile("/Users/zhanghao/code/golang/src/crd/admission-resource/deploy/test/mutate.yaml")
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("template").Parse(string(data))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmpl)
	var objectURl strings.Builder
	if err := tmpl.Execute(&objectURl, string(data)); err != nil{
		fmt.Println(err)
	}
	fmt.Println(objectURl.String())
}