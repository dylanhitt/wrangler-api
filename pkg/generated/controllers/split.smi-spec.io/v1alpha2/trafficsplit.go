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

package v1alpha2

import (
	"context"
	"time"

	v1alpha2 "github.com/deislabs/smi-sdk-go/pkg/apis/split/v1alpha2"
	clientset "github.com/rancher/wrangler-api/pkg/generated/clientset/versioned/typed/split/v1alpha2"
	informers "github.com/rancher/wrangler-api/pkg/generated/informers/externalversions/split/v1alpha2"
	listers "github.com/rancher/wrangler-api/pkg/generated/listers/split/v1alpha2"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type TrafficSplitHandler func(string, *v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error)

type TrafficSplitController interface {
	generic.ControllerMeta
	TrafficSplitClient

	OnChange(ctx context.Context, name string, sync TrafficSplitHandler)
	OnRemove(ctx context.Context, name string, sync TrafficSplitHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() TrafficSplitCache
}

type TrafficSplitClient interface {
	Create(*v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error)
	Update(*v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha2.TrafficSplit, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha2.TrafficSplitList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.TrafficSplit, err error)
}

type TrafficSplitCache interface {
	Get(namespace, name string) (*v1alpha2.TrafficSplit, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha2.TrafficSplit, error)

	AddIndexer(indexName string, indexer TrafficSplitIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha2.TrafficSplit, error)
}

type TrafficSplitIndexer func(obj *v1alpha2.TrafficSplit) ([]string, error)

type trafficSplitController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.TrafficSplitsGetter
	informer          informers.TrafficSplitInformer
	gvk               schema.GroupVersionKind
}

func NewTrafficSplitController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.TrafficSplitsGetter, informer informers.TrafficSplitInformer) TrafficSplitController {
	return &trafficSplitController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromTrafficSplitHandlerToHandler(sync TrafficSplitHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha2.TrafficSplit
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha2.TrafficSplit))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *trafficSplitController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha2.TrafficSplit))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateTrafficSplitDeepCopyOnChange(client TrafficSplitClient, obj *v1alpha2.TrafficSplit, handler func(obj *v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error)) (*v1alpha2.TrafficSplit, error) {
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

func (c *trafficSplitController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *trafficSplitController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *trafficSplitController) OnChange(ctx context.Context, name string, sync TrafficSplitHandler) {
	c.AddGenericHandler(ctx, name, FromTrafficSplitHandlerToHandler(sync))
}

func (c *trafficSplitController) OnRemove(ctx context.Context, name string, sync TrafficSplitHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromTrafficSplitHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *trafficSplitController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *trafficSplitController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *trafficSplitController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *trafficSplitController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *trafficSplitController) Cache() TrafficSplitCache {
	return &trafficSplitCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *trafficSplitController) Create(obj *v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error) {
	return c.clientGetter.TrafficSplits(obj.Namespace).Create(obj)
}

func (c *trafficSplitController) Update(obj *v1alpha2.TrafficSplit) (*v1alpha2.TrafficSplit, error) {
	return c.clientGetter.TrafficSplits(obj.Namespace).Update(obj)
}

func (c *trafficSplitController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.TrafficSplits(namespace).Delete(name, options)
}

func (c *trafficSplitController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha2.TrafficSplit, error) {
	return c.clientGetter.TrafficSplits(namespace).Get(name, options)
}

func (c *trafficSplitController) List(namespace string, opts metav1.ListOptions) (*v1alpha2.TrafficSplitList, error) {
	return c.clientGetter.TrafficSplits(namespace).List(opts)
}

func (c *trafficSplitController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.TrafficSplits(namespace).Watch(opts)
}

func (c *trafficSplitController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.TrafficSplit, err error) {
	return c.clientGetter.TrafficSplits(namespace).Patch(name, pt, data, subresources...)
}

type trafficSplitCache struct {
	lister  listers.TrafficSplitLister
	indexer cache.Indexer
}

func (c *trafficSplitCache) Get(namespace, name string) (*v1alpha2.TrafficSplit, error) {
	return c.lister.TrafficSplits(namespace).Get(name)
}

func (c *trafficSplitCache) List(namespace string, selector labels.Selector) ([]*v1alpha2.TrafficSplit, error) {
	return c.lister.TrafficSplits(namespace).List(selector)
}

func (c *trafficSplitCache) AddIndexer(indexName string, indexer TrafficSplitIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha2.TrafficSplit))
		},
	}))
}

func (c *trafficSplitCache) GetByIndex(indexName, key string) (result []*v1alpha2.TrafficSplit, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha2.TrafficSplit, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha2.TrafficSplit))
	}
	return result, nil
}
