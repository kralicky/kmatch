package matchers

import (
	"context"
	"fmt"
	"reflect"

	"emperror.dev/errors"

	"github.com/onsi/gomega"
	gtypes "github.com/onsi/gomega/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrUnsupportedObjectType = errors.New("unsupported object type")
	ErrNotAClientObject      = errors.New("target is not a client.Object")
)

type ClientObject struct {
	Client client.Client
	Object client.Object
}

// Object returns a function that will look up the object with the given name
// and namespace in the cluster using the provided client.
func Object(cli client.Client, obj client.Object) func() (client.Object, error) {
	return func() (client.Object, error) {
		copied := obj.DeepCopyObject().(client.Object)
		err := cli.Get(context.Background(), client.ObjectKeyFromObject(obj), copied)
		return copied, err
	}
}

type NameMatcher struct {
	Name string
}

func (o NameMatcher) Match(target interface{}) (success bool, err error) {
	if obj, ok := target.(client.Object); ok {
		return obj.GetName() == o.Name, nil
	}
	val := reflect.ValueOf(target)
	nameField := val.FieldByName("Name")
	if nameField.IsValid() {
		var name string
		if nameField.Kind() == reflect.Ptr {
			name = nameField.Elem().String()
		} else {
			name = nameField.String()
		}
		return name == o.Name, nil
	}
	objectMetaField := val.FieldByName("ObjectMeta")
	if objectMetaField.IsValid() {
		nameField := objectMetaField.FieldByName("Name")
		if nameField.IsValid() {
			return nameField.String() == o.Name, nil
		}
	}
	return false, fmt.Errorf(
		"%w %T in NameMatcher (allowed types: client.Object)",
		ErrUnsupportedObjectType, target)
}

func (o NameMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have name " + o.Name
}

func (o NameMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have name " + o.Name
}

func HaveName(name string) gtypes.GomegaMatcher {
	return &NameMatcher{Name: name}
}

type NamespaceMatcher struct {
	Namespace string
}

func (o NamespaceMatcher) Match(target interface{}) (success bool, err error) {
	if target, ok := target.(client.Object); ok {
		return target.GetNamespace() == o.Namespace, nil
	}
	val := reflect.ValueOf(target)
	nameField := val.FieldByName("Namespace")
	if nameField.IsValid() {
		var name string
		if nameField.Kind() == reflect.Ptr {
			name = nameField.Elem().String()
		} else {
			name = nameField.String()
		}
		return name == o.Namespace, nil
	}
	objectMetaField := val.FieldByName("ObjectMeta")
	if objectMetaField.IsValid() {
		nameField := objectMetaField.FieldByName("Name")
		if nameField.IsValid() {
			return nameField.String() == o.Namespace, nil
		}
	}
	return false, fmt.Errorf(
		"%w %T in NamespaceMatcher (allowed types: client.Object)",
		ErrUnsupportedObjectType, target)
}

func (o NamespaceMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have namespace " + o.Namespace
}

func (o NamespaceMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have namespace " + o.Namespace
}

func HaveNamespace(namespace string) gtypes.GomegaMatcher {
	return &NamespaceMatcher{Namespace: namespace}
}

type ImageMatcher struct {
	Image string
}

func (o ImageMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.Container:
		return t.Image == o.Image, nil
	default:
		return false, fmt.Errorf(
			"%w %T in ImageMatcher (allowed types: corev1.Container)",
			ErrUnsupportedObjectType, target)
	}
}

func (o ImageMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have image " + o.Image
}

func (o ImageMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have image " + o.Image
}

func HaveImage(image string) gtypes.GomegaMatcher {
	return &ImageMatcher{Image: image}
}

type OwnershipMatcher struct {
	Name string
	UID  types.UID
	GVK  schema.GroupVersionKind
}

func (o OwnershipMatcher) Match(target interface{}) (success bool, err error) {
	if target, ok := target.(client.Object); ok {
		ownerRefs := target.GetOwnerReferences()
		if len(ownerRefs) == 0 {
			return false, nil
		}
		for _, ref := range ownerRefs {
			if ref.UID == o.UID &&
				ref.Kind == o.GVK.Kind &&
				ref.APIVersion == o.GVK.GroupVersion().String() &&
				ref.Name == o.Name {
				return true, nil
			}
		}
	} else {
		return false, ErrNotAClientObject
	}
	return false, nil
}

func (o OwnershipMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to be owned by " + o.Name
}

func (o OwnershipMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to be owned by " + o.Name
}

func HaveOwner(owner client.Object) gtypes.GomegaMatcher {
	return &OwnershipMatcher{
		Name: owner.GetName(),
		UID:  owner.GetUID(),
		GVK:  owner.GetObjectKind().GroupVersionKind(),
	}
}

type LabelMatcher struct {
	labels map[string]string
}

func (o LabelMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case *appsv1.Deployment:
		return reflect.DeepEqual(o.labels, t.Labels) &&
			reflect.DeepEqual(o.labels, t.Spec.Selector.MatchLabels) &&
			reflect.DeepEqual(o.labels, t.Spec.Template.Labels), nil
	case *appsv1.StatefulSet:
		return reflect.DeepEqual(o.labels, t.Labels) &&
			reflect.DeepEqual(o.labels, t.Spec.Selector.MatchLabels) &&
			reflect.DeepEqual(o.labels, t.Spec.Template.Labels), nil
	case *appsv1.DaemonSet:
		return reflect.DeepEqual(o.labels, t.Labels) &&
			reflect.DeepEqual(o.labels, t.Spec.Selector.MatchLabels) &&
			reflect.DeepEqual(o.labels, t.Spec.Template.Labels), nil
	case *corev1.Pod:
		return reflect.DeepEqual(o.labels, t.Labels), nil
	case *corev1.Namespace:
		return reflect.DeepEqual(o.labels, t.Labels), nil
	default:
		return false, fmt.Errorf(
			"%w %T in LabelMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet, *appsv1.DaemonSet, *corev1.Pod, *corev1.Namespace)",
			ErrUnsupportedObjectType, target)
	}
}

func (o LabelMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have labels " + fmt.Sprint(target)
}

func (o LabelMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have labels " + fmt.Sprint(target)
}

func HaveLabels(keysAndValues ...string) gtypes.GomegaMatcher {
	lm := &LabelMatcher{labels: make(map[string]string)}
	for i := 0; i < len(keysAndValues); i += 2 {
		lm.labels[keysAndValues[i]] = keysAndValues[i+1]
	}
	return lm
}

type ReplicaCountMatcher struct {
	ReplicaCount int32
}

func (o ReplicaCountMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case *appsv1.Deployment:
		return pointer.Int32Equal(t.Spec.Replicas, &o.ReplicaCount), nil
	case *appsv1.StatefulSet:
		return pointer.Int32Equal(t.Spec.Replicas, &o.ReplicaCount), nil
	default:
		return false, fmt.Errorf(
			"%w %T in ReplicaCountMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet)",
			ErrUnsupportedObjectType, target)
	}
}

func (o ReplicaCountMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a replica count of" + fmt.Sprint(o.ReplicaCount)
}

func (o ReplicaCountMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have replica count of" + fmt.Sprint(o.ReplicaCount)
}

func HaveReplicaCount(replicaCount int32) gtypes.GomegaMatcher {
	return &ReplicaCountMatcher{ReplicaCount: replicaCount}
}

type VolumeSourceMatcher struct {
	Source string
}

func (o VolumeSourceMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.Volume:
		return !reflect.ValueOf(t.VolumeSource).FieldByName(o.Source).IsNil(), nil
	case corev1.PersistentVolume:
		return !reflect.ValueOf(t.Spec.PersistentVolumeSource).FieldByName(o.Source).IsNil(), nil
	default:
		return false, fmt.Errorf(
			"%w %T in VolumeSourceMatcher (allowed types: corev1.Volume, corev1.PersistentVolume)",
			ErrUnsupportedObjectType, target)
	}
}

func (o VolumeSourceMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a volume with a volume source of " + o.Source
}

func (o VolumeSourceMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a volume with a volume source of " + o.Source
}

type VolumeSourceObjectMatcher struct {
	VolumeSource corev1.VolumeSource
}

func (o VolumeSourceObjectMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.Volume:
		return reflect.DeepEqual(o.VolumeSource, t.VolumeSource), nil
	default:
		return false, fmt.Errorf(
			"%w %T in VolumeSourceObjectMatcher (allowed types: corev1.Volume)",
			ErrUnsupportedObjectType, target)
	}
}

func (o VolumeSourceObjectMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching volume source object"
}

func (o VolumeSourceObjectMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching volume source object"
}

// source is the field name in a VolumeSource or PersistentVolumeSource, or
// a corev1.VolumeSource object to compare against
func HaveVolumeSource(source interface{}) gtypes.GomegaMatcher {
	switch t := source.(type) {
	case string:
		return &VolumeSourceMatcher{Source: t}
	case corev1.VolumeSource:
		return &VolumeSourceObjectMatcher{VolumeSource: t}
	default:
		panic("HaveVolumeSource called with an unsupported type")
	}
}

type StorageClassMatcher struct {
	StorageClass string
}

func (o StorageClassMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.PersistentVolumeClaim:
		return t.Spec.StorageClassName != nil &&
			*t.Spec.StorageClassName == o.StorageClass, nil
	case corev1.PersistentVolume:
		return t.Spec.StorageClassName == o.StorageClass, nil
	default:
		return false, fmt.Errorf(
			"%w %T in StorageClassMatcher (allowed types: corev1.PersistentVolumeClaim, corev1.PersistentVolume)",
			ErrUnsupportedObjectType, target)
	}
}

func (o StorageClassMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a storage class of " + o.StorageClass
}

func (o StorageClassMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a storage class of " + o.StorageClass
}

func HaveStorageClass(storageClass string) gtypes.GomegaMatcher {
	return &StorageClassMatcher{StorageClass: storageClass}
}

type VolumeMatcher struct {
	Matcher gtypes.GomegaMatcher
}

func matchAny(matcher gtypes.GomegaMatcher, slice interface{}) (_ bool, err error) {
	val := reflect.ValueOf(slice)
	for i := 0; i < val.Len(); i++ {
		if s, e := matcher.Match(val.Index(i).Interface()); s {
			return true, nil
		} else if e != nil {
			err = errors.Append(err, errors.WithDetails(e, "index", i))
		}
	}
	return
}

func (o VolumeMatcher) Match(target interface{}) (bool, error) {
	switch t := target.(type) {
	case *appsv1.Deployment:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Volumes)
	case *appsv1.StatefulSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Volumes)
	case *appsv1.DaemonSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Volumes)
	case *corev1.Pod:
		return matchAny(o.Matcher, t.Spec.Volumes)
	default:
		return false, fmt.Errorf(
			"%w %T in VolumeMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet, *appsv1.DaemonSet, *corev1.Pod)",
			ErrUnsupportedObjectType, target)
	}
}

func (o VolumeMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching volume"
}

func (o VolumeMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching volume"
}

func HaveMatchingVolume(matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	return &VolumeMatcher{Matcher: matcher}
}

type PersistentVolumeMatcher struct {
	matcher gtypes.GomegaMatcher
}

func (o PersistentVolumeMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case *appsv1.StatefulSet:
		return matchAny(o.matcher, t.Spec.VolumeClaimTemplates)
	default:
		return false, fmt.Errorf(
			"%w %T in PersistentVolumeMatcher (allowed types: *appsv1.StatefulSet)",
			ErrUnsupportedObjectType, target)
	}
}

func (o PersistentVolumeMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching persistent volume"
}

func (o PersistentVolumeMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching persistent volume"
}

func HaveMatchingPersistentVolume(matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	var pvcName string
	if nameMatcher, ok := matcher.(*NameMatcher); ok {
		pvcName = nameMatcher.Name
	} else {
		// check for And, Or, etc.
		matchers := reflect.ValueOf(matcher).Elem().FieldByName("Matchers")
		if matchers.IsValid() {
			for i := 0; i < matchers.Len(); i++ {
				matcher := matchers.Index(i).Interface()
				if nameMatcher, ok := matcher.(*NameMatcher); ok {
					pvcName = nameMatcher.Name
				}
			}
		}
	}
	if pvcName == "" {
		panic("Can't check matching persistent volumes without a name matcher. Try including HasName() somewhere in your match expression.")
	}

	return gomega.And(
		&PersistentVolumeMatcher{matcher: matcher},
		&VolumeMatcher{
			Matcher: gomega.And(
				HaveName(pvcName),
				gomega.Or(
					HaveVolumeSource(corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
							ReadOnly:  false,
						},
					}),
					HaveVolumeSource(corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
							ReadOnly:  true,
						},
					}),
				),
			),
		},
	)
}

type InitContainerMatcher struct {
	Matcher gtypes.GomegaMatcher
}

func (o InitContainerMatcher) Match(target interface{}) (_ bool, err error) {
	switch t := target.(type) {
	case *appsv1.Deployment:
		return matchAny(o.Matcher, t.Spec.Template.Spec.InitContainers)
	case *appsv1.StatefulSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.InitContainers)
	case *appsv1.DaemonSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.InitContainers)
	case *corev1.Pod:
		return matchAny(o.Matcher, t.Spec.InitContainers)
	default:
		return false, fmt.Errorf(
			"%w %T in InitContainerMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet, *appsv1.DaemonSet, *corev1.Pod)",
			ErrUnsupportedObjectType, target)
	}
}

func (o InitContainerMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching init container"
}

func (o InitContainerMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching init container"
}

func HaveMatchingInitContainer(matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	return &InitContainerMatcher{Matcher: matcher}
}

type ContainerMatcher struct {
	Matcher gtypes.GomegaMatcher
}

func (o ContainerMatcher) Match(target interface{}) (_ bool, err error) {
	switch t := target.(type) {
	case *appsv1.Deployment:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Containers)
	case *appsv1.StatefulSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Containers)
	case *appsv1.DaemonSet:
		return matchAny(o.Matcher, t.Spec.Template.Spec.Containers)
	case *corev1.Pod:
		return matchAny(o.Matcher, t.Spec.Containers)
	default:
		return false, fmt.Errorf(
			"%w %T in ContainerMatcher (allowed types: *appsv1.Deployment, *appsv1.StatefulSet, *appsv1.DaemonSet, *corev1.Pod)",
			ErrUnsupportedObjectType, target)
	}
}

func (o ContainerMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching container"
}

func (o ContainerMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching container"
}

func HaveMatchingContainer(matcher gtypes.GomegaMatcher) gtypes.GomegaMatcher {
	return &ContainerMatcher{Matcher: matcher}
}

type EnvVarMatcher struct {
	KeysAndValues map[string]*string
}

func (o EnvVarMatcher) Match(target interface{}) (success bool, err error) {
	envVars := map[string]string{}
	switch t := target.(type) {
	case corev1.Container:
		for _, env := range t.Env {
			envVars[env.Name] = env.Value
		}
		for k, v := range o.KeysAndValues {
			if envValue, ok := envVars[k]; ok {
				if v != nil && envValue != *v {
					return false, nil
				}
			} else {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf(
			"%w %T in EnvVarMatcher (allowed types: corev1.Container)",
			ErrUnsupportedObjectType, target)
	}
}

func (o EnvVarMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have matching environment variables"
}

func (o EnvVarMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have matching environment variables"
}

// keysAndValues should alternate between key names (string) and values which
// could be any string-convertible type or nil, indicating any value is allowed.
func HaveEnv(keysAndValues ...interface{}) gtypes.GomegaMatcher {
	matcher := &EnvVarMatcher{
		KeysAndValues: make(map[string]*string),
	}
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			panic("key names must be strings")
		}
		if v := keysAndValues[i+1]; v == nil {
			matcher.KeysAndValues[key] = nil
		} else {
			matcher.KeysAndValues[key] = pointer.String(fmt.Sprint(v))
		}
	}
	return matcher
}

type VolumeMountMatcher struct {
	VolumeMountNames []string
}

func (o VolumeMountMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.Container:
		for _, expected := range o.VolumeMountNames {
			for _, actual := range t.VolumeMounts {
				if actual.Name == expected {
					return true, nil
				}
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf(
			"%w %T in VolumeMountMatcher (allowed types: corev1.Container)",
			ErrUnsupportedObjectType, target)
	}
}

func (o VolumeMountMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching volume mount"
}

func (o VolumeMountMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching volume mount"
}

func HaveVolumeMounts(volumeMountNames ...string) gtypes.GomegaMatcher {
	return &VolumeMountMatcher{
		VolumeMountNames: volumeMountNames,
	}
}

type PortMatcher struct {
	Ports []intstr.IntOrString
}

func (o PortMatcher) Match(target interface{}) (success bool, err error) {
	switch t := target.(type) {
	case corev1.Container:
		for _, port := range o.Ports {
			if port.Type == intstr.Int {
				for _, tp := range t.Ports {
					if tp.ContainerPort == port.IntVal {
						return true, nil
					}
				}
			} else {
				for _, tp := range t.Ports {
					if tp.Name == port.StrVal {
						return true, nil
					}
				}
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf(
			"%w %T in PortMatcher (allowed types: corev1.Container)",
			ErrUnsupportedObjectType, target)
	}
}

func (o PortMatcher) FailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " to have a matching port"
}

func (o PortMatcher) NegatedFailureMessage(target interface{}) (message string) {
	return "expected " + target.(client.Object).GetName() + " not to have a matching port"
}

// intsOrStrings can accept either an int port number or string port name.
func HavePorts(intsOrStrings ...interface{}) gtypes.GomegaMatcher {
	pm := &PortMatcher{}
	for _, item := range intsOrStrings {
		switch i := item.(type) {
		case int:
			pm.Ports = append(pm.Ports, intstr.FromInt(i))
		case string:
			pm.Ports = append(pm.Ports, intstr.FromString(i))
		}
	}
	return pm
}
