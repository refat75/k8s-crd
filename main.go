package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/util/retry"
	"os"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// client-go
// Dynamic Create-Update-Delete operation for Custom Resource Definition

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	customRes := schema.GroupVersionResource{
		Group:    "music.sportshead.dev",
		Version:  "v1",
		Resource: "songs",
	}

	song := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "music.sportshead.dev/v1",
			"kind":       "Song",
			"metadata": map[string]interface{}{
				"name": "my-favourite-song",
			},
			"spec": map[string]interface{}{
				"title":  "Shape of You",
				"artist": "Ed Sheeran",
				"rating": 4,
				"genres": []interface{}{"Pop", "Dancehill"},
			},
		},
	}

	//Create CRD -- Song Resource
	fmt.Println("Creating custom resource(song)....")
	result, err := client.Resource(customRes).Namespace(apiv1.NamespaceDefault).Create(context.TODO(), song, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created song %q.\n", result.GetName())
	prompt()

	//Update CRD
	fmt.Println("Updating Custom Resource(song)....")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.Resource(customRes).Namespace(apiv1.NamespaceDefault).Get(context.TODO(), "my-favourite-song", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get resource: %v", getErr))
		}

		//update rating to 5
		if err := unstructured.SetNestedField(result.Object, int64(5), "spec", "rating"); err != nil {
			panic(fmt.Errorf("Failed to set rating: %v", err))
		}
		_, updateErr := client.Resource(customRes).Namespace(apiv1.NamespaceDefault).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Failed to get resource: %v", retryErr))
	}
	fmt.Printf("Updated song %q.\n", result.GetName())
	prompt()

	//Delete Deployment
	fmt.Println("Deleting Custom Resource(song)....")
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(customRes).Namespace(apiv1.NamespaceDefault).Delete(context.TODO(), "my-favourite-song", deleteOptions); err != nil {
		panic(err.Error())
	}
	fmt.Println("Custom Resource(song) Deleted")
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
