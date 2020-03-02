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

package v1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/storage/v1"
	clientset "k8s.io/client-go/kubernetes/typed/storage/v1"
	listers "k8s.io/client-go/listers/storage/v1"
	"k8s.io/client-go/tools/cache"
)

type StorageClassHandler func(string, *v1.StorageClass) (*v1.StorageClass, error)

type StorageClassController interface {
	generic.ControllerMeta
	StorageClassClient

	OnChange(ctx context.Context, name string, sync StorageClassHandler)
	OnRemove(ctx context.Context, name string, sync StorageClassHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() StorageClassCache
}

type StorageClassClient interface {
	Create(*v1.StorageClass) (*v1.StorageClass, error)
	Update(*v1.StorageClass) (*v1.StorageClass, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v1.StorageClass, error)
	List(opts metav1.ListOptions) (*v1.StorageClassList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.StorageClass, err error)
}

type StorageClassCache interface {
	Get(name string) (*v1.StorageClass, error)
	List(selector labels.Selector) ([]*v1.StorageClass, error)

	AddIndexer(indexName string, indexer StorageClassIndexer)
	GetByIndex(indexName, key string) ([]*v1.StorageClass, error)
}

type StorageClassIndexer func(obj *v1.StorageClass) ([]string, error)

type storageClassController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.StorageClassesGetter
	informer          informers.StorageClassInformer
	gvk               schema.GroupVersionKind
}

func NewStorageClassController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.StorageClassesGetter, informer informers.StorageClassInformer) StorageClassController {
	return &storageClassController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromStorageClassHandlerToHandler(sync StorageClassHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.StorageClass
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.StorageClass))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *storageClassController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.StorageClass))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateStorageClassDeepCopyOnChange(client StorageClassClient, obj *v1.StorageClass, handler func(obj *v1.StorageClass) (*v1.StorageClass, error)) (*v1.StorageClass, error) {
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

func (c *storageClassController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *storageClassController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *storageClassController) OnChange(ctx context.Context, name string, sync StorageClassHandler) {
	c.AddGenericHandler(ctx, name, FromStorageClassHandlerToHandler(sync))
}

func (c *storageClassController) OnRemove(ctx context.Context, name string, sync StorageClassHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromStorageClassHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *storageClassController) Enqueue(name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), "", name)
}

func (c *storageClassController) EnqueueAfter(name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), "", name, duration)
}

func (c *storageClassController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *storageClassController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *storageClassController) Cache() StorageClassCache {
	return &storageClassCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *storageClassController) Create(obj *v1.StorageClass) (*v1.StorageClass, error) {
	return c.clientGetter.StorageClasses().Create(obj)
}

func (c *storageClassController) Update(obj *v1.StorageClass) (*v1.StorageClass, error) {
	return c.clientGetter.StorageClasses().Update(obj)
}

func (c *storageClassController) Delete(name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.StorageClasses().Delete(name, options)
}

func (c *storageClassController) Get(name string, options metav1.GetOptions) (*v1.StorageClass, error) {
	return c.clientGetter.StorageClasses().Get(name, options)
}

func (c *storageClassController) List(opts metav1.ListOptions) (*v1.StorageClassList, error) {
	return c.clientGetter.StorageClasses().List(opts)
}

func (c *storageClassController) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.StorageClasses().Watch(opts)
}

func (c *storageClassController) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.StorageClass, err error) {
	return c.clientGetter.StorageClasses().Patch(name, pt, data, subresources...)
}

type storageClassCache struct {
	lister  listers.StorageClassLister
	indexer cache.Indexer
}

func (c *storageClassCache) Get(name string) (*v1.StorageClass, error) {
	return c.lister.Get(name)
}

func (c *storageClassCache) List(selector labels.Selector) ([]*v1.StorageClass, error) {
	return c.lister.List(selector)
}

func (c *storageClassCache) AddIndexer(indexName string, indexer StorageClassIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.StorageClass))
		},
	}))
}

func (c *storageClassCache) GetByIndex(indexName, key string) (result []*v1.StorageClass, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1.StorageClass, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1.StorageClass))
	}
	return result, nil
}
