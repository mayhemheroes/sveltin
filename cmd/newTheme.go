/**
 * Copyright © 2021 Mirco Veltri <github@mircoveltri.me>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package cmd ...
package cmd

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/sveltinio/sveltin/common"
	"github.com/sveltinio/sveltin/config"
	"github.com/sveltinio/sveltin/helpers/factory"
	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/sveltinlib/composer"
	"github.com/sveltinio/sveltin/sveltinlib/css"
	"github.com/sveltinio/sveltin/sveltinlib/shell"
	"github.com/sveltinio/sveltin/sveltinlib/sveltinerr"
	"github.com/sveltinio/sveltin/utils"
)

//=============================================================================

var newThemeCmd = &cobra.Command{
	Use:     "theme <name>",
	Aliases: []string{"t"},
	Short:   "Command to create a new theme",
	Long: resources.GetASCIIArt() + `
This command help you creating new themes for projects, so that can be shared with others and reused.

Examples:

sveltin new theme paper
sveltin new theme paper --css tailwindcss
`,
	Run: NewThemeCmdRun,
}

// NewThemeCmdRun is the actual work function.
func NewThemeCmdRun(cmd *cobra.Command, args []string) {
	// Exit if running the command from an existing sveltin project folder.
	isValidForThemeMaker()

	themeName, err := promptThemeName(args)
	utils.ExitIfError(err)
	log.Info(themeName)

	projectName := themeName + "_project"

	cssLibName, err := promptCSSLibName(withCSSLib)
	utils.ExitIfError(err)
	log.Info(cssLibName)

	npmClient := getSelectedNPMClient()
	npmClientName = npmClient.Name

	log.Plain(utils.Underline("A Starter project will be created"))

	// Clone starter template github repository
	themeStarterTemplate := appTemplatesMap[ThemeStarter]
	log.Info(fmt.Sprintf("Cloning the %s repos", themeStarterTemplate.Name))
	gitClient := shell.NewGitClient()
	err = gitClient.RunGitClone(themeStarterTemplate.URL, pathMaker.GetProjectRoot(projectName), true)
	//err = utils.GitClone(themeStarterTemplate.URL, pathMaker.GetProjectRoot(projectName))
	utils.ExitIfError(err)

	// GET FOLDER: <project_name>
	projectFolder := fsManager.GetFolder(projectName)

	// NEW FILE: config/defaults.js
	f := fsManager.NewConfigFile(projectName, Defaults, CliVersion)
	// NEW FOLDER: config
	configFolder := composer.NewFolder(Config)
	configFolder.Add(f)
	projectFolder.Add(configFolder)

	// NEW FOLDER: themes
	themesFolder := composer.NewFolder(Themes)

	themeData := &config.ThemeData{
		ID:     config.BlankTheme,
		IsNew:  true,
		Name:   themeName,
		CSSLib: cssLibName,
	}
	newThemeFolder := makeThemeFolderStructure(themeData)
	themesFolder.Add(newThemeFolder)
	// ADD themes folder to the project
	projectFolder.Add(themesFolder)

	// SET FOLDER STRUCTURE
	rootFolder := fsManager.GetFolder(Root)
	rootFolder.Add(projectFolder)

	// GENERATE FOLDER STRUCTURE
	sfs := factory.NewThemeArtifact(&resources.SveltinFS, AppFs)
	err = rootFolder.Create(sfs)
	utils.ExitIfError(err)

	// SETUP THE CSS LIB
	log.Info("Setting up the CSS Lib")
	tplData := config.TemplateData{
		ProjectName: projectName,
		NPMClient:   npmClient.ToString(),
		PortNumber:  withPortNumber,
		Theme:       themeData,
	}
	err = setupThemeCSSLib(&resources.SveltinFS, AppFs, &conf, &tplData)
	utils.ExitIfError(err)

	log.Success("Done")

	// NEXT STEPS
	log.Plain(utils.Underline("Next Steps"))
	log.Plain(common.HelperTextNewTheme(projectName))
}

func newThemeCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&withCSSLib, "css", "c", "", "The name of the CSS framework to use. Possible values: vanillacss, tailwindcss, bulma, bootstrap, scss")
	cmd.Flags().StringVarP(&npmClientName, "npmClient", "n", "", "The name of your preferred npm client")
}

func init() {
	newThemeCmdFlags(newThemeCmd)
	newCmd.AddCommand(newThemeCmd)
}

//=============================================================================

// isValidForThemeMaker returns error if find the package.json file within the current folder.
func isValidForThemeMaker() {
	pwd, _ := os.Getwd()
	pathToPkgJSON := filepath.Join(pwd, "package.json")
	exists, _ := afero.Exists(AppFs, pathToPkgJSON)
	if exists {
		err := sveltinerr.NewNotEmptyProjectError(pathToPkgJSON)
		jww.FATAL.Fatalf("\x1b[31;1m✘ %s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	}
}

func promptThemeName(inputs []string) (string, error) {
	switch numOfArgs := len(inputs); {
	case numOfArgs < 1:
		themeNamePromptContent := config.PromptContent{
			ErrorMsg: "Please, provide a name for the theme.",
			Label:    "What's the theme name?",
		}
		result, err := common.PromptGetInput(themeNamePromptContent, nil, "")
		if err != nil {
			return "", err
		}
		return utils.ToSlug(result), nil
	case numOfArgs == 1:
		return utils.ToSlug(inputs[0]), nil
	default:
		err := errors.New("something went wrong: value not valid")
		return "", sveltinerr.NewDefaultError(err)
	}
}

func setupThemeCSSLib(efs *embed.FS, fs afero.Fs, conf *config.SveltinConfig, tplData *config.TemplateData) error {
	switch tplData.Theme.CSSLib {
	case VanillaCSS:
		vanillaCSS := css.NewVanillaCSS(efs, fs, conf, tplData)
		return vanillaCSS.Setup(false)
	case Scss:
		scss := css.NewScss(efs, fs, conf, tplData)
		return scss.Setup(false)
	case TailwindCSS:
		tailwind := css.NewTailwindCSS(efs, fs, conf, tplData)
		return tailwind.Setup(false)
	case Bulma:
		bulma := css.NewBulma(efs, fs, conf, tplData)
		return bulma.Setup(false)
	case Bootstrap:
		boostrap := css.NewBootstrap(efs, fs, conf, tplData)
		return boostrap.Setup(false)
	default:
		return sveltinerr.NewOptionNotValidError(tplData.Theme.CSSLib, []string{"vanillacss", "tailwindcss", "bulma", "bootstrap", "scss"})
	}
}
