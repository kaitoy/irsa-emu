package webhook

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kaitoy/irsa-emu/pkg/logging"
	"github.com/kaitoy/irsa-emu/pkg/util"
	kwlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const sharedCredentialsDir = "/var/run/secrets/irsa-emu.kaitoy.xyz/shared_credentials"

// GetMutatorFunc returns a mutator func that insert a sidecar to the given pod.
func GetMutatorFunc(
	sidecarImageRepo string,
	sidecarImageTag string,
	awsEnvvarsSecret string,
	stsEndpointURL string,
) kwmutating.Mutator {
	return kwmutating.MutatorFunc(func(
		_ context.Context,
		admissionReview *kwmodel.AdmissionReview,
		obj metav1.Object,
	) (*kwmutating.MutatorResult, error) {
		namespace := util.GetNamespace(obj, admissionReview)
		logger := logging.GetLogger().WithValues(
			kwlog.Kv{
				"namespace": namespace,
				"name":      util.GetName(obj),
			},
		)

		pod, ok := obj.(*corev1.Pod)
		if !ok {
			logger.Infof("Skip an object that's not a pod.")
			return &kwmutating.MutatorResult{}, nil
		}
		logger.Infof("Handling a pod.")

		roleARN, err := getAWSRoleARN(pod, namespace)
		if err != nil {
			return nil, fmt.Errorf(
				"get AWS Role ARN for pod %s/%s: %w",
				namespace,
				util.GetName(pod),
				err,
			)
		}
		if roleARN == "" {
			logger.Infof("Skip a pod that doesn't have the annotations 'eks.amazonaws.com/role-arn' in its ServiceAccount.")
			return &kwmutating.MutatorResult{}, nil
		}
		logger.Infof("Role ARN: %s", roleARN)

		volName := addSharedCredentialsDir(pod)
		addEnvVar(pod)
		addSidecar(pod, sidecarImageRepo, sidecarImageTag, volName, roleARN, awsEnvvarsSecret, stsEndpointURL)

		logger.Infof("Mutated a pod.")
		return &kwmutating.MutatorResult{MutatedObject: pod}, nil
	})
}

func getAWSRoleARN(pod *corev1.Pod, namespace string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", fmt.Errorf("get in-cluster config for k8s: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("create k8s client: %w", err)
	}

	sa, err := clientset.CoreV1().ServiceAccounts(namespace).
		Get(context.Background(), pod.Spec.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("get ServiceAccount %s/%s: %w", namespace, pod.Spec.ServiceAccountName, err)
	}

	if sa.Annotations == nil {
		return "", nil
	}
	if arn, ok := sa.Annotations["eks.amazonaws.com/role-arn"]; ok {
		return arn, nil
	}
	return "", nil
}

func addSharedCredentialsDir(pod *corev1.Pod) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 5)
	rand.Read(b)
	volName := fmt.Sprintf("aws-creds-dir-%x", b)

	pod.Spec.Volumes = append(
		pod.Spec.Volumes,
		corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	)

	var ctrs []corev1.Container
	for _, ctr := range pod.Spec.Containers {
		ctr.VolumeMounts = append(
			ctr.VolumeMounts,
			corev1.VolumeMount{
				Name:      volName,
				MountPath: sharedCredentialsDir,
				ReadOnly:  true,
			},
		)
		ctrs = append(ctrs, ctr)
	}
	pod.Spec.Containers = ctrs

	return volName
}

func addEnvVar(pod *corev1.Pod) {
	var ctrs []corev1.Container
	for _, ctr := range pod.Spec.Containers {
		ctr.Env = append(
			ctr.Env,
			corev1.EnvVar{
				Name:  "AWS_SHARED_CREDENTIALS_FILE",
				Value: sharedCredentialsDir + "/credentials",
			},
		)
		ctrs = append(ctrs, ctr)
	}
	pod.Spec.Containers = ctrs
}

func addSidecar(
	pod *corev1.Pod,
	sidecarImageRepo string,
	sidecarImageTag string,
	volName string,
	roleARN string,
	awsEnvvarsSecret string,
	stsEndpointURL string,
) {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 5)
	rand.Read(b)
	ctrName := fmt.Sprintf("irsa-emu-creds-injector-%x", b)

	args := []string{roleARN}
	if len(stsEndpointURL) > 0 {
		args = append(args, stsEndpointURL)
	}

	pod.Spec.Containers = append(
		pod.Spec.Containers,
		corev1.Container{
			Name:  ctrName,
			Image: fmt.Sprintf("%s:%s", sidecarImageRepo, sidecarImageTag),
			Args:  args,
			EnvFrom: []corev1.EnvFromSource{
				{
					SecretRef: &corev1.SecretEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: awsEnvvarsSecret,
						},
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      volName,
					MountPath: "/shared_credentials",
					ReadOnly:  false,
				},
			},
		},
	)
}
