To generate code for this project:
Run following commands in ~/go/src/github.com/rahulsidgondapatil/k8s-CRD-Operator-Demo directory
1) go mod init (optional if go.mod is absent)
2) go mod vendor (optional if vendor directory or k8s.io/code-generator is absent)
3) bash -x vendor/k8s.io/code-generator/generate-groups.sh all \
    "github.com/rahulsidgondapatil/k8s-CRD-Operator-Demo/pkg/client" \
    "github.com/rahulsidgondapatil/k8s-CRD-Operator-Demo/pkg/apis" \
    customcontroller:v1alpha1   \
    --go-header-file ./hack/boilerplate.go.txt
