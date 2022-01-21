package pkg

import (
	"strconv"
	"strings"
)

func ResourceConvert(resource string) int {
	if strings.HasSuffix(resource, "m") { // 1m
		c := strings.Replace(resource, "m","",-1)
		cpu, _ := strconv.Atoi(c)
		return cpu
	}

	if strings.HasSuffix(resource, "Gi"){ // 1Gi
		m := strings.Replace(resource, "Gi","", -1)
		memory, _ := strconv.Atoi(m)
		return memory
	}

	if strings.HasSuffix(resource, "Mi"){ // 1Mi
		m := strings.Replace(resource, "Mi","", -1)
		memory, _ := strconv.Atoi(m)
		return memory
	}

	cpu, _ := strconv.Atoi(resource) // 1c
	return cpu * 1000
}


//func ISValueInList(slice []string, val string) (bool) {
//	for _, item := range slice {
//		if item == val {
//			return true
//		}
//	}
//	return false
//}