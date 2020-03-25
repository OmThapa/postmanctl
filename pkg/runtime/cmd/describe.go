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

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kevinswiber/postmanctl/pkg/sdk/resources"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("describe called")
	},
}

func init() {
	describeCmd.AddCommand(
		generateDescribeSubcommand(resources.CollectionType, "collections", []string{"collection", "co"}),
	)
	rootCmd.AddCommand(describeCmd)
}

func generateDescribeSubcommand(t resources.ResourceType, use string, aliases []string) *cobra.Command {
	return &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeResource(t, args...)
		},
	}
}

func describeResource(resourceType resources.ResourceType, args ...string) error {
	switch resourceType {
	case resources.CollectionType:
		r := make(resources.CollectionSlice, len(args))
		for i, id := range args {
			resource, err := service.Collection(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		out, err := describeCollections(r)
		if err != nil {
			return err
		}

		fmt.Println(out)
	case resources.EnvironmentType:
		r := make(resources.EnvironmentSlice, len(args))
		for i, id := range args {
			resource, err := service.Environment(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.MockType:
		r := make(resources.MockSlice, len(args))
		for i, id := range args {
			resource, err := service.Mock(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.MonitorType:
		r := make(resources.MonitorSlice, len(args))
		for i, id := range args {
			resource, err := service.Monitor(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.APIType:
		r := make(resources.APISlice, len(args))
		for i, id := range args {
			resource, err := service.API(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.APIVersionType:
		apiID := args[0]
		ids := args[1:]

		r := make(resources.APIVersionSlice, len(ids))
		for i, id := range ids {
			resource, err := service.APIVersion(context.Background(), apiID, id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.WorkspaceType:
		r := make(resources.WorkspaceSlice, len(args))
		for i, id := range args {
			resource, err := service.Workspace(context.Background(), id)

			if err != nil {
				return handleResponseError(err)
			}

			r[i] = resource
		}

		printGetOutput(r)
	case resources.UserType:
		resource, err := service.User(context.Background())

		if err != nil {
			return handleResponseError(err)
		}

		printGetOutput(resource)
	case resources.SchemaType:
		apiID := args[0]
		apiVersionID := args[1]
		id := args[2]

		resource, err := service.Schema(context.Background(), apiID, apiVersionID, id)

		if err != nil {
			return handleResponseError(err)
		}

		printGetOutput(resource)
	default:
		return fmt.Errorf("invalid resource type: %s", resourceType.String())
	}

	return nil
}

func describeCollections(r resources.CollectionSlice) (string, error) {
	return tabbedString(func(out io.Writer) error {
		for _, c := range r {
			out.Write([]byte("Info:\n"))
			out.Write([]byte(fmt.Sprintf("  ID:\t%s\n", c.Info.PostmanId)))
			out.Write([]byte(fmt.Sprintf("  Name:\t%s\n", c.Info.Name)))
			out.Write([]byte(fmt.Sprintf("  Schema:\t%s\n", c.Info.Schema)))
			writeCollectionItemOrItemGroup(out, c.Item, "")
		}
		return nil
	})
}

func writeCollectionItemOrItemGroup(out io.Writer, c []interface{}, prefix string) {
	if len(c) > 0 {
		out.Write([]byte(fmt.Sprintf("%sItems:\n", prefix)))
		for idx, v := range c {
			var m map[string]interface{}
			m = v.(map[string]interface{})
			if _, ok := m["item"]; ok {
				writeCollectionItemGroup(out, m, prefix, idx)
				//writeCollectionItemOrItemGroup(out, m["item"].([]interface{}), prefix)
			} else {
				writeCollectionItem(out, m, prefix, idx)
				//fmt.Println("it's an item")
			}
		}
	}
}

func writeCollectionItemGroup(out io.Writer, c map[string]interface{}, prefix string, idx int) {
	tmp, err := json.Marshal(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var ci resources.ItemGroup
	if err := json.Unmarshal(tmp, &ci); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	level1Prefix := prefix + "  "
	level2Prefix := level1Prefix + "  "
	level3Prefix := level2Prefix + "  "
	level4Prefix := level3Prefix + "  "
	level5Prefix := level4Prefix + "  "
	level6Prefix := level5Prefix + "  "

	out.Write([]byte(fmt.Sprintf("%s%d:\n", level1Prefix, idx)))
	out.Write([]byte(fmt.Sprintf("%sType:\tfolder\n", level2Prefix)))
	out.Write([]byte(fmt.Sprintf("%sName:\t%s\n", level2Prefix, ci.Name)))
	writeCollectionItemOrItemGroup(out, ci.Item, level2Prefix)

	if len(ci.Event) > 0 {
		out.Write([]byte(fmt.Sprintf("%sEvents:\n", level2Prefix)))
	}

	for idx, ev := range ci.Event {
		out.Write([]byte(fmt.Sprintf("%s%d:\n", level3Prefix, idx)))
		out.Write([]byte(fmt.Sprintf("%sListen:\t%s\n", level4Prefix, ev.Listen)))
		out.Write([]byte(fmt.Sprintf("%sScript:\n", level4Prefix)))
		out.Write([]byte(fmt.Sprintf("%sID:\t%s\n", level5Prefix, ev.Script.Id)))
		out.Write([]byte(fmt.Sprintf("%sType:\t%s\n", level5Prefix, ev.Script.Type)))
		out.Write([]byte(fmt.Sprintf("%sExec:\n", level5Prefix)))
		if vSlice, ok := ev.Script.Exec.([]interface{}); ok {
			var execSlice = make([]string, len(vSlice))
			for i, s := range vSlice {
				execSlice[i] = s.(string)
			}

			out.Write([]byte(fmt.Sprintf("%s\t%s\n", level6Prefix, strings.Join(execSlice, "\n\t"))))
		} else if t, ok := ev.Script.Exec.(string); ok {
			out.Write([]byte(fmt.Sprintf("%s\t%s\n", level6Prefix, t)))
		}
	}
}

func writeCollectionItem(out io.Writer, c map[string]interface{}, prefix string, idx int) {
	tmp, err := json.Marshal(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var ci resources.Item
	if err := json.Unmarshal(tmp, &ci); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	level1Prefix := prefix + "  "
	level2Prefix := level1Prefix + "  "
	level3Prefix := level2Prefix + "  "
	level4Prefix := level3Prefix + "  "
	level5Prefix := level4Prefix + "  "
	level6Prefix := level5Prefix + "  "

	out.Write([]byte(fmt.Sprintf("%s%d:\n", level1Prefix, idx)))
	out.Write([]byte(fmt.Sprintf("%sID:\t%s\n", level2Prefix, c["_postman_id"])))
	out.Write([]byte(fmt.Sprintf("%sType:\trequest\n", level2Prefix)))
	out.Write([]byte(fmt.Sprintf("%sName:\t%s\n", level2Prefix, ci.Name)))

	if len(ci.Event) > 0 {
		out.Write([]byte(fmt.Sprintf("%sEvents:\n", level2Prefix)))
	}

	for idx, ev := range ci.Event {
		out.Write([]byte(fmt.Sprintf("%s%d:\n", level3Prefix, idx)))
		out.Write([]byte(fmt.Sprintf("%sListen:\t%s\n", level4Prefix, ev.Listen)))
		out.Write([]byte(fmt.Sprintf("%sScript:\n", level4Prefix)))
		out.Write([]byte(fmt.Sprintf("%sID:\t%s\n", level5Prefix, ev.Script.Id)))
		out.Write([]byte(fmt.Sprintf("%sType:\t%s\n", level5Prefix, ev.Script.Type)))
		out.Write([]byte(fmt.Sprintf("%sExec:\n", level5Prefix)))
		if vSlice, ok := ev.Script.Exec.([]interface{}); ok {
			var execSlice = make([]string, len(vSlice))
			for i, s := range vSlice {
				execSlice[i] = s.(string)
			}

			out.Write([]byte(fmt.Sprintf("%s\t%s\n", level6Prefix, strings.Join(execSlice, "\n\t"))))
		} else if t, ok := ev.Script.Exec.(string); ok {
			out.Write([]byte(fmt.Sprintf("%s\t%s\n", level6Prefix, t)))
		}
	}
}

func tabbedString(f func(io.Writer) error) (string, error) {
	out := new(tabwriter.Writer)
	buf := &bytes.Buffer{}
	out.Init(buf, 0, 8, 2, ' ', 0)

	err := f(out)
	if err != nil {
		return "", err
	}

	out.Flush()
	str := string(buf.String())
	return str, nil
}
