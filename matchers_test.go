package kmatch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/kralicky/kmatch"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
)

var _ = Describe("Matchers", func() {
	It("should match deployments", func() {
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "foo",
				Namespace:   "bar",
				Labels:      map[string]string{"app": "foo"},
				Annotations: map[string]string{"foo": "app"},
				Finalizers:  []string{"foo", "bar"},
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
			HaveAnnotations("foo", "app"),
			Not(HaveAnnotations("app", "foo")),
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
			HaveFinalizers("foo"),
			Not(HaveFinalizers("fizz")),
			ConsistOfFinalizers("foo", "bar"),
			Not(ConsistOfFinalizers("foo", "bar", "fizz")),
		))
	})
	It("should match statefulsets", func() {
		statefulSet := &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "foo",
				Namespace:   "bar",
				Labels:      map[string]string{"app": "foo"},
				Annotations: map[string]string{"foo": "app"},
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
			HaveAnnotations("foo", "app"),
			Not(HaveAnnotations("app", "foo")),
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
				Name:        "foo",
				Namespace:   "bar",
				Labels:      map[string]string{"app": "foo"},
				Annotations: map[string]string{"foo": "app"},
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
			HaveAnnotations("foo", "app"),
			Not(HaveAnnotations("app", "foo")),
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
				Name:        "foo",
				Namespace:   "bar",
				Labels:      map[string]string{"app": "foo"},
				Annotations: map[string]string{"foo": "app"},
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
			HaveAnnotations("foo", "app"),
			Not(HaveAnnotations("app", "foo")),
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
				Name:        "foo",
				Namespace:   "bar",
				Labels:      map[string]string{"app": "foo"},
				Annotations: map[string]string{"foo": "app"},
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
			HaveAnnotations("foo", "app"),
			Not(HaveAnnotations("app", "foo")),
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

	It("should match meta.RESTMapping's", func() {
		var m *meta.RESTMapping
		Expect(m).NotTo(Exist())
		m = &meta.RESTMapping{
			GroupVersionKind: schema.GroupVersionKind{
				Group:   "helm.cattle.io",
				Version: "v1",
				Kind:    "HelmChart",
			},
		}
		Expect(m).To(Exist())
	})

	When("we match workload rollouts", func() {
		It("should match deployment rollout status", func() {
			By("verifying an empty spec has a successful rollout")
			deploy := &appsv1.Deployment{}
			Expect(deploy).To(HaveSuccessfulRollout())

			By("verifying a deployment with minimal criteria for a successful rollout")
			deployRollout := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "foo",
					Generation: 2,
				},
				Status: appsv1.DeploymentStatus{
					ObservedGeneration: 2,
				},
			}
			Expect(deployRollout).To(HaveSuccessfulRollout())

			By("verifying a deployment eventually has a successful rollout")
			expectedReplicas := int32(3)
			deployReplica := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "foo",
					Generation: 2,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32(int32(1)),
				},
				Status: appsv1.DeploymentStatus{
					ObservedGeneration: 1,
					UpdatedReplicas:    expectedReplicas,
					Replicas:           0,
					AvailableReplicas:  0,
				},
			}
			deployReplica.Status.ObservedGeneration = 2
			Expect(deployReplica).NotTo(HaveSuccessfulRollout())

			Expect(deployReplica).NotTo(HaveSuccessfulRollout())
			deployReplica.Spec.Replicas = pointer.Int32(expectedReplicas)
			Expect(deployReplica).NotTo(
				HaveSuccessfulRollout(),
				"spec.Replicas does not match the desired replicas",
			)
			deployReplica.Status.Replicas = expectedReplicas
			Expect(deployReplica).NotTo(
				HaveSuccessfulRollout(),
				"status.AvailableReplicas does not match the update replicas",
			)
			deployReplica.Status.AvailableReplicas = expectedReplicas
			Expect(deployReplica).To(HaveSuccessfulRollout())

		})

		It("should match stateful set rollout status", func() {
			By("verifying we throw a failure if the stateful set update strategy is not rolling")
			ss := &appsv1.StatefulSet{}
			Expect(InterceptGomegaFailure(func() {
				Expect(ss).To(HaveSuccessfulRollout())
			})).To(MatchError("rollout status is only available for RollingUpdate strategy type"))

			By("verifying a stateful set eventually has a successful rollout")
			ssReplica := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "foo",
					Generation: 2,
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32(3),
					UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
						Type: appsv1.RollingUpdateStatefulSetStrategyType,
					},
				},
				Status: appsv1.StatefulSetStatus{
					CurrentRevision:    "1",
					UpdateRevision:     "2",
					ObservedGeneration: 1,
				},
			}
			Expect(ssReplica).NotTo(HaveSuccessfulRollout())
			ssReplica.Status.ObservedGeneration = 2
			Expect(ssReplica).NotTo(HaveSuccessfulRollout())
			ssReplica.Status.Replicas = 3
			ssReplica.Status.ReadyReplicas = 3
			ssReplica.Status.CurrentReplicas = 3
			Expect(ssReplica).NotTo(HaveSuccessfulRollout())
			ssReplica.Status.CurrentRevision = "2"
			Expect(ssReplica).To(HaveSuccessfulRollout())

			By("verifying a partioned rollout for a statefulset is eventually successful")
			expectedReplicas := int32(3)
			ssPartitioned := &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "foo",
					Generation: 2,
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: pointer.Int32(expectedReplicas),
					UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
						Type: appsv1.RollingUpdateStatefulSetStrategyType,
						RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
							Partition: pointer.Int32(1),
						},
					},
				},
				Status: appsv1.StatefulSetStatus{
					ReadyReplicas:      expectedReplicas,
					CurrentReplicas:    expectedReplicas,
					CurrentRevision:    "1",
					UpdatedReplicas:    0,
					ObservedGeneration: 2,
				},
			}
			Expect(ssPartitioned).NotTo(HaveSuccessfulRollout())
			ssPartitioned.Status.UpdatedReplicas = expectedReplicas - 1
			Expect(ssPartitioned).To(HaveSuccessfulRollout())

		})

		It("should match daemonset rollout status", func() {
			By("verifying we throw a failure if the daemonsest update strategy is not rolling")
			ds := &appsv1.DaemonSet{}
			Expect(InterceptGomegaFailure(func() {
				Expect(ds).To(HaveSuccessfulRollout())
			})).To(MatchError("rollout status is only available for RollingUpdate strategy type"))

			By("verifying a daemonset eventually has a successful rollout")
			dsReplica := &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "foo",
					Generation: 2,
				},
				Spec: appsv1.DaemonSetSpec{
					UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
						Type: appsv1.RollingUpdateDaemonSetStrategyType,
					},
				},
				Status: appsv1.DaemonSetStatus{
					DesiredNumberScheduled: 3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        0,
					ObservedGeneration:     1,
				},
			}
			Expect(dsReplica).NotTo(HaveSuccessfulRollout())
			dsReplica.Status.ObservedGeneration = 2
			Expect(dsReplica).NotTo(HaveSuccessfulRollout())
			dsReplica.Status.UpdatedNumberScheduled = 3
			Expect(dsReplica).NotTo(HaveSuccessfulRollout())
			dsReplica.Status.NumberAvailable = 3
			Expect(dsReplica).To(HaveSuccessfulRollout())
		})

		It("should throw a failure if we try and match to an non workload object", func() {
			pod := &corev1.Pod{}
			Expect(InterceptGomegaFailure(func() {
				Expect(pod).To(HaveSuccessfulRollout())
			})).To(
				MatchError(
					"unsupported object type *v1.Pod in ReplicaCountMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet, *appsv1.DaemonSet)",
				),
			)
		})
	})
})
