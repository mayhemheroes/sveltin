/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package migrations

import (
	"fmt"
	"path"
	"strings"

	"github.com/sveltinio/sveltin/resources"
	"github.com/sveltinio/sveltin/utils"
)

// OverwriteSveltinDTS is the struct representing the migration add the sveltin.json file.
type OverwriteSveltinDTS struct {
	Mediator IMigrationMediator
	Services *MigrationServices
	Data     *MigrationData
}

// MakeMigration implements IMigrationFactory interface.
func (m *OverwriteSveltinDTS) MakeMigration(migrationManager *MigrationManager, services *MigrationServices, data *MigrationData) IMigration {
	return &OverwriteSveltinDTS{
		Mediator: migrationManager,
		Services: services,
		Data:     data,
	}
}

// MakeMigration implements IMigration interface.
func (m *OverwriteSveltinDTS) getServices() *MigrationServices { return m.Services }
func (m *OverwriteSveltinDTS) getData() *MigrationData         { return m.Data }

// Migrate return error if migration execution over up and down methods fails.
func (m OverwriteSveltinDTS) Migrate() error {
	if err := m.up(); err != nil {
		return err
	}
	if err := m.down(); err != nil {
		return err
	}
	return nil
}

func (m *OverwriteSveltinDTS) up() error {
	if !m.Mediator.canRun(m) {
		return nil
	}

	exists, err := utils.FileExists(m.getServices().fs, m.Data.TargetPath)
	if !exists {
		return err
	}

	if exists {
		fileContent, err := retrieveFileContent(m.getServices().fs, m.getData().TargetPath)
		if err != nil {
			return err
		}

		gatekeeper := patterns[sveltindts]
		if mustMigrate(fileContent, gatekeeper) {
			localFilePath :=
				strings.Replace(m.Data.TargetPath, m.getServices().pathMaker.GetRootFolder(), "", 1)
			m.getServices().logger.Info(fmt.Sprintf("Migrating %s", localFilePath))
			saveTo := path.Join(m.Services.pathMaker.GetSrcFolder())
			return m.Services.fsManager.CopyFileFromEmbed(&resources.SveltinStaticFS, m.Services.fs, resources.SveltinFilesFS, "sveltin_d_ts", saveTo)
		}
	}

	return nil
}

func (m *OverwriteSveltinDTS) down() error {
	if err := m.Mediator.notifyAboutCompletion(); err != nil {
		return err
	}
	return nil
}

func (m *OverwriteSveltinDTS) allowUp() error {
	if err := m.up(); err != nil {
		return err
	}
	return nil
}

func (m *OverwriteSveltinDTS) runMigration(content []byte, file string) ([]byte, error) {
	return nil, nil
}

//=============================================================================
