package kmatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kralicky/kmatch"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var _ = Describe("Matchers", func() {
	It("should match deployments", func() {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				Labels:    map[string]string{"app": "foo"},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32(50),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "foo",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:         "foo",
								Image:        "someimage",
								Ports:        []corev1.ContainerPort{{ContainerPort: 8080}},
								Env:          []corev1.EnvVar{{Name: "FOO", Value: "BAR"}},
								VolumeMounts: []corev1.VolumeMount{{Name: "foo", MountPath: "/foo"}},
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("500m"),
										corev1.ResourceMemory: resource.MustParse("2Gi"),
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "foo",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						NodeSelector: map[string]string{
							"foo": "bar",
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "foo",
								Operator: corev1.TolerationOpExists,
							},
							{
								Key:      "bar",
								Operator: corev1.TolerationOpExists,
							},
						},
					},
				},
			},
		}
		Expect(deployment).To(And(
			HaveName("foo"),
			Not(HaveName("bar")),
			HaveNamespace("bar"),
			Not(HaveNamespace("baz")),
			HaveLabels("app", "foo"),
			Not(HaveLabels("app", "bar")),
			HaveReplicaCount(50),
			Not(HaveReplicaCount(100)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("EmptyDir"),
			)),
			Not(HaveMatchingVolume(HaveName("bar"))),
			HaveNodeSelector("foo", "bar"),
			Not(HaveNodeSelector("fizz", "buzz")),
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo"),
				HaveLimits("cpu"),
				HaveLimits(corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				}),
			)),
			HaveTolerations("foo", corev1.Toleration{
				Key:      "bar",
				Operator: corev1.TolerationOpExists,
			}),
			Not(HaveTolerations("baz")),
		))
	})
	It("should match statefulsets", func() {
		statefulSet := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				Labels:    map[string]string{"app": "foo"},
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas: pointer.Int32(50),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "foo",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "foo",
								Image: "someimage",
								Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
								Env:   []corev1.EnvVar{{Name: "FOO", Value: "BAR"}},
								VolumeMounts: []corev1.VolumeMount{
									{Name: "foo", MountPath: "/foo"},
									{Name: "bar", MountPath: "/bar"},
								},
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("500m"),
										corev1.ResourceMemory: resource.MustParse("2Gi"),
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "foo",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
							{
								Name: "bar",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: "bar",
									},
								},
							},
						},
						NodeSelector: map[string]string{
							"foo": "bar",
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "foo",
								Operator: corev1.TolerationOpExists,
							},
							{
								Key:      "bar",
								Operator: corev1.TolerationOpExists,
							},
						},
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "bar",
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							StorageClassName: pointer.String("some-storage-class"),
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "foo",
								},
							},
						},
					},
				},
			},
		}
		Expect(statefulSet).To(And(
			HaveName("foo"),
			Not(HaveName("bar")),
			HaveNamespace("bar"),
			Not(HaveNamespace("baz")),
			HaveLabels("app", "foo"),
			Not(HaveLabels("app", "bar")),
			HaveReplicaCount(50),
			Not(HaveReplicaCount(100)),
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo", "bar"),
				HaveLimits("cpu"),
				HaveLimits(corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				}),
			)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("EmptyDir"),
			)),
			HaveMatchingPersistentVolume(And(
				HaveName("bar"),
				HaveStorageClass("some-storage-class"),
			)),
			HaveNodeSelector("foo", "bar"),
			Not(HaveNodeSelector("fizz", "buzz")),
			HaveTolerations("foo", corev1.Toleration{
				Key:      "bar",
				Operator: corev1.TolerationOpExists,
			}),
			Not(HaveTolerations("baz")),
		))
	})
	It("should match daemonsets", func() {
		daemonset := &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				Labels:    map[string]string{"app": "foo"},
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "foo",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "foo",
								Image: "someimage",
								Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
								Env:   []corev1.EnvVar{{Name: "FOO", Value: "BAR"}},
								VolumeMounts: []corev1.VolumeMount{
									{Name: "foo", MountPath: "/foo"},
									{Name: "bar", MountPath: "/bar"},
								},
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("500m"),
										corev1.ResourceMemory: resource.MustParse("2Gi"),
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "foo",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: "foo",
									},
								},
							},
							{
								Name: "bar",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "bar",
										},
									},
								},
							},
						},
						NodeSelector: map[string]string{
							"foo": "bar",
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "foo",
								Operator: corev1.TolerationOpExists,
							},
							{
								Key:      "bar",
								Operator: corev1.TolerationOpExists,
							},
						},
					},
				},
			},
		}
		Expect(daemonset).To(And(
			HaveName("foo"),
			HaveNamespace("bar"),
			HaveLabels("app", "foo"),
			Not(HaveLabels("app", "bar")),
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo", "bar"),
				HaveLimits("cpu"),
				HaveLimits(corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				}),
			)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource(corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "foo",
					},
				}),
			)),
			HaveMatchingVolume(And(
				HaveName("bar"),
				HaveVolumeSource(corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "bar",
						},
					},
				}),
			)),
			HaveNodeSelector("foo", "bar"),
			Not(HaveNodeSelector("fizz", "buzz")),
			HaveTolerations("foo", corev1.Toleration{
				Key:      "bar",
				Operator: corev1.TolerationOpExists,
			}),
			Not(HaveTolerations("baz")),
		))
	})
	It("should match jobs", func() {
		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				Labels:    map[string]string{"app": "foo"},
			},
			Spec: batchv1.JobSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "foo",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:         "foo",
								Image:        "someimage",
								Ports:        []corev1.ContainerPort{{ContainerPort: 8080}},
								Env:          []corev1.EnvVar{{Name: "FOO", Value: "BAR"}},
								VolumeMounts: []corev1.VolumeMount{{Name: "foo", MountPath: "/foo"}},
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										corev1.ResourceCPU:    resource.MustParse("500m"),
										corev1.ResourceMemory: resource.MustParse("2Gi"),
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "foo",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							},
						},
						NodeSelector: map[string]string{
							"foo": "bar",
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "foo",
								Operator: corev1.TolerationOpExists,
							},
							{
								Key:      "bar",
								Operator: corev1.TolerationOpExists,
							},
						},
					},
				},
			},
		}
		Expect(job).To(And(
			HaveName("foo"),
			Not(HaveName("bar")),
			HaveNamespace("bar"),
			Not(HaveNamespace("baz")),
			HaveLabels("app", "foo"),
			Not(HaveLabels("app", "bar")),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("EmptyDir"),
			)),
			Not(HaveMatchingVolume(HaveName("bar"))),
			HaveNodeSelector("foo", "bar"),
			Not(HaveNodeSelector("fizz", "buzz")),
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo"),
				HaveLimits("cpu"),
				HaveLimits(corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				}),
			)),
			HaveTolerations("foo", corev1.Toleration{
				Key:      "bar",
				Operator: corev1.TolerationOpExists,
			}),
			Not(HaveTolerations("baz")),
		))
	})
	It("should match pods", func() {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				Labels:    map[string]string{"app": "foo"},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "foo",
						Image: "someimage",
						Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
						Env: []corev1.EnvVar{
							{
								Name:  "FOO",
								Value: "BAR",
							},
							{
								Name: "BAZ",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "BAZ",
										},
										Key: "BAT",
									},
								},
							},
							{
								Name: "LOREM",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "LOREM",
										},
										Key: "IPSUM",
									},
								},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "foo", MountPath: "/foo"},
							{Name: "bar", MountPath: "/bar"},
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("500m"),
								corev1.ResourceMemory: resource.MustParse("2Gi"),
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "foo",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "foo",
							},
						},
					},
					{
						Name: "bar",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "bar",
								},
							},
						},
					},
				},
				NodeSelector: map[string]string{
					"foo": "bar",
				},
				Tolerations: []corev1.Toleration{
					{
						Key:      "foo",
						Operator: corev1.TolerationOpExists,
					},
					{
						Key:      "bar",
						Operator: corev1.TolerationOpExists,
					},
				},
			},
		}
		Expect(pod).To(And(
			HaveName("foo"),
			HaveNamespace("bar"),
			HaveLabels("app", "foo"),
			Not(HaveLabels("app", "bar")),
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv(
					"FOO",
					nil,
					"BAZ",
					corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "BAZ",
						},
						Key: "BAT",
					},
					"LOREM",
					corev1.ConfigMapKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "LOREM",
						},
						Key: "IPSUM",
					},
				),
				HaveVolumeMounts("foo", "bar"),
				HaveLimits("cpu"),
				HaveLimits(corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("500m"),
					corev1.ResourceMemory: resource.MustParse("2Gi"),
				}),
			)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("Secret"),
			)),
			HaveMatchingVolume(And(
				HaveName("bar"),
				HaveVolumeSource("ConfigMap"),
			)),
			HaveNodeSelector("foo", "bar"),
			Not(HaveNodeSelector("fizz", "buzz")),
			HaveTolerations("foo", corev1.Toleration{
				Key:      "bar",
				Operator: corev1.TolerationOpExists,
			}),
			Not(HaveTolerations("baz")),
		))
	})
})
