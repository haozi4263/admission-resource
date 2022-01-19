package pkg

import (
	"github.com/spf13/viper"
	"strings"
)

func GetLabels(appName,limitCpu, limitMem string) (labels map[string]string, annotations map[string]string) {
	config := viper.New()
	config.AddConfigPath("/etc/webhook/conf/")
	config.SetConfigName("mutate")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	config.Set("labels.baymax.io/groupname", appName)

	config.Set("annotations.baymax.io/register-plugin-extra.cmdb.tags.groupName",appName)
	config.Set("annotations.baymax.io/register-plugin-extra.cmdb.cpu",limitCpu)
	config.Set("annotations.baymax.io/register-plugin-extra.cmdb.mem",limitMem)

	allAnnoKey:= config.AllKeys()
	labels = map[string]string{}
	annotations = map[string]string{}
	for _, key := range allAnnoKey{
		if strings.HasPrefix(key, "annotations"){
			annotations[key] = config.GetString(key)
		}
		if strings.HasPrefix(key, "labels"){
			labels[key] = config.GetString(key)
		}
	}
	return labels, annotations
}
