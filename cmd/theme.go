/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/sveltinio/sveltin/config"
	"github.com/sveltinio/sveltin/helpers/factory"
	"github.com/sveltinio/sveltin/internal/composer"
	"github.com/sveltinio/sveltin/internal/css"
	sveltinerr "github.com/sveltinio/sveltin/internal/errors"
	"github.com/sveltinio/sveltin/internal/markup"
	"github.com/sveltinio/sveltin/internal/shell"
	"github.com/sveltinio/sveltin/internal/tpltypes"
	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/tui/activehelps"
	"github.com/sveltinio/sveltin/tui/feedbacks"
	"github.com/sveltinio/sveltin/tui/prompts"
	"github.com/sveltinio/sveltin/utils"
)

var (
	// How to use the command.
	themeCmdExample = "sveltin new theme paper --css tailwindcss"
	// Short description shown in the 'help' output.
	themeCmdShortMsg = "Create a new theme reusable theme"
	// Long message shown in the 'help <this-command>' output.
	themeCmdLongMsg = utils.MakeCmdLongMsg("Command used to create a new theme for projects so that can be shared with others and reused.")
)

//=============================================================================

var themeCmd = &cobra.Command{
	Use:               "theme [name]",
	Aliases:           []string{"t"},
	Example:           themeCmdExample,
	Short:             themeCmdShortMsg,
	Long:              themeCmdLongMsg,
	ValidArgsFunction: themeCmdValidArgFunction,
	PersistentPreRun:  themeCmdPreRunHook,
	Run:               ThemeCmdRun,
}

// ThemeCmdRun is the actual work function.
func ThemeCmdRun(cmd *cobra.Command, args []string) {
	themeName, err := prompts.AskThemeName(args)
	utils.ExitIfError(err)
	cfg.log.Info(themeName)

	projectName := themeName + "_project"

	cssLibName, err := prompts.SelectCSSLibHandler(withCSSLib)
	utils.ExitIfError(err)
	cfg.log.Info(cssLibName)

	npmClient := getSelectedNPMClient(npmClientName, cfg.log)
	npmClientName = npmClient.Name

	cfg.log.Plain(markup.H1("A Starter project will be created"))

	// Clone starter template github repository
	themeStarterTemplate := cfg.startersMap[ThemeStarter]
	cfg.log.Info(fmt.Sprintf("Cloning the %s repos", themeStarterTemplate.Name))

	gitClient := shell.NewGitClient()
	err = gitClient.RunClone(themeStarterTemplate.URL, cfg.pathMaker.GetProjectRoot(projectName), true)
	utils.ExitIfError(err)

	// NEW FILE: config/defaults.js
	f := cfg.fsManager.NewConfigFile(projectName, Defaults, CliVersion)
	// NEW FOLDER: config
	configFolder := composer.NewFolder(ConfigFolder)
	configFolder.Add(f)

	// MAKE FOLDER STRUCTURE: themes/<theme_name> folder
	themeData := &tpltypes.ThemeData{
		ID:     tpltypes.BlankTheme,
		IsNew:  true,
		Name:   themeName,
		CSSLib: cssLibName,
	}
	themeFolder, err := makeProjectFolderStructure(ThemesFolder, "", themeData)
	utils.ExitIfError(err)

	// SET FOLDER STRUCTURE
	projectFolder := cfg.fsManager.GetFolder(projectName)
	projectFolder.Add(configFolder)
	projectFolder.Add(themeFolder)

	rootFolder := cfg.fsManager.GetFolder(RootFolder)
	rootFolder.Add(projectFolder)

	// GENERATE THE FOLDER TREE
	sfs := factory.NewThemeArtifact(&resources.SveltinTemplatesFS, cfg.fs)
	err = rootFolder.Create(sfs)
	utils.ExitIfError(err)

	// SETUP THE CSS LIB
	cfg.log.Info("Setting up the CSS Lib")
	tplData := config.TemplateData{
		ProjectName: projectName,
		NPMClient: &tpltypes.NPMClientData{
			Name:    npmClient.Name,
			Version: npmClient.Version,
			Info:    npmClient.ToString(),
		},
		Vite: &tpltypes.ViteData{
			Port: withPortNumber,
		},
		Theme: themeData,
	}
	err = setupThemeCSSLib(&resources.SveltinTemplatesFS, cfg, &tplData)
	utils.ExitIfError(err)

	cfg.log.Success("Done\n")

	// NEXT STEPS
	projectConfigSummary := config.NewProjectConfig(projectName, cssLibName, themeName, npmClient.Desc)
	feedbacks.ShowNewThemeHelpMessage(projectConfigSummary)
}

// Command initialization.
func init() {
	themeCmdFlags(themeCmd)
	//rootCmd.AddCommand(themeCmd)
}

//=============================================================================

// Assign flags to the command.
func themeCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&withCSSLib, "css", "c", "",
		fmt.Sprintf("The CSS lib to use. Valid: %s, %s, %s, %s, %s, %s",
			css.Bootstrap, css.Bulma, css.Scss, css.TailwindCSS, css.UnoCSS, css.VanillaCSS))
	cmd.Flags().StringVarP(&npmClientName, "npmClient", "n", "", "The name of your preferred npm client")
}

// Adding Active Help messages enhancing shell completions.
func themeCmdValidArgFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var comps []string
	if len(args) == 0 {
		comps = cobra.AppendActiveHelp(comps, activehelps.Hint("You must choose a name for the theme"))
	} else {
		comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments but accepts flags")
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

// Run before the main Run function of init command to check and alert about newer version.
func themeCmdPreRunHook(cmd *cobra.Command, args []string) {
	pwd, _ := os.Getwd()
	pathToPkgJSON := filepath.Join(pwd, "package.json")
	exists, _ := afero.Exists(cfg.fs, pathToPkgJSON)
	if exists {
		err := sveltinerr.NewNotEmptyProjectError(pathToPkgJSON)
		log.Fatalf("\x1b[31;1m✘ %s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	}
}

//=============================================================================

func setupThemeCSSLib(efs *embed.FS, cfg appConfig, tplData *config.TemplateData) error {
	switch tplData.Theme.CSSLib {
	case css.Bootstrap:
		boostrap := css.NewBootstrap(efs, cfg.fs, cfg.settings, tplData)
		return boostrap.Setup(false)
	case css.Bulma:
		bulma := css.NewBulma(efs, cfg.fs, cfg.settings, tplData)
		return bulma.Setup(false)
	case css.Scss:
		scss := css.NewScss(efs, cfg.fs, cfg.settings, tplData)
		return scss.Setup(false)
	case css.TailwindCSS:
		tailwind := css.NewTailwindCSS(efs, cfg.fs, cfg.settings, tplData)
		return tailwind.Setup(false)
	case css.VanillaCSS:
		vanillaCSS := css.NewVanillaCSS(efs, cfg.fs, cfg.settings, tplData)
		return vanillaCSS.Setup(false)
	default:
		return sveltinerr.NewOptionNotValidError(
			tplData.Theme.CSSLib,
			[]string{
				css.Bootstrap,
				css.Bulma,
				css.Scss,
				css.TailwindCSS,
				css.UnoCSS,
				css.VanillaCSS,
			})
	}
}