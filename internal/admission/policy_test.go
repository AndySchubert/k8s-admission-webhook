package admission

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestValidatePod_AllGood(t *testing.T) {
	pod := validPod("nginx:1.25")

	result := ValidatePod(pod, DefaultConfig())

	if !result.Allowed {
		t.Fatalf("expected pod to be allowed, got denied: %s", result.Message)
	}
}

func TestValidatePod_LatestTag(t *testing.T) {
	pod := validPod("nginx:latest")

	result := ValidatePod(pod, DefaultConfig())

	if result.Allowed {
		t.Fatalf("expected pod to be denied for latest tag")
	}
}

func TestValidatePod_NoTag(t *testing.T) {
	pod := validPod("nginx")

	result := ValidatePod(pod, DefaultConfig())

	if result.Allowed {
		t.Fatalf("expected pod to be denied for missing tag")
	}
}

func TestValidatePod_MissingResources(t *testing.T) {
	pod := validPod("nginx:1.25")
	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{}

	result := ValidatePod(pod, DefaultConfig())

	if result.Allowed {
		t.Fatalf("expected pod to be denied for missing resources")
	}
}

func TestValidatePod_LatestTagDisabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Policies.DenyLatestTag = false
	pod := validPod("nginx:latest")

	result := ValidatePod(pod, cfg)

	if !result.Allowed {
		t.Fatalf("expected pod to be allowed when denyLatestTag is disabled, got: %s", result.Message)
	}
}

func TestValidatePod_RequireResourcesDisabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Policies.RequireResources = false
	pod := validPod("nginx:1.25")
	pod.Spec.Containers[0].Resources = corev1.ResourceRequirements{}

	result := ValidatePod(pod, cfg)

	if !result.Allowed {
		t.Fatalf("expected pod to be allowed when requireResources is disabled, got: %s", result.Message)
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
