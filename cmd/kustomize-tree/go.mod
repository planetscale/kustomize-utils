module github.com/planetscale/kustomize-utils/cmd/kustomize-tree

go 1.14

require (
	github.com/planetscale/kustomize-utils v0.0.0-00010101000000-000000000000
	github.com/spf13/pflag v1.0.5
	github.com/xlab/treeprint v1.0.0
	sigs.k8s.io/kustomize/v3 v3.3.1
)

replace github.com/xlab/treeprint => github.com/dctrwatson/treeprint v1.1.0

replace github.com/planetscale/kustomize-utils => ../../
