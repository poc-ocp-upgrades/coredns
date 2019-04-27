package object

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type ToFunc func(interface{}) interface{}
type Empty struct{}

func (e *Empty) GetObjectKind() schema.ObjectKind {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return schema.EmptyObjectKind
}
func (e *Empty) GetGenerateName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (e *Empty) SetGenerateName(name string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetUID() types.UID {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (e *Empty) SetUID(uid types.UID) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetGeneration() int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return 0
}
func (e *Empty) SetGeneration(generation int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetSelfLink() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (e *Empty) SetSelfLink(selfLink string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetCreationTimestamp() v1.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1.Time{}
}
func (e *Empty) SetCreationTimestamp(timestamp v1.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetDeletionTimestamp() *v1.Time {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &v1.Time{}
}
func (e *Empty) SetDeletionTimestamp(timestamp *v1.Time) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetDeletionGracePeriodSeconds() *int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetDeletionGracePeriodSeconds(*int64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetLabels() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetLabels(labels map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetAnnotations() map[string]string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetAnnotations(annotations map[string]string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetInitializers() *v1.Initializers {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetInitializers(initializers *v1.Initializers) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetFinalizers() []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetFinalizers(finalizers []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetOwnerReferences() []v1.OwnerReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (e *Empty) SetOwnerReferences([]v1.OwnerReference) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
func (e *Empty) GetClusterName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (e *Empty) SetClusterName(clusterName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
}
