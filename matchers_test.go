package kmatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kralicky/kmatch"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
			HaveMatchingContainer(And(
				HaveName("foo"),
				HavePorts(8080),
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo"),
			)),
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
			)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("EmptyDir"),
			)),
			HaveMatchingPersistentVolume(And(
				HaveName("bar"),
				HaveStorageClass("some-storage-class"),
			)),
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
						Env:   []corev1.EnvVar{{Name: "FOO", Value: "BAR"}},
						VolumeMounts: []corev1.VolumeMount{
							{Name: "foo", MountPath: "/foo"},
							{Name: "bar", MountPath: "/bar"},
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
				HaveEnv("FOO", "BAR"),
				HaveVolumeMounts("foo", "bar"),
			)),
			HaveMatchingVolume(And(
				HaveName("foo"),
				HaveVolumeSource("Secret"),
			)),
			HaveMatchingVolume(And(
				HaveName("bar"),
				HaveVolumeSource("ConfigMap"),
			)),
		))
	})
})
