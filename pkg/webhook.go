package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"
)

var (
	runtimeScheme = runtime.NewScheme()
	codeFactory   = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codeFactory.UniversalDeserializer()
)

type WhSvrParam struct {
	Port     int
	CertFile string
	KeyFile  string
}

type WebhookServer struct {
	Server              *http.Server // http server
	RESOURCE_MULTIPLE []string     // cpu:memory 4:3 limit/request资源比例
}

func (s *WebhookServer) Handler(writer http.ResponseWriter, request *http.Request) {
	var body []byte
	if request.Body != nil {
		if data, err := ioutil.ReadAll(request.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		klog.Error("empty data body")
		http.Error(writer, "empty data body", http.StatusBadRequest)
		return
	}

	// 校验 content-type
	contentType := request.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("Content-Type is %s, but expect application/json", contentType)
		http.Error(writer, "Content-Type invalid, expect application/json", http.StatusBadRequest)
		return
	}

	// 数据序列化（validate、mutate）请求的数据都是 AdmissionReview
	var admissionResponse *admissionv1.AdmissionResponse
	requestedAdmissionReview := admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Errorf("Can't decode body: %v", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}
	} else {
		// 序列化成功，也就是说获取到了请求的 AdmissionReview 的数据
		if request.URL.Path == "/mutate" {
			//admissionResponse = s.mutate(&requestedAdmissionReview)
		} else if request.URL.Path == "/validate" {
			admissionResponse = s.validate(&requestedAdmissionReview)
		}
	}

	// 构造返回的 AdmissionReview 这个结构体
	responseAdmissionReview := admissionv1.AdmissionReview{}
	// admission/v1
	responseAdmissionReview.APIVersion = requestedAdmissionReview.APIVersion
	responseAdmissionReview.Kind = requestedAdmissionReview.Kind
	if admissionResponse != nil {
		responseAdmissionReview.Response = admissionResponse
		if requestedAdmissionReview.Request != nil { // 返回相同的 UID
			responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		}

	}

	klog.Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))
	// send response
	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Errorf("Can't encode response: %v", err)
		http.Error(writer, fmt.Sprintf("Can't encode response: %v", err), http.StatusBadRequest)
		return
	}
	klog.Info("Ready to write response...")

	if _, err := writer.Write(respBytes); err != nil {
		klog.Errorf("Can't write response: %v", err)
		http.Error(writer, fmt.Sprintf("Can't write reponse: %v", err), http.StatusBadRequest)
	}
}

func (s *WebhookServer) validate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	// Deployment、Sts -> request/limit mutate
	req := ar.Request

	var (
		resource *corev1.ResourceRequirements
	)

	klog.Infof("AdmissionReview for Kind=%s, Namespace=%s Name=%s UID=%s",
		req.Kind.Kind, req.Namespace, req.Name, req.UID)

	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			klog.Errorf("Can't not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
		}
		for _, dep := range deployment.Spec.Template.Spec.Containers {
			resource = &dep.Resources
		}
	case "StatefulSet":
		var statefulset appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &statefulset); err != nil {
			klog.Errorf("Can't not unmarshal raw object: %v", err)
			return &admissionv1.AdmissionResponse{
				Allowed: false,
				Result: &metav1.Status{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
		}
		for _, sts := range statefulset.Spec.VolumeClaimTemplates {
			resource = &sts.Spec.Resources
		}

	default:
		return &admissionv1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("resoucres request: %s and limit: %s is in line with expectations"),
			},
		}
	}

	requestCpu := ResourceConvert(resource.Requests.Cpu().String())
	requestMem := ResourceConvert(resource.Requests.Memory().String())
	limitCpu := ResourceConvert(resource.Limits.Cpu().String())
	limitMem := ResourceConvert(resource.Limits.Memory().String())

	cpuMultiple,_ := strconv.Atoi(s.RESOURCE_MULTIPLE[0])
	MemMultiple,_ := strconv.Atoi(s.RESOURCE_MULTIPLE[1])

	if limitCpu/requestCpu > cpuMultiple || limitMem/requestMem > MemMultiple {
		klog.Info("limit/request资源比例大于4倍不符合资源限制要求!")
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code: http.StatusBadRequest,
				Message: fmt.Sprintf("limit_cpu/request_cpu: %v or limit_mem/limit_cpu: %v above 4倍",
					limitCpu/requestCpu, limitMem/requestMem),
			},
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Result: &metav1.Status{
			Code:    http.StatusOK,
			Message: fmt.Sprintf("resoucres limit/request结果小于等于4符合资源限制要求"),
		},
	}
}
