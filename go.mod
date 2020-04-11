module github.com/planetscale/kustomize-utils

go 1.14

require (
	github.com/xlab/treeprint v1.0.0
	sigs.k8s.io/kustomize/v3 v3.3.1
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/xlab/treeprint => github.com/dctrwatson/treeprint v1.1.0
