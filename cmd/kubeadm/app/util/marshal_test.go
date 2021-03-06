/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/scheme"
	kubeadmapiv1alpha2 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha2"
)

func TestMarshalUnmarshalYaml(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "someName",
			Namespace: "testNamespace",
			Labels: map[string]string{
				"test": "yes",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	bytes, err := MarshalToYaml(pod, corev1.SchemeGroupVersion)
	if err != nil {
		t.Fatalf("unexpected error marshalling: %v", err)
	}

	t.Logf("\n%s", bytes)

	obj2, err := UnmarshalFromYaml(bytes, corev1.SchemeGroupVersion)
	if err != nil {
		t.Fatalf("unexpected error marshalling: %v", err)
	}

	pod2, ok := obj2.(*corev1.Pod)
	if !ok {
		t.Fatal("did not get a Pod")
	}

	if pod2.Name != pod.Name {
		t.Errorf("expected %q, got %q", pod.Name, pod2.Name)
	}

	if pod2.Namespace != pod.Namespace {
		t.Errorf("expected %q, got %q", pod.Namespace, pod2.Namespace)
	}

	if !reflect.DeepEqual(pod2.Labels, pod.Labels) {
		t.Errorf("expected %v, got %v", pod.Labels, pod2.Labels)
	}

	if pod2.Spec.RestartPolicy != pod.Spec.RestartPolicy {
		t.Errorf("expected %q, got %q", pod.Spec.RestartPolicy, pod2.Spec.RestartPolicy)
	}
}

func TestMarshalUnmarshalToYamlForCodecs(t *testing.T) {
	cfg := &kubeadmapiv1alpha2.MasterConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "MasterConfiguration",
			APIVersion: kubeadmapiv1alpha2.SchemeGroupVersion.String(),
		},
		API: kubeadmapiv1alpha2.API{
			AdvertiseAddress: "10.100.0.1",
			BindPort:         4332,
		},
		NodeRegistration: kubeadmapiv1alpha2.NodeRegistrationOptions{
			Name:      "testNode",
			CRISocket: "/var/run/cri.sock",
		},
		Networking: kubeadmapiv1alpha2.Networking{
			ServiceSubnet: "10.100.0.0/24",
			PodSubnet:     "10.100.1.0/24",
		},
	}
	scheme.Scheme.Default(cfg)

	bytes, err := MarshalToYamlForCodecs(cfg, kubeadmapiv1alpha2.SchemeGroupVersion, scheme.Codecs)
	if err != nil {
		t.Fatalf("unexpected error marshalling MasterConfiguration: %v", err)
	}
	t.Logf("\n%s", bytes)

	obj, err := UnmarshalFromYamlForCodecs(bytes, kubeadmapiv1alpha2.SchemeGroupVersion, scheme.Codecs)
	if err != nil {
		t.Fatalf("unexpected error unmarshalling MasterConfiguration: %v", err)
	}

	cfg2, ok := obj.(*kubeadmapiv1alpha2.MasterConfiguration)
	if !ok || cfg2 == nil {
		t.Fatal("did not get MasterConfiguration back")
	}
	if !reflect.DeepEqual(*cfg, *cfg2) {
		t.Errorf("expected %v, got %v", *cfg, *cfg2)
	}
}
