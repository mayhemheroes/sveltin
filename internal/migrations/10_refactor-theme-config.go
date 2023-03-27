/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package migrations

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/utils"
)

// RefactorThemeConfig is the struct representing the migration update the defaults.js.ts file.
type RefactorThemeConfig struct {
	Mediator IMigrationMediator
	Services *MigrationServices
	Data     *MigrationData
}

// MakeMigration implements IMigrationFactory interface.
func (m *RefactorThemeConfig) MakeMigration(migrationManager *MigrationManager, services *MigrationServices, data *MigrationData) IMigration {
	return &RefactorThemeConfig{
		Mediator: migrationManager,
		Services: services,
		Data:     data,
	}
}

// implements IMigration interface.
func (m *RefactorThemeConfig) getServices() *MigrationServices { return m.Services }
func (m *RefactorThemeConfig) getData() *MigrationData         { return m.Data }

// Migrate return error if migration execution over up and down methods fails.
func (m RefactorThemeConfig) Migrate() error {
	if err := m.up(); err != nil {
		return err
	}
	if err := m.down(); err != nil {
		return err
	}
	return nil
}

func (m *RefactorThemeConfig) up() error {
	if !m.Mediator.canRun(m) {
		return nil
	}

	exists, err := utils.FileExists(m.getServices().fs, m.Data.TargetPath)
	if err != nil {
		return err
	}

	if exists {
		fileContent, err := retrieveFileContent(m.getServices().fs, m.getData().TargetPath)
		if err != nil {
			return err
		}

		migrationTriggers := []string{
			patterns[themeConfigConst],
		}
		if patternsMatched(fileContent, migrationTriggers, findStringMatcher) {
			localFilePath :=
				strings.Replace(m.Data.TargetPath, m.getServices().pathMaker.GetRootFolder(), "", 1)
			m.getServices().logger.Info(fmt.Sprintf("Migrating %s", localFilePath))
			if _, err := m.runMigration(fileContent, ""); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RefactorThemeConfig) down() error {
	if err := m.Mediator.notifyAboutCompletion(); err != nil {
		return err
	}
	return nil
}

func (m *RefactorThemeConfig) allowUp() error {
	if err := m.up(); err != nil {
		return err
	}
	return nil
}

func (m *RefactorThemeConfig) runMigration(content []byte, filepath string) ([]byte, error) {
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var prevLine string
		if i > 0 {
			prevLine = lines[i-1]
		}

		rules := []*migrationRule{
			newConstNameRule(line),
			newExportLineRule(line),
			newThemeNameRule(line, prevLine),
		}
		if res, ok := applyMigrationRules(rules); ok {
			lines[i] = res
		} else {
			lines[i] = line
		}
	}

	output := strings.Join(lines, "\n")
	err := m.getServices().fs.Remove(m.Data.TargetPath)
	if err != nil {
		return nil, err
	}

	if err = afero.WriteFile(m.getServices().fs, m.Data.TargetPath, []byte(output), 0644); err != nil {
		return nil, err
	}
	return nil, nil
}

//=============================================================================

func newConstNameRule(line string) *migrationRule {
	return &migrationRule{
		value:           line,
		trigger:         patterns[themeConfigConst],
		replaceFullLine: true,
		replacerFunc: func(string) string {
			return `import { theme } from '../../sveltin.json';

const themeConfig = {`
		},
	}
}

func newExportLineRule(line string) *migrationRule {
	return &migrationRule{
		value:           line,
		trigger:         patterns[themeConfigExport],
		replaceFullLine: false,
		replacerFunc: func(string) string {
			return "export { themeConfig }"
		},
	}
}

func newThemeNameRule(line, prevLine string) *migrationRule {
	return &migrationRule{
		value:           line,
		trigger:         patterns[themeNameProp],
		replaceFullLine: true,
		replacerFunc: func(string) string {
			if !strings.Contains(prevLine, "author:") && strings.Contains(line, "name:") {
				return "\tname: theme.name,"
			}
			return line
		},
	}
}
