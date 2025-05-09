// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	context "context"
	time "time"

	versioned "github.com/refat75/codegen/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/refat75/codegen/pkg/generated/informers/externalversions/internalinterfaces"
	musicsportsheaddevv1 "github.com/refat75/codegen/pkg/generated/listers/music.sportshead.dev/v1"
	apismusicsportsheaddevv1 "github.com/refat75/k8s-crd/pkg/apis/music.sportshead.dev/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SongInformer provides access to a shared informer and lister for
// Songs.
type SongInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() musicsportsheaddevv1.SongLister
}

type songInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSongInformer constructs a new informer for Song type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSongInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSongInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSongInformer constructs a new informer for Song type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSongInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MusicV1().Songs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MusicV1().Songs(namespace).Watch(context.TODO(), options)
			},
		},
		&apismusicsportsheaddevv1.Song{},
		resyncPeriod,
		indexers,
	)
}

func (f *songInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSongInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *songInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apismusicsportsheaddevv1.Song{}, f.defaultInformer)
}

func (f *songInformer) Lister() musicsportsheaddevv1.SongLister {
	return musicsportsheaddevv1.NewSongLister(f.Informer().GetIndexer())
}
