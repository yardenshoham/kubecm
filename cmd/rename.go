/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var oldName string
var newName string

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename the contexts of kubeconfig",
	Long: `
# Renamed dev to test
kubecm rename -o dev -n test
# Renamed current-context name to dev
kubecm rename -n dev -c
`,
	Run: func(cmd *cobra.Command, args []string) {
		cover, _ = cmd.Flags().GetBool("cover")
		if cover && oldName != "" {
			fmt.Println("parameter `-c` and `-n` cannot be set at the same time")
			os.Exit(1)
		} else {
			config, err := LoadClientConfig(cfgFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			if cover {
				for key, obj := range config.Contexts {
					if current := config.CurrentContext; key == current {
						config.Contexts[newName] = obj
						delete(config.Contexts, key)
						config.CurrentContext = newName
						fmt.Println(fmt.Sprintf("Rename %s to %s", key, newName))
						break
					}
				}
			} else {
				if obj, ok := config.Contexts[oldName]; ok {
					config.Contexts[newName] = obj
					delete(config.Contexts, oldName)
					if config.CurrentContext == oldName {
						config.CurrentContext = newName
					}
				} else {
					fmt.Println(fmt.Sprintf("Can not find context: %s", oldName))
					err := Formatable(nil)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					os.Exit(-1)
				}
			}
			err = ModifyKubeConfig(config)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		err := Formatable(nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&oldName, "old", "o", "", "Old context name")
	renameCmd.Flags().StringVarP(&newName, "new", "n", "", "New context name")
	renameCmd.Flags().BoolP("cover", "c", false, "")
	renameCmd.MarkFlagRequired("new")
}