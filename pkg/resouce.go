package pkg

import (
	"github.com/spf13/viper"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	"strings"
)

type InitContainers struct {
	InitContainers  Containers
}
type Containers struct {
	Command []string
	Image string
	ImagePullPolicy string
	Name string
	EnvFrom Configmap
}
type Configmap struct {
	ConfigMapRef
}
type ConfigMapRef struct {
	Name string
}



func NetInitContainers(init *InitContainers) []corev1.Container  {
	if init == nil {
		return nil
	}
	return []corev1.Container{
		corev1.Container{
			Name: init.InitContainers.Name,
			Image: init.InitContainers.Image,
			Command: init.InitContainers.Command,
			EnvFrom: []corev1.EnvFromSource{
				corev1.EnvFromSource{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: init.InitContainers.EnvFrom.Name,
						},
					},
					},
				},
			},
	}
}

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
		"initContainers": false,
	}


	if config.GetBool("mutate.initContainers") {
		var init InitContainers
		if err := config.Unmarshal(&init); err != nil{
			klog.Infof("Unmarshal initContainers failed error: %v", err)
		}else {
			required["initContainers"] = init
		}
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
