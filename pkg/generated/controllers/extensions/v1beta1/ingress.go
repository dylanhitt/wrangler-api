/*
Copyright The Kubernetes Authors.

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

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/extensions/v1beta1"
	clientset "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	listers "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

type IngressHandler func(string, *v1beta1.Ingress) (*v1beta1.Ingress, error)

type IngressController interface {
	generic.ControllerMeta
	IngressClient

	OnChange(ctx context.Context, name string, sync IngressHandler)
	OnRemove(ctx context.Context, name string, sync IngressHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() IngressCache
}

type IngressClient interface {
	Create(*v1beta1.Ingress) (*v1beta1.Ingress, error)
	Update(*v1beta1.Ingress) (*v1beta1.Ingress, error)
	UpdateStatus(*v1beta1.Ingress) (*v1beta1.Ingress, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Ingress, error)
	List(namespace string, opts metav1.ListOptions) (*v1beta1.IngressList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Ingress, err error)
}

type IngressCache interface {
	Get(namespace, name string) (*v1beta1.Ingress, error)
	List(namespace string, selector labels.Selector) ([]*v1beta1.Ingress, error)

	AddIndexer(indexName string, indexer IngressIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.Ingress, error)
}

type IngressIndexer func(obj *v1beta1.Ingress) ([]string, error)

type ingressController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.IngressesGetter
	informer          informers.IngressInformer
	gvk               schema.GroupVersionKind
}

func NewIngressController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.IngressesGetter, informer informers.IngressInformer) IngressController {
	return &ingressController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromIngressHandlerToHandler(sync IngressHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1beta1.Ingress
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1beta1.Ingress))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *ingressController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1beta1.Ingress))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateIngressDeepCopyOnChange(client IngressClient, obj *v1beta1.Ingress, handler func(obj *v1beta1.Ingress) (*v1beta1.Ingress, error)) (*v1beta1.Ingress, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *ingressController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *ingressController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *ingressController) OnChange(ctx context.Context, name string, sync IngressHandler) {
	c.AddGenericHandler(ctx, name, FromIngressHandlerToHandler(sync))
}

func (c *ingressController) OnRemove(ctx context.Context, name string, sync IngressHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromIngressHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *ingressController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *ingressController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *ingressController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *ingressController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *ingressController) Cache() IngressCache {
	return &ingressCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *ingressController) Create(obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return c.clientGetter.Ingresses(obj.Namespace).Create(obj)
}

func (c *ingressController) Update(obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return c.clientGetter.Ingresses(obj.Namespace).Update(obj)
}

func (c *ingressController) UpdateStatus(obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	return c.clientGetter.Ingresses(obj.Namespace).UpdateStatus(obj)
}

func (c *ingressController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Ingresses(namespace).Delete(name, options)
}

func (c *ingressController) Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Ingress, error) {
	return c.clientGetter.Ingresses(namespace).Get(name, options)
}

func (c *ingressController) List(namespace string, opts metav1.ListOptions) (*v1beta1.IngressList, error) {
	return c.clientGetter.Ingresses(namespace).List(opts)
}

func (c *ingressController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Ingresses(namespace).Watch(opts)
}

func (c *ingressController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Ingress, err error) {
	return c.clientGetter.Ingresses(namespace).Patch(name, pt, data, subresources...)
}

type ingressCache struct {
	lister  listers.IngressLister
	indexer cache.Indexer
}

func (c *ingressCache) Get(namespace, name string) (*v1beta1.Ingress, error) {
	return c.lister.Ingresses(namespace).Get(name)
}

func (c *ingressCache) List(namespace string, selector labels.Selector) ([]*v1beta1.Ingress, error) {
	return c.lister.Ingresses(namespace).List(selector)
}

func (c *ingressCache) AddIndexer(indexName string, indexer IngressIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1beta1.Ingress))
		},
	}))
}

func (c *ingressCache) GetByIndex(indexName, key string) (result []*v1beta1.Ingress, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1beta1.Ingress))
	}
	return result, nil
}

type IngressStatusHandler func(obj *v1beta1.Ingress, status v1beta1.IngressStatus) (v1beta1.IngressStatus, error)

type IngressGeneratingHandler func(obj *v1beta1.Ingress, status v1beta1.IngressStatus) ([]runtime.Object, v1beta1.IngressStatus, error)

func RegisterIngressStatusHandler(ctx context.Context, controller IngressController, condition condition.Cond, name string, handler IngressStatusHandler) {
	statusHandler := &ingressStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromIngressHandlerToHandler(statusHandler.sync))
}

func RegisterIngressGeneratingHandler(ctx context.Context, controller IngressController, apply apply.Apply,
	condition condition.Cond, name string, handler IngressGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &ingressGeneratingHandler{
		IngressGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterIngressStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type ingressStatusHandler struct {
	client    IngressClient
	condition condition.Cond
	handler   IngressStatusHandler
}

func (a *ingressStatusHandler) sync(key string, obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	obj.Status = newStatus
	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(obj, "", nil)
		} else {
			a.condition.SetError(obj, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, obj.Status) {
		var newErr error
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type ingressGeneratingHandler struct {
	IngressGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *ingressGeneratingHandler) Handle(obj *v1beta1.Ingress, status v1beta1.IngressStatus) (v1beta1.IngressStatus, error) {
	objs, newStatus, err := a.IngressGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	apply := a.apply

	if !a.opts.DynamicLookup {
		apply = apply.WithStrictCaching()
	}

	if !a.opts.AllowCrossNamespace && !a.opts.AllowClusterScoped {
		apply = apply.WithSetOwnerReference(true, false).
			WithDefaultNamespace(obj.GetNamespace()).
			WithListerNamespace(obj.GetNamespace())
	}

	if !a.opts.AllowClusterScoped {
		apply = apply.WithRestrictClusterScoped()
	}

	return newStatus, apply.
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
