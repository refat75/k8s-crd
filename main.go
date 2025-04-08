package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"k8s.io/client-go/util/retry"
	"os"

	songv1 "github.com/refat75/codegen/pkg/apis/music.sportshead.dev/v1"
	clientset "github.com/refat75/codegen/pkg/generated/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	client, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	song := &songv1.Song{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-favourite-song",
		},
		Spec: songv1.SongSpec{
			Title:  "Shape of You",
			Artist: "Ed Sheeran",
			Rating: 4,
			Genres: []string{"Pop", "Dancehill"},
		},
	}

	//Create CRD -- Song Resource
	fmt.Println("Creating custom resource(song)....")
	result, err := client.MusicV1().Songs("default").Create(context.TODO(), song, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created song %q.\n", result.GetName())
	prompt()

	//Update CRD
	fmt.Println("Updating Custom Resource(song)....")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.MusicV1().Songs("default").Get(context.TODO(), "my-favourite-song", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get resource: %v", getErr))
		}

		result.Spec.Rating = 5
		_, updateErr := client.MusicV1().Songs("default").Update(context.TODO(), result, metav1.UpdateOptions{})
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
	if err := client.MusicV1().Songs("default").Delete(context.TODO(), "my-favourite-song", deleteOptions); err != nil {
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
