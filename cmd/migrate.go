package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"core/database"
	"core/models"
)

var (
	tenants bool
)

var cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the models",
	Long:  `Migrate the models to the database schema`,
	Run: func(cmd *cobra.Command, args []string) {
		db := database.GetCentralConnection()
		if !tenants {
			database.MigrateCentral(db)
			return
		}
		var tenantsList []models.Tenant
		if err := db.Find(&tenantsList).Error; err != nil {
			fmt.Println("Error fetching tenants:", err)
			return
		}
		for _, tenant := range tenantsList {
			fmt.Printf("Migrating database for tenant %s\n", *tenant.StoreName)
			db := database.GetTenantConnection(*tenant.StoreName)
			database.MigrateTenants(db)
		}
	},
}

func init() {
	cmdMigrate.Flags().BoolVarP(&tenants, "tenants", "t", false, "specify whether to migrate for tenants")
	rootCmd.AddCommand(cmdMigrate)
}
