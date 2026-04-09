# Object commands

## enable remote-obc

Install the ObjectBucket (OB) and ObjectBucketClaim (OBC) CRDs,
and if the CRDs already exist on the cluster, they will be updated.

### Example

```bash
odf object enable remote-obc

# Info: CRD "objectbuckets.objectbucket.io" installed successfully
# Info: CRD "objectbucketclaims.objectbucket.io" installed successfully
```

## disable remote-obc

Remove the ObjectBucket (OB) and ObjectBucketClaim (OBC) CRDs from the cluster.

### Example

```bash
odf object disable remote-obc

# Info: CRD "objectbuckets.objectbucket.io" deleted successfully
# Info: CRD "objectbucketclaims.objectbucket.io" deleted successfully
```
