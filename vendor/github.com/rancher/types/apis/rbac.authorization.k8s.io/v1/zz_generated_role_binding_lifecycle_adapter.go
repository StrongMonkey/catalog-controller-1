package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type RoleBindingLifecycle interface {
	Create(obj *v1.RoleBinding) error
	Remove(obj *v1.RoleBinding) error
	Updated(obj *v1.RoleBinding) error
}

type roleBindingLifecycleAdapter struct {
	lifecycle RoleBindingLifecycle
}

func (w *roleBindingLifecycleAdapter) Create(obj runtime.Object) error {
	return w.lifecycle.Create(obj.(*v1.RoleBinding))
}

func (w *roleBindingLifecycleAdapter) Finalize(obj runtime.Object) error {
	return w.lifecycle.Remove(obj.(*v1.RoleBinding))
}

func (w *roleBindingLifecycleAdapter) Updated(obj runtime.Object) error {
	return w.lifecycle.Updated(obj.(*v1.RoleBinding))
}

func NewRoleBindingLifecycleAdapter(name string, client RoleBindingInterface, l RoleBindingLifecycle) RoleBindingHandlerFunc {
	adapter := &roleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, adapter, client.ObjectClient())
	return func(key string, obj *v1.RoleBinding) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
