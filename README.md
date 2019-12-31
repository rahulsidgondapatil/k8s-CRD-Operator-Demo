To generate code for this project:
Run following commands in ~/go/src/github.com/rahulsidgondapatil/sample-depSvcController directory
1) go mod init (optional if go.mod is absent)
2) go mod vendor (optional if vendor directory or k8s.io/code-generator is absent)
3) bash -x vendor/k8s.io/code-generator/generate-groups.sh all \
    "github.com/rahulsidgondapatil/sample-customController/pkg/client" \
    "github.com/rahulsidgondapatil/sample-customController/pkg/apis" \
    customcontroller:v1alpha1
