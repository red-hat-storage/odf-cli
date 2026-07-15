package root

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestDefaultClientConfigSupportsImpersonationFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	clientConfig := defaultClientConfig(flags)

	err := flags.Parse([]string{
		"--server", "https://example.com",
		"--insecure-skip-tls-verify",
		"--as", "alice",
		"--as-group", "developers",
		"--as-group", "admins",
		"--as-uid", "1000",
	})
	require.NoError(t, err)

	config, err := clientConfig.ClientConfig()
	require.NoError(t, err)
	require.Equal(t, "alice", config.Impersonate.UserName)
	require.Equal(t, []string{"developers", "admins"}, config.Impersonate.Groups)
	require.Equal(t, "1000", config.Impersonate.UID)
}

func TestDefaultClientConfigUsesODFDefaultNamespace(t *testing.T) {
	originalNamespace := StorageClusterNamespace
	StorageClusterNamespace = "openshift-storage"
	t.Cleanup(func() {
		StorageClusterNamespace = originalNamespace
	})

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	clientConfig := defaultClientConfig(flags)

	namespace, overridden, err := clientConfig.Namespace()
	require.NoError(t, err)
	require.False(t, overridden)
	require.Equal(t, "openshift-storage", namespace)
}

func TestDefaultClientConfigUsesNamespaceFlag(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	clientConfig := defaultClientConfig(flags)

	err := flags.Parse([]string{"--namespace", "custom-storage"})
	require.NoError(t, err)

	namespace, overridden, err := clientConfig.Namespace()
	require.NoError(t, err)
	require.True(t, overridden)
	require.Equal(t, "custom-storage", namespace)
}
