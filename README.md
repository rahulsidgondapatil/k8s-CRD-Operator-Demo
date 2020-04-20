This project demonstrates: How to create your own k8s Custom Resource Definition(CRD) and a k8s Operator for it.

How to generate code for this project:
Run following commands in ~/go/src/github.com/<your-directory-name>/k8s-crd-operator-demo directory
1) go mod init (optional if go.mod is absent)
2) go mod vendor (optional if vendor directory or k8s.io/code-generator is absent)
3) bash -x vendor/k8s.io/code-generator/generate-groups.sh all \
    "github.com/rahulsidgondapatil/k8s-crd-operator-demo/pkg/client" \
    "github.com/rahulsidgondapatil/k8s-crd-operator-demo/pkg/apis" \
    customcontroller:v1alpha1   \
    --go-header-file ./hack/boilerplate.go.txt
