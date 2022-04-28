# prow-plugin


## Creating a new package

Being the name of the package `horologium` and the version `0.1.0`

### Initializing the structure

```
❯ hack/create-package.sh horologium 0.1.0

```

### Check the vendir manifests

On `config/upstream` copy the correct upstream manifests. Edit the `vendir.yml` file to match the new files here with a manual type.

```
apiVersion: vendir.k14s.io/v1alpha1
kind: Config
minimumRequiredVersion: 0.12.0
directories:
  - path: config/upstream
    contents:
      - path: horologium_deployment.yaml
        manual: {}
      - path: horologium_rbac.yaml
        manual: {}
```

Run:

```
vendir sync
```

### Resolve image digest and lock

Lock the image digest from the repo

```
> hack/lock.sh horologium 0.1.0
```

### Push the imgpkg to the repository

Push the package bundle to public.ecr.aws/t0q8k6g2/prow-plugin

```
hack/push-package.sh hook 0.1.0
```

#### Save the imgpkg hash

Save the imgpkg has from the above command output and update packages/${package}/${version}/package.yaml

### Push repo bundle

Update repo bundle in repos/prow.yaml.

```
cd hack
go run generate-package-repository.go prow 0.1.0
```
