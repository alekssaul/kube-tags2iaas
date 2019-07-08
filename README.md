# kube-tags2iaas

Experimental Project to sync Kubernetes Node labels with IaaS provider(s)

## Why do I need this

Today Pivotal Container Service (PKS), supports tagging of Virtual Machines on per "foundation" basis. Meaning the "custom tags" to PKS managed Kubernetes clusters cannot be individually assigned. i.e. master and worker nodes `Cluster01` created by PKS instance `api.pks.example.com` will get assigned mostly the same tags as master and worker nodes on `Cluster02`. The operator cannot have custom tags on a per cluster basis.

## What does this code do

This is a sample project that can be deployed onto a PKS managed cluster. This (relatively) simple golang application:
1. Reads the Kubernetes API to detect which nodes make up the cluster
1. Makes a call to Cloud Provider API (Azure) to tag the VMs specified by the operator
1. Stores tag information and other relavent information (i.e. last time it synced with Cloud API) as Kubernetes Node object annotations.

![Alt text](docs/img/Kube-2Iaas-bootstrap.png?raw=true )

Once the initial bootstrapping is complete, app watches for state changes to the Node objects in Kubernetes API in order to re-act to these events.

![Alt text](docs/img/Kube-2Iaas-watch.png?raw=true )

## Getting Started

Kubernetes deployment specs are stored in `deployments` folder of this repository.

Copy the example/template files

```sh
cp ./deployments/configmap.example ./deployments/configmap.yaml
cp ./deployments/secret.example ./deployments/secret.yaml
```

and edit those files.

`configmap.yaml` contains the Infrastructure TAGs to use in json array format. i.e.

```yaml
    foo: hello
    bar: world
```

becomes

```json
"{\"foo\": \"hello\",\"bar\": \"world\"}"
```

with proper escaping. This will tag the VMs with `foo=hello` and `bar=world` tags on Azure.

`secrets.yaml` contains the Secrets to connect to Azure API. 

```text
  AZURE_CLIENT_ID:              // Azure Client ID
  AZURE_CLIENT_SECRET:          // Azure Client Secret
  AZURE_GROUP_NAME:             // Azure Resource Group Name
  AZURE_LOCATION_DEFAULT:       // Azure Region
  AZURE_SUBSCRIPTION_ID:        // Azure Subscription ID
  AZURE_TENANT_ID:              // Azure Tenant ID
```

In general, above ServicePrincipal needs access to read Virtual Machines and update them with Tags.

Note that the Kubernetes secret object data values must be base64 encoded. i.e.

```sh
echo -n $AZURE_TENANT_ID | base64
```

output can be placed into the secret object.

Once the ConfigMap and Secret objects are updated, create a namespace (or use the default namespace) and update the `rbac.yaml` to make sure service account for the application is referred in correct namespace (line number 18).

```sh
kubectl create namespace kube2iaas
vi ./deployment/rbac.yaml
```

Finially, create the deployment and validate the output from the application

```sh
kubectl apply -f ./deployment
kubectl -n kube2iaas logs -l app=kube-tags2iaas
```
