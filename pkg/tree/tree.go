/*
Copyright 2020 PlanetScale Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tree

import (
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"

	"github.com/planetscale/kustomize-utils/pkg/kustomize"
)

type ResourceTree struct {
	ldr ifc.Loader

	relativeBase string
	basesOnly    bool

	tree treeprint.Tree
}

func NewResourceTree(ldr ifc.Loader, relativeBase string, basesOnly bool) *ResourceTree {
	if relativeBase != "" && !strings.HasSuffix(relativeBase, "/") {
		relativeBase += "/"
	}

	tree := treeprint.New()
	tree.SetValue(strings.TrimPrefix(ldr.Root(), relativeBase))

	return &ResourceTree{
		ldr:          ldr,
		relativeBase: relativeBase,
		basesOnly:    basesOnly,
		tree:         tree,
	}
}

func (rt *ResourceTree) Build() error {
	return rt.loadResources(rt.ldr, rt.tree)
}

func (rt *ResourceTree) loadResources(ldr ifc.Loader, branch treeprint.Tree) error {
	kb, err := kustomize.LoadKustomization(ldr)
	if err != nil {
		return err
	}

	trimmedRoot := strings.TrimPrefix(ldr.Root(), rt.relativeBase)

	for _, path := range kb.Resources {
		rldr, err := ldr.New(path)
		if err != nil {
			// err means it's a resource not a base
			if !rt.basesOnly {
				branch.AddNode(filepath.Join(trimmedRoot, path))
			}
		} else {
			branch2 := branch.AddBranch(strings.TrimPrefix(rldr.Root(), rt.relativeBase))
			err2 := rt.loadResources(rldr, branch2)
			if err2 != nil {
				return err2
			}
		}
	}

	return nil
}

func (rt *ResourceTree) String() string {
	return rt.tree.String()
}

func (rt *ResourceTree) Walk(walkFn treeprint.TreeWalkFn) error {
	return rt.tree.Walk(walkFn)
}
