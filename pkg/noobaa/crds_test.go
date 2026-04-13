package noobaa

import (
	"context"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	fakeapiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLoadEmbeddedCRD(t *testing.T) {
	for _, path := range embeddedCRDPaths {
		t.Run(path, func(t *testing.T) {
			crd, err := loadEmbeddedCRD(path)
			if err != nil {
				t.Fatalf("loadEmbeddedCRD: %v", err)
			}
			if crd.Name == "" {
				t.Fatal("expected non-empty CRD name")
			}
		})
	}
}

func TestApplyCRD(t *testing.T) {
	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: "objectbuckets.objectbucket.io"},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "objectbucket.io",
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Kind:   "ObjectBucket",
				Plural: "objectbuckets",
			},
			Scope: apiextensionsv1.ClusterScoped,
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{Name: "v1alpha1", Served: true, Storage: true},
			},
		},
	}

	t.Run("create new CRD", func(t *testing.T) {
		client := fakeapiextensions.NewSimpleClientset() //nolint:staticcheck // Skip lint for unit-test
		if err := applyCRD(context.Background(), client, crd.DeepCopy()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(
			context.Background(), crd.Name, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("CRD not found after create: %v", err)
		}
		if got.Name != crd.Name {
			t.Errorf("expected %q, got %q", crd.Name, got.Name)
		}
	})

	t.Run("update existing CRD", func(t *testing.T) {
		client := fakeapiextensions.NewSimpleClientset(crd.DeepCopy()) //nolint:staticcheck // Skip lint for unit-test

		updated := crd.DeepCopy()
		updated.Spec.Names.Kind = "ObjectBucketUpdated"

		if err := applyCRD(context.Background(), client, updated); err != nil {
			t.Fatalf("unexpected error on update: %v", err)
		}

		got, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(
			context.Background(), crd.Name, metav1.GetOptions{})
		if err != nil {
			t.Fatalf("CRD not found after update: %v", err)
		}
		if got.Spec.Names.Kind != "ObjectBucketUpdated" {
			t.Errorf("expected kind %q, got %q", "ObjectBucketUpdated", got.Spec.Names.Kind)
		}
	})
}
