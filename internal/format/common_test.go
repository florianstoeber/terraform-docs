/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package format

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/terraform-docs/terraform-docs/internal/terraform"
	"github.com/terraform-docs/terraform-docs/internal/testutil"
)

func TestCommonSort(t *testing.T) {
	tests := map[string]struct {
		options terraform.Options
	}{
		"NoSort": {
			options: terraform.Options{},
		},
		"SortByName": {
			options: terraform.Options{
				SortBy: &terraform.SortBy{
					Name: true,
				},
			},
		},
		"SortByRequired": {
			options: terraform.Options{
				SortBy: &terraform.SortBy{
					Required: true,
				},
			},
		},
		"SortByType": {
			options: terraform.Options{
				SortBy: &terraform.SortBy{
					Type: true,
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			options, err := terraform.NewOptions().With(&tt.options)
			assert.Nil(err)

			module, err := testutil.GetModule(options)
			assert.Nil(err)

			type Expected struct {
				Inputs       []string `json:"inputs"`
				Modules      []string `json:"modules"`
				Outputs      []string `json:"outputs"`
				Providers    []string `json:"providers"`
				Requirements []string `json:"requirements"`
				Resources    []string `json:"resources"`
			}

			golden, err := testutil.GetExpected("common", "sort-"+name)
			assert.Nil(err)

			var expected Expected

			err = json.Unmarshal([]byte(golden), &expected)
			assert.Nil(err)

			for ii, i := range module.Inputs {
				assert.Equal(i.Name, expected.Inputs[ii])
			}
			for ii, m := range module.ModuleCalls {
				assert.Equal(m.Name+"-"+m.Source, expected.Modules[ii])
			}
			for ii, o := range module.Outputs {
				assert.Equal(o.Name, expected.Outputs[ii])
			}
			for ii, p := range module.Providers {
				assert.Equal(p.FullName(), expected.Providers[ii])
			}
			for ii, r := range module.Requirements {
				assert.Equal(r.Name, expected.Requirements[ii])
			}
			for ii, r := range module.Resources {
				assert.Equal(r.FullType()+"__"+r.Mode, expected.Resources[ii])
			}
		})
	}
}

func TestCommonHeaderFrom(t *testing.T) {
	tests := map[string]struct {
		options terraform.Options
	}{
		"HeaderFromADOCFile": {
			options: terraform.Options{
				HeaderFromFile: "doc.adoc",
			},
		},
		"HeaderFromMDFile": {
			options: terraform.Options{
				HeaderFromFile: "doc.md",
			},
		},
		"HeaderFromTFFile": {
			options: terraform.Options{
				HeaderFromFile: "doc.tf",
			},
		},
		"HeaderFromTXTFile": {
			options: terraform.Options{
				HeaderFromFile: "doc.txt",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			expected, err := testutil.GetExpected("common", "header-"+name)
			assert.Nil(err)

			options, err := terraform.NewOptions().WithOverwrite(&tt.options)
			assert.Nil(err)

			module, err := testutil.GetModule(options)
			assert.Nil(err)

			assert.Equal(module.Header, expected)
		})
	}
}
