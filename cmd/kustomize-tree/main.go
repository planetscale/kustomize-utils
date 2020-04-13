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
package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/xlab/treeprint"
	"sigs.k8s.io/kustomize/v3/k8sdeps/validator"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/loader"

	"github.com/planetscale/kustomize-utils/pkg/tree"
)

var (
	relativeBase *string = flag.String("relative-base", ".", "show absolute paths unless paths below relative-base")
	listFiles    *bool   = flag.BoolP("list-files", "l", false, "list files instead of printing a tree")
	basesOnly    *bool   = flag.Bool("bases-only", false, "show only kustomize bases")
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s [overlay path]:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()
	var overlayP = "."
	if len(args) > 0 {
		overlayP = args[0]
	}

	fSys := fs.MakeFsOnDisk()
	ldr, err := loader.NewLoader(
		loader.RestrictionRootOnly,
		validator.NewKustValidator(),
		overlayP,
		fSys)
	if err != nil {
		log.Fatal(err)
	}

	if *relativeBase != "" {
		crb, _, err := fSys.CleanedAbs(*relativeBase)
		if err != nil {
			log.Fatal(err)
		}
		cwd := crb.String()
		relativeBase = &cwd
	}

	rt := tree.NewResourceTree(ldr, *relativeBase, *basesOnly)

	if err = rt.Build(); err != nil {
		log.Fatal(err)
	}

	if *listFiles {
		rt.Walk(func(v *treeprint.Vertex, level int) error {
			fmt.Println(v.GetValue().(string))
			return nil
		})
	} else {
		fmt.Println(rt.String())
	}
}
