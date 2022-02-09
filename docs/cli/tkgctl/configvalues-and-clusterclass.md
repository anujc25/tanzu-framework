# Config Values and Clusterclass

# Management cluster creation

In this section we will be explaining how the user configuration will be passed to create clusterclass using packages.

1. User passes the cluster configuration in the flat `key:value` format.

```
CLUSTER_NAME: aws-mc-1
CLUSTER_PLAN: dev
INFRASTRUCTURE_PROVIDER: aws

AWS_VPC_ID: ""
AWS_SSH_KEY_NAME: tkgtest
AWS_PRIVATE_NODE_CIDR: 10.0.1.0/24
...
```

1. Tanzu CLI management-cluster creation workflow
    1. create kind cluster
    1. install CAPI and CAPx providers to the kind cluster
    1. **install clusterclass package to the kind cluster for the given infrastructure**
    1. **apply `cluster` resource to the kind cluster to create management cluster**
    1. wait for management-cluster to get created
    1. once management-cluster gets created
    1. **install clusterclass package to the management cluster for the given infrastructure**
    1. move all cluster-api resources from the bootstrap cluster to management cluster