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

#### Save the imgpkg hash
