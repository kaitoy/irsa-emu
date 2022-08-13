package util

import (
	kwmodel "github.com/slok/kubewebhook/v2/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetName gets the name or the generated name of the given k8s object.
func GetName(obj metav1.Object) string {
	name := obj.GetName()
	if len(name) > 0 {
		return name
	} else {
		return obj.GetGenerateName()
	}
}

// GetNamespace gets the namespace of the given k8s object.
func GetNamespace(obj metav1.Object, admissionReview *kwmodel.AdmissionReview) string {
	ns := obj.GetNamespace()
	if len(ns) > 0 {
		return ns
	} else if len(admissionReview.Namespace) > 0 {
		return admissionReview.Namespace
	} else {
		return "default"
	}
}
