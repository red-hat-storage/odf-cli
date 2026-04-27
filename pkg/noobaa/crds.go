package noobaa

import (
	"context"
	"embed"
	"fmt"

	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

//go:embed crds/*.yaml
var embeddedCRDs embed.FS

var embeddedCRDPaths = []string{
	"crds/objectbucket.io_objectbuckets_crd.yaml",
	"crds/objectbucket.io_objectbucketclaims_crd.yaml",
}

const (
	obCRDName  = "objectbuckets.objectbucket.io"
	obcCRDName = "objectbucketclaims.objectbucket.io"
)

func InstallCRDs(ctx context.Context, client apiextensionsclient.Interface) {
	for _, path := range embeddedCRDPaths {
		crd, err := loadEmbeddedCRD(path)
		if err != nil {
			logging.Fatal(err)
		}

		if err := applyCRD(ctx, client, crd); err != nil {
			logging.Fatal(fmt.Errorf("failed to install CRD %q: %v", crd.Name, err))
		}
		logging.Info("CRD %q installed successfully", crd.Name)
	}
}

func UninstallCRDs(ctx context.Context, client apiextensionsclient.Interface) {
	for _, name := range []string{obCRDName, obcCRDName} {
		err := client.ApiextensionsV1().CustomResourceDefinitions().Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				logging.Info("CRD %q not found, skipping\n", name)
				continue
			}
			logging.Fatal(fmt.Errorf("failed to delete CRD %q: %v", name, err))
		}
		logging.Info("CRD %q deleted successfully", name)
	}
}

func loadEmbeddedCRD(path string) (*apiextensionsv1.CustomResourceDefinition, error) {
	body, err := embeddedCRDs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read embedded CRD %q: %w", path, err)
	}

	crd := &apiextensionsv1.CustomResourceDefinition{}
	if err := yaml.Unmarshal(body, crd); err != nil {
		return nil, fmt.Errorf("parse embedded CRD %q: %w", path, err)
	}

	return crd, nil
}

func applyCRD(ctx context.Context, client apiextensionsclient.Interface, crd *apiextensionsv1.CustomResourceDefinition) error {
	_, err := client.ApiextensionsV1().CustomResourceDefinitions().Create(ctx, crd, metav1.CreateOptions{})
	if err == nil || !k8serrors.IsAlreadyExists(err) {
		return err
	}

	logging.Info("CRD %q already exists, updating\n", crd.Name)
	existing, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, crd.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	crd.ResourceVersion = existing.ResourceVersion
	_, err = client.ApiextensionsV1().CustomResourceDefinitions().Update(ctx, crd, metav1.UpdateOptions{})
	return err
}
