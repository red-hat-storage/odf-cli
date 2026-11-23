#!/usr/bin/env bash

# Source: https://github.com/rook/kubectl-rook-ceph/blob/master/tests/github-action-helper.sh

set -xeEo pipefail

#############
# VARIABLES #
#############
: "${FUNCTION:=${1}}"


create_extra_disk() {
  sudo apt install -y targetcli-fb open-iscsi
  truncate -s 75G ~/iscsi-disk.img
  sudo targetcli /backstores/fileio create disk1 ~/iscsi-disk.img 75G
  local target_iqn=iqn.2026-02.target.local:disk1
  sudo targetcli /iscsi create ${target_iqn}
  sudo targetcli /iscsi/${target_iqn}/tpg1/luns create /backstores/fileio/disk1
  local init_iqn=iqn.2026-02.initiator.local
  echo "InitiatorName=${init_iqn}" | sudo tee /etc/iscsi/initiatorname.iscsi >/dev/null
  sudo targetcli /iscsi/${target_iqn}/tpg1/acls create ${init_iqn}
  sudo targetcli /iscsi/${target_iqn}/tpg1/acls/${init_iqn} create tpg_lun_or_backstore=lun0 mapped_lun=0
  sudo iscsiadm -m discovery -t sendtargets -p 127.0.0.1
  sudo iscsiadm -m node --login
}

# Source: https://github.com/rook/rook
find_extra_block_dev() {
  # shellcheck disable=SC2005 # redirect doesn't work with sudo, so use echo
  echo "$(sudo lsblk)" >/dev/stderr # print lsblk output to stderr for debugging in case of future errors
  # relevant lsblk --pairs example: (MOUNTPOINT identifies boot partition)(PKNAME is Parent dev ID)
  # NAME="sda15" SIZE="106M" TYPE="part" MOUNTPOINT="/boot/efi" PKNAME="sda"
  # NAME="sdb"   SIZE="75G"  TYPE="disk" MOUNTPOINT=""          PKNAME=""
  # NAME="sdb1"  SIZE="75G"  TYPE="part" MOUNTPOINT="/mnt"      PKNAME="sdb"
  boot_dev="$(sudo lsblk --noheading --list --output MOUNTPOINT,PKNAME | grep boot | awk '{print $2}')"
  echo "  == find_extra_block_dev(): boot_dev='$boot_dev'" >/dev/stderr # debug in case of future errors
  # --nodeps ignores partitions
  extra_dev="$(sudo lsblk --noheading --list --nodeps --output KNAME | egrep -v "($boot_dev|loop|nbd)" | head -1)"
  if [ -z "$extra_dev" ]; then
    create_extra_disk >/dev/stderr
    extra_dev="$(sudo lsblk --noheading --list --nodeps --output KNAME | egrep -v "($boot_dev|loop|nbd)" | head -1)"
  fi
  echo "  == find_extra_block_dev(): extra_dev='$extra_dev'" >/dev/stderr # debug in case of future errors
  echo "$extra_dev"                                                       # output of function
}

: "${BLOCK:=$(find_extra_block_dev)}"

# Default namespace values
DEFAULT_OPERATOR_NS="rook-ceph"
DEFAULT_CLUSTER_NS="rook-ceph"

DEFAULT_TIMEOUT=600

# Tool versions
CRICTL_VERSION="v1.31.1"
MINIKUBE_VERSION="v1.35.0"
CRI_DOCKERD_VERSION="v0.3.21"
CNI_PLUGIN_VERSION="v1.6.0"

# Source: https://github.com/rook/kubectl-rook-ceph/blob/master/tests/github-action-helper.sh
install_minikube_with_none_driver() {
    local kubernetes_version="$1"

    if [[ -z "$kubernetes_version" ]]; then
        echo "Error: Kubernetes version must be specified" >&2
        echo "Usage: install_minikube_with_none_driver <k8s-version>" >&2
        return 1
    fi

    echo "Installing minikube with Kubernetes $kubernetes_version"

    # Install system dependencies
    sudo apt update
    sudo apt install -y conntrack socat

    # Install minikube
    echo "Installing minikube $MINIKUBE_VERSION..."
    if ! curl -LO "https://storage.googleapis.com/minikube/releases/$MINIKUBE_VERSION/minikube-linux-amd64"; then
        echo "Failed to download minikube binary" >&2
        return 1
    fi
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
    rm -f minikube-linux-amd64

    # Install container runtime interface
    local cri_dockerd_ver="${CRI_DOCKERD_VERSION#v}"
    echo "Installing cri-dockerd $CRI_DOCKERD_VERSION..."
    if ! curl -LO "https://github.com/Mirantis/cri-dockerd/releases/download/${CRI_DOCKERD_VERSION}/cri-dockerd_${cri_dockerd_ver}.3-0.ubuntu-focal_amd64.deb"; then
        echo "Failed to download cri-dockerd" >&2
        return 1
    fi
    sudo dpkg -i "cri-dockerd_${cri_dockerd_ver}.3-0.ubuntu-focal_amd64.deb"
    rm -f "cri-dockerd_${cri_dockerd_ver}.3-0.ubuntu-focal_amd64.deb"

    # Install crictl (container runtime CLI)
    echo "Installing crictl $CRICTL_VERSION..."
    if ! curl -LO "https://github.com/kubernetes-sigs/cri-tools/releases/download/$CRICTL_VERSION/crictl-$CRICTL_VERSION-linux-amd64.tar.gz"; then
        echo "Failed to download crictl" >&2
        return 1
    fi
    sudo tar zxf "crictl-$CRICTL_VERSION-linux-amd64.tar.gz" -C /usr/local/bin
    rm -f "crictl-$CRICTL_VERSION-linux-amd64.tar.gz"

    # Configure system settings
    sudo sysctl fs.protected_regular=0

    # Install CNI plugins
    echo "Installing CNI plugins $CNI_PLUGIN_VERSION..."
    local cni_plugin_tar="cni-plugins-linux-amd64-$CNI_PLUGIN_VERSION.tgz"
    local cni_plugin_install_dir="/opt/cni/bin"

    if ! curl -LO "https://github.com/containernetworking/plugins/releases/download/$CNI_PLUGIN_VERSION/$cni_plugin_tar"; then
        echo "Failed to download CNI plugins" >&2
        return 1
    fi
    sudo mkdir -p "$cni_plugin_install_dir"
    sudo tar -xf "$cni_plugin_tar" -C "$cni_plugin_install_dir"
    rm "$cni_plugin_tar"

    # Start minikube cluster
    echo "Starting minikube cluster..."
    export MINIKUBE_HOME=$HOME CHANGE_MINIKUBE_NONE_USER=true KUBECONFIG=$HOME/.kube/config
    if ! minikube start \
        --kubernetes-version="$kubernetes_version" \
        --driver=none \
        --container-runtime=docker \
        --memory 6g \
        --cpus=2 \
        --addons ingress \
        --cni=calico; then
        echo "ERROR: minikube start failed" >&2
        return 1
    fi

    echo "Minikube installation completed successfully!"
}

download_and_modify_yaml() {
    local url="$1"
    local output_file="$2"
    local operator_ns="${3:-$DEFAULT_OPERATOR_NS}"
    local cluster_ns="${4:-$DEFAULT_CLUSTER_NS}"

    echo "Downloading $output_file from $url"

    if ! curl -fL "$url" -o "$output_file"; then
        echo "Failed to download $output_file from $url" >&2
        return 1
    fi

    if [[ "$operator_ns" != "$DEFAULT_OPERATOR_NS" || "$cluster_ns" != "$DEFAULT_CLUSTER_NS" ]]; then
        sed -i "s|rook-ceph # namespace:operator|${operator_ns} # namespace:operator|g" "$output_file"
        sed -i "s|rook-ceph # namespace:cluster|${cluster_ns} # namespace:cluster|g" "$output_file"
        sed -i "s|namespace: rook-ceph|namespace: ${operator_ns}|g" "$output_file"
    fi
}

# Apply YAML with kubectl
apply_yaml() {
    local file="$1"
    kubectl create -f "$file"
}

# Apply YAML from URL directly
apply_yaml_from_url() {
    local url="$1"
    kubectl create -f "$url"
}

# Source: https://github.com/rook/rook
use_local_disk() {
  BLOCK_DATA_PART="$(block_dev)1"
  sudo apt purge snapd -y
  sudo dmsetup version || true
  sudo swapoff --all --verbose
  if mountpoint -q /mnt; then
    sudo umount /mnt
    # search for the device since it keeps changing between sda and sdb
    sudo wipefs --all --force "$BLOCK_DATA_PART"
  else
    # it's the hosted runner!
    sudo sgdisk --zap-all -- "$(block_dev)"
    sudo dd if=/dev/zero of="$(block_dev)" bs=1M count=10 oflag=direct,dsync
    sudo parted -s "$(block_dev)" mklabel gpt
  fi
  sudo lsblk
}

deploy_rook() {
    local operator_ns="${1:-$DEFAULT_OPERATOR_NS}"
    local cluster_ns="${2:-$DEFAULT_CLUSTER_NS}"

    echo "Starting Rook-Ceph deployment"
    echo "Operator namespace: $operator_ns"
    echo "Cluster namespace: $cluster_ns"

    # Create custom namespaces if needed
    if [[ "$operator_ns" != "$DEFAULT_OPERATOR_NS" ]]; then
        echo "Creating operator namespace: $operator_ns"
        kubectl create namespace "$operator_ns" || echo "Namespace $operator_ns already exists"
    fi
    if [[ "$cluster_ns" != "$DEFAULT_CLUSTER_NS" && "$cluster_ns" != "$operator_ns" ]]; then
        echo "Creating cluster namespace: $cluster_ns"
        kubectl create namespace "$cluster_ns" || echo "Namespace $cluster_ns already exists"
    fi

    echo "Deploying Rook common resources..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/common.yaml" "common.yaml" "$operator_ns" "$cluster_ns"
    apply_yaml "common.yaml"

    echo "Deploying Custom Resource Definitions (CRDs)..."
    apply_yaml_from_url "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/crds.yaml"

    echo "Deploying Rook operator..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/operator.yaml" "operator.yaml" "$operator_ns" "$cluster_ns"
    apply_yaml "operator.yaml"

    echo "Deploying CSI operator..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/csi-operator.yaml" "csi-operator.yaml" "$operator_ns" "$cluster_ns"
    apply_yaml "csi-operator.yaml"

    # Wait for operator to be ready before proceeding
    echo "Waiting for Rook operator to become ready..."
    wait_for_operator_pod_to_be_ready_state "$operator_ns"

    echo "Deploying Ceph cluster with device filter for $BLOCK..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/cluster-test.yaml" "cluster-test.yaml" "$operator_ns" "$cluster_ns"
    sed -i "s|#deviceFilter:|deviceFilter: ${BLOCK/\/dev\//}|g" cluster-test.yaml
    apply_yaml "cluster-test.yaml"

    echo "Waiting for Ceph cluster to become ready..."
    wait_for_ceph_cluster_to_be_ready_state "$cluster_ns"

    echo "Deploying Ceph toolbox for management..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/toolbox.yaml" "toolbox.yaml" "$operator_ns" "$cluster_ns"
    apply_yaml "toolbox.yaml"

    # Deploy storage class
    echo "Deploying RBD storage class..."
    download_and_modify_yaml "https://raw.githubusercontent.com/rook/rook/master/deploy/examples/csi/rbd/storageclass-test.yaml" "storageclass-rbd.yaml" "$operator_ns" "$cluster_ns"
    if [[ "$operator_ns" != "$DEFAULT_OPERATOR_NS" ]]; then
        sed -i "s|provisioner: rook-ceph.rbd.csi.ceph.com|provisioner: ${operator_ns}.rbd.csi.ceph.com|g" storageclass-rbd.yaml
    fi
    apply_yaml "storageclass-rbd.yaml"

    echo "Rook-Ceph deployment completed successfully!"
}

# CI deploys Rook only. Install a dummy StorageCluster API so
# `odf get health` can list CRs.
install_dummy_storagecluster() {
    local cluster_ns="${1:-openshift-storage}"

    echo "Installing dummy StorageCluster CRD and CR in namespace $cluster_ns"

    kubectl apply -f - <<'EOF'
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: storageclusters.ocs.openshift.io
spec:
  group: ocs.openshift.io
  names:
    kind: StorageCluster
    listKind: StorageClusterList
    plural: storageclusters
    singular: storagecluster
  scope: Namespaced
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          x-kubernetes-preserve-unknown-fields: true
EOF

    kubectl wait --for=condition=Established crd/storageclusters.ocs.openshift.io --timeout=60s

    kubectl apply -f - <<EOF
apiVersion: ocs.openshift.io/v1
kind: StorageCluster
metadata:
  name: ocs-storagecluster
  namespace: ${cluster_ns}
spec: {}
EOF
}

#################
# WAIT FUNCTIONS
#################

# Wait for ceph cluster to be ready
wait_for_ceph_cluster_to_be_ready_state() {
    local cluster_ns="$1"

    echo "Waiting for CephCluster to be ready in namespace $cluster_ns"
    if ! kubectl wait --for=condition=Ready cephcluster my-cluster -n "$cluster_ns" --timeout=${DEFAULT_TIMEOUT}s; then
        echo "CephCluster failed to become ready, current status:"
        kubectl get cephcluster -A
        exit 1
    fi
}

# Wait for operator pod to be ready
wait_for_operator_pod_to_be_ready_state() {
    local operator_ns="$1"

    echo "Waiting for operator pod to be ready in namespace $operator_ns"
    kubectl wait --for=condition=Ready pod -l app=rook-ceph-operator -n "$operator_ns" --timeout=${DEFAULT_TIMEOUT}s
}

########
# MAIN #
########

FUNCTION="$1"
shift # remove function arg now that we've recorded it

# Call the function with the remainder of the user-provided args
if ! "$FUNCTION" "$@"; then
    echo "Function '$FUNCTION' failed" >&2
    exit 1
fi

echo "Function '$FUNCTION' completed successfully!"
