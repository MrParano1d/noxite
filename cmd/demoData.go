/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/mrparano1d/noxite/pkg/app"
	"github.com/mrparano1d/noxite/pkg/core/entities"

	"github.com/spf13/cobra"
)

func adminRoleJSON() string {
	return `{
		"CreateUser": true,
		"GetUser": true,
		"UpdateUser": true,
		"DeleteUser": true,

		"CreateRole": true,
		"GetRole": true,
		"UpdateRole": true,
		"DeleteRole": true,

		"PublishPackage": true,
		"GetPackage": true,
		"UpdatePackage": true,
		"UnpublishPackage": true
	}`
}

func userRoleJson() string {
	return `{
		"CreateUser": false,
		"GetUser": false,
		"UpdateUser": false,
		"DeleteUser": false,

		"CreateRole": false,
		"GetRole": false,
		"UpdateRole": false,
		"DeleteRole": false,

		"PublishPackage": true,
		"GetPackage": true,
		"UpdatePackage": true,
		"UnpublishPackage": true
	}`
}

// demoDataCmd represents the demoData command
var demoDataCmd = &cobra.Command{
	Use:   "demoData",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		entClient := app.EntClient()

		adminPermissions := entities.Permissions{}
		err := json.Unmarshal([]byte(adminRoleJSON()), &adminPermissions)
		if err != nil {
			panic(err)
		}

		adminRole, err := entClient.Role.Create().SetName("admin").SetPermissions(adminPermissions).Save(cmd.Context())
		if err != nil {
			panic(err)
		}

		_, err = entClient.User.Create().SetName("admin").SetEmail("admin@example.com").SetPassword([]byte("AdminPassw0rd!")).SetRoleID(adminRole.ID).Save(cmd.Context())
		if err != nil {
			panic(err)
		}

		userPermissions := entities.Permissions{}
		err = json.Unmarshal([]byte(userRoleJson()), &userPermissions)
		if err != nil {
			panic(err)
		}

		userRole, err := entClient.Role.Create().SetName("user").SetPermissions(userPermissions).Save(cmd.Context())
		if err != nil {
			panic(err)
		}

		_, err = entClient.User.Create().SetName("demo").SetEmail("demo@example.com").SetPassword([]byte("DemoPassw0rd!")).SetRoleID(userRole.ID).Save(cmd.Context())
		if err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(demoDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// demoDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// demoDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
