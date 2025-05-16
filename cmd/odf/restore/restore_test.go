package restore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseFullyQualifiedCRD(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectResource string
		expectGroup    string
		expectVersion  string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "valid CRD: cephclusters.ceph.rook.io",
			input:          "cephclusters.ceph.rook.io",
			expectResource: "cephclusters",
			expectGroup:    "ceph.rook.io",
			expectVersion:  "v1",
		},
		{
			name:           "valid CRD: storageclusters.ocs.openshift.io",
			input:          "storageclusters.ocs.openshift.io",
			expectResource: "storageclusters",
			expectGroup:    "ocs.openshift.io",
			expectVersion:  "v1",
		},
		{
			name:           "valid CRD: cephconnections.csi.ceph.io",
			input:          "cephconnections.csi.ceph.io",
			expectResource: "cephconnections",
			expectGroup:    "csi.ceph.io",
			expectVersion:  "v1alpha1",
		},
		{
			name:           "valid CRD: storagesystems.odf.openshift.io",
			input:          "storagesystems.odf.openshift.io",
			expectResource: "storagesystems",
			expectGroup:    "odf.openshift.io",
			expectVersion:  "v1alpha1",
		},
		{
			name:           "invalid format: missing dot",
			input:          "invalidformat",
			expectErr:      true,
			expectedErrMsg: "invalid CRD format",
		},
		{
			name:           "unsupported group",
			input:          "foo.unsupported.group",
			expectErr:      true,
			expectedErrMsg: "unsupported group",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			resource, group, version, err := parseFullyQualifiedCRD(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResource, resource)
				assert.Equal(t, tt.expectGroup, group)
				assert.Equal(t, tt.expectVersion, version)
			}
		})
	}
}
