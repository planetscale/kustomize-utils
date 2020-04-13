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
package kustomize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/pgmconfig"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

func LoadKustomization(ldr ifc.Loader) (*types.Kustomization, error) {
	content, err := readKustFile(ldr)
	if err != nil {
		return nil, err
	}
	content = types.FixKustomizationPreUnmarshalling(content)
	var k types.Kustomization
	err = unmarshal(content, &k)
	if err != nil {
		return nil, err
	}
	k.FixKustomizationPostUnmarshalling()
	errs := k.EnforceFields()
	if len(errs) > 0 {
		return nil, fmt.Errorf(
			"Failed to read kustomization file under %s:\n"+
				strings.Join(errs, "\n"), ldr.Root())
	}
	return &k, nil
}

func unmarshal(y []byte, o interface{}) error {
	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(j))
	dec.DisallowUnknownFields()
	return dec.Decode(o)
}

func readKustFile(ldr ifc.Loader) ([]byte, error) {
	var content []byte
	match := 0
	for _, kf := range pgmconfig.RecognizedKustomizationFileNames() {
		c, err := ldr.Load(kf)
		if err == nil {
			match += 1
			content = c
		}
	}

	switch match {
	case 0:
		return nil, fmt.Errorf("Unable to find one of these in directory '%s':\n%s\n",
			ldr.Root(), strings.Join(pgmconfig.RecognizedKustomizationFileNames(), "\n"))
	case 1:
		return content, nil
	default:
		return nil, fmt.Errorf("Found multiple kustomization files under: %s\n", ldr.Root())
	}
}
