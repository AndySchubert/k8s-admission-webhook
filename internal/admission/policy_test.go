package admission

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestValidatePod_AllGood(t *testing.T) {
	pod := validPod("nginx:1.25")

	result := ValidatePod(pod)

	if !result.Allowed {
		t.Fatalf("expected pod to be allowed, got denied: %s", result.Message)
	}
}

func TestValidatePod_LatestTag(t *testing.T) {
	pod := validPod("nginx:latest")

	result := ValidatePod(pod)

	if result.Allowed {
		t.Fatalf("expected pod to be denied for latest tag")
	}
}

func TestValidatePod_NoTag(t *testing.T) {
	pod := validPod("nginx")

	result := ValidatePod(pod)

	if result.Allowed {
		t.Fatalf("expected pod to be denied for missing tag")
	}
}

func TestValidatePod_MissingResources(t *testing.T) {
	pod := validPod("nginx:1.25")
	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{}

	result := ValidatePod(pod)

	if result.Allowed {
		t.Fatalf("expected pod to be denied for missing resources")
	}
}

func validPod(image string) *corev1.Pod {
	return &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "app",
					Image: image,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
		},
	}
}
