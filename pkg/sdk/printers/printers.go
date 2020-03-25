/*
Copyright © 2020 Kevin Swiber <kswiber@gmail.com>

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

package printers

import (
	"io"

	"github.com/kevinswiber/postmanctl/pkg/sdk/resources"
)

// ResourcePrinter writes an output of API resources.
type ResourcePrinter interface {
	PrintResource(resources.Formatter, io.Writer)
}

// PrintOptions holds various options used in ResourcePrinters.
type PrintOptions struct {
	NoHeaders bool
}