/**
 * Copyright © 2021-present Sveltin contributors <github@sveltin.io>
 *
 * Use of this source code is governed by Apache 2.0 license
 * that can be found in the LICENSE file.
 */

// Package migrations implements the Mediator design pattern used to manage migrations over sveltin versions.
package migrations

import (
	"bytes"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/sveltinio/sveltin/internal/fsm"
	"github.com/sveltinio/sveltin/internal/pathmaker"
)

type matcherFunc = func([]byte, string, string) bool

// IMigration is the interface defining the methods to be implemented by single migration.
type IMigration interface {
	Migrate() error
	getServices() *MigrationServices
	getData() *MigrationData
	up() error
	down() error
	allowUp() error
	runMigration([]byte, string) ([]byte, error)
}

// MigrationServices contains references to services used by the migrations.
type MigrationServices struct {
	fs        afero.Fs
	fsManager *fsm.SveltinFSManager
	pathMaker *pathmaker.SveltinPathMaker
	logger    *log.Logger
}

// NewMigrationServices creates an instance of MigrationService struct.
func NewMigrationServices(fs afero.Fs, fsm *fsm.SveltinFSManager, pathmaker *pathmaker.SveltinPathMaker, logger *log.Logger) *MigrationServices {
	return &MigrationServices{
		fs:        fs,
		fsManager: fsm,
		pathMaker: pathmaker,
		logger:    logger,
	}
}

// MigrationData is the struct with data used by migrations.
type MigrationData struct {
	TargetPath           string
	CliVersion           string
	ProjectCliVersion    string
	EnableUpdateNotifier string
	LastCheck            string
}

// MigrationRule is the struct with settings to be matched for running the migration.
type migrationRule struct {
	value           string
	trigger         string
	replaceFullLine bool
	replacerFunc    func(string) string
}

//=============================================================================

func mustMigrate(content []byte, gatekeeper string) bool {
	return !bytes.Contains(content, []byte(gatekeeper))
}

func patternsMatched(content []byte, patterns []string, matcher matcherFunc) bool {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		for _, pattern := range patterns {
			if matcher(content, pattern, line) {
				return true
			}
			continue
		}
	}
	return false
}

func applyMigrationRules(rules []*migrationRule) (string, bool) {
	for _, r := range rules {
		expression := regexp.MustCompile(r.trigger)

		if expression.MatchString(r.value) {
			if r.replaceFullLine {
				return r.replacerFunc(r.value), true
			}
			return expression.ReplaceAllStringFunc(r.value, r.replacerFunc), true
		}
	}
	return "", false
}

func findStringMatcher(content []byte, pattern, line string) bool {
	rule := regexp.MustCompile(pattern)
	matches := rule.FindString(line)
	return len(matches) > 0
}

//=============================================================================

func retrieveFileContent(fs afero.Fs, pathToFile string) ([]byte, error) {
	content, err := afero.ReadFile(fs, pathToFile)
	if err != nil {
		return nil, afero.ErrFileNotFound
	}
	return content, nil
}

func overwriteFile(m IMigration, content []byte) error {
	err := m.getServices().fs.Remove(m.getData().TargetPath)
	if err != nil {
		return err
	}

	if err = afero.WriteFile(m.getServices().fs, m.getData().TargetPath, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func appendToFile(fs afero.Fs, filename string, contentToAppend []string, logger *log.Logger) {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger.Fatalf("failed opening file: %s", err)
	}
	defer f.Close()

	for _, content := range contentToAppend {
		if _, err = f.WriteString(content); err != nil {
			logger.Fatalf("failed writing to file: %s", err)
		}
	}
}

func getTextInBetween(text string, start string, end string) string {
	startIndex := strings.Index(text, start)
	if startIndex == -1 {
		return ""
	}
	endIndex := strings.Index(text, end) + len(end)
	if endIndex == -1 {
		return ""
	}
	return text[startIndex:endIndex]
}

func replaceTextInBetween(old, new, start, end string) string {
	textInBetween := getTextInBetween(old, start, end)
	return strings.ReplaceAll(old, textInBetween, new)
}
