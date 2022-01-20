package pkg

import (
	"github.com/spf13/viper"
	admissionv1 "k8s.io/api/admission/v1"
	"strings"
)

func GetLabels(ar *admissionv1.AdmissionReview, limitCpu, limitMem string) (labels map[string]string, annotations map[string]string, required map[string]interface{}) {
	req := ar.Request
	config := viper.New()
	config.AddConfigPath("/etc/webhook/conf/")
	config.SetConfigName("mutate")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	required = map[string]interface{}{
		"labels": false,
		"annotations": false,
		"ns":[]string{},
	}

	required["ns"] = config.GetStringSlice("mutate.namespaces")
	if config.GetBool("mutate.labels") {
		required["labels"] = true
	}
	if config.GetBool("mutate.annotations"){
		required["annotations"] = true
	}


	allAnnotationKey := config.AllKeys()
	labels = map[string]string{}
	annotations = map[string]string{}
	for _, key := range allAnnotationKey {
		if strings.HasPrefix(key, "annotations") {
			if strings.HasSuffix(key, ".groupname") {
				config.Set(key, req.Name)
			}
			if strings.HasSuffix(key, ".mem") {
				config.Set(key, limitMem)
			}
			if strings.HasSuffix(key, ".cpu") {
				config.Set(key, limitCpu)
			}

			keys := strings.Replace(key, "annotations.", "", -1)
			annotations[keys] = config.GetString(key)
		}
		if strings.HasPrefix(key, "labels") {
			if strings.HasSuffix(key, "groupname") {
				config.Set(key, req.Name)
			}
			keys := strings.Replace(key, "labels.", "", -1)
			labels[keys] = config.GetString(key)
		}
	}
	return labels, annotations, required
}
