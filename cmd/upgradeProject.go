/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/sveltinio/prompti/confirm"
	"github.com/sveltinio/sveltin/internal/markup"
	"github.com/sveltinio/sveltin/internal/migrations"
	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/tui/feedbacks"
	"github.com/sveltinio/sveltin/utils"
)

// Migration identifiers.
const (
	ProjectSettingsMigrationID string = "projectSettings"
	DefaultsConfigMigrationID  string = "defaultsConfig"
	ThemeConfigMigrationID     string = "themeConfig"
	DotEnvMigrationID          string = "dotenv"
	PackageJSONID              string = "packagejson"
	MDsveXID                   string = "mdsvex"
)

//=============================================================================

var upgradeProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Upgrade your project to the latest Sveltin version",
	Long: resources.GetASCIIArt() + `
Command used to upgrade your project files to the latest Sveltin version.
`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(0),
	Run:                   RunUpgradeProjectCmd,
}

// RunUpgradeProjectCmd is the actual work function.
func RunUpgradeProjectCmd(cmd *cobra.Command, args []string) {
	// Exit if running sveltin commands from a not valid directory.
	isValidProject(false)

	feedbacks.ShowUpgradeCommandMessages()

	isConfirm, err := confirm.Run(&confirm.Config{Question: "Continue?"})
	utils.ExitIfError(err)

	if isConfirm {
		cwd, _ := os.Getwd()
		cfg.log.Plain(markup.H1(fmt.Sprintf("Upgrading your project to sveltin v%s", CliVersion)))

		migrationManager := migrations.NewMigrationManager()
		migrationServices := migrations.NewMigrationServices(cfg.fs, cfg.fsManager, cfg.pathMaker, cfg.log)

		/** FILE: <project_root>/sveltin.json */
		pathToFile := path.Join(cwd, ProjectSettingsFile)
		migrationData := &migrations.MigrationData{
			PathToFile:        pathToFile,
			CliVersion:        CliVersion,
			ProjectCliVersion: cfg.projectSettings.Sveltin.Version,
		}
		migrationFactory, err := migrations.GetMigrationFactory(migrations.ProjectSettingsMigrationID)
		utils.ExitIfError(err)
		migration := migrationFactory.MakeMigration(migrationManager, migrationServices, migrationData)
		// execute the migration.
		err = migration.Execute()
		utils.ExitIfError(err)

		// Load project settings file after sveltin.json file creation
		cfg.projectSettings, err = loadProjectSettings(ProjectSettingsFile)
		utils.ExitIfError(err)

		migrationIdPathToFileMap := map[string]string{
			migrations.DefaultsConfigMigrationID: path.Join(cwd, cfg.pathMaker.GetConfigFolder(), DefaultsConfigFile),
			migrations.ThemeConfigMigrationID:    path.Join(cwd, cfg.pathMaker.GetThemesFolder(), cfg.projectSettings.Theme.Name, cfg.settings.GetThemeConfigFilename()),
			migrations.DotEnvMigrationID:         path.Join(cwd, DotEnvProdFile),
			migrations.PackageJSONMigrationID:    path.Join(cwd, PackageJSONFile),
			migrations.MDsveXMigrationID:         path.Join(cwd, MDsveXFile),
			migrations.SvelteConfigMigrationID:   path.Join(cwd, SvelteConfigFile),
			migrations.LayoutMigrationID:         path.Join(cwd, cfg.pathMaker.GetRoutesFolder(), LayoutTSFile),
		}

		for id, pathToFile := range migrationIdPathToFileMap {
			migrationData := &migrations.MigrationData{
				PathToFile: pathToFile,
			}
			migrationFactory, err := migrations.GetMigrationFactory(id)
			utils.ExitIfError(err)
			migration := migrationFactory.MakeMigration(migrationManager, migrationServices, migrationData)
			// execute the migration.
			err = migration.Execute()
			utils.ExitIfError(err)
		}

		cfg.log.Success(fmt.Sprintf("Your project is ready for sveltin v%s\n", CliVersion))
	}
}

func init() {
	upgradeCmd.AddCommand(upgradeProjectCmd)
}
