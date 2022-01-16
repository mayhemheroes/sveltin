/*
Copyright © 2021 Mirco Veltri <github@mircoveltri.me>

Use of this source code is governed by Apache 2.0 license
that can be found in the LICENSE file.
*/
package fsm

import (
	"path/filepath"
	"strings"

	"github.com/sveltinio/sveltin/config"
	"github.com/sveltinio/sveltin/helpers"
	"github.com/sveltinio/sveltin/sveltinlib/composer"
	"github.com/sveltinio/sveltin/sveltinlib/pathmaker"
)

type SveltinFSManager struct {
	maker *pathmaker.SveltinPathMaker
}

func NewSveltinFSManager(maker *pathmaker.SveltinPathMaker) *SveltinFSManager {
	return &SveltinFSManager{
		maker: maker,
	}
}

func (s *SveltinFSManager) GetFolder(name string) *composer.Folder {
	switch name {
	case "root":
		return composer.GetRootFolder(s.maker)
	case "config":
		return composer.GetConfigFolder(s.maker)
	case "content":
		return composer.GetContentFolder(s.maker)
	case "routes":
		return composer.GetRoutesFolder(s.maker)
	case "api":
		return composer.GetAPIFolder(s.maker)
	case "lib":
		return composer.GetLibFolder(s.maker)
	case "static":
		return composer.GetStaticFolder(s.maker)
	case "themes":
		return composer.GetThemesFolder(s.maker)
	default:
		return composer.NewFolder(name)
	}
}

func (s *SveltinFSManager) NewResourceContentFolder(name string, resource string) *composer.Folder {
	return composer.NewFolder(filepath.Join(resource, name))
}

func (s *SveltinFSManager) NewResourceContentFile(name string, template string) *composer.File {
	return &composer.File{
		Name:       s.maker.GetResourceContentFilename(),
		TemplateId: template,
		TemplateData: &config.TemplateData{
			Name: name,
		},
	}
}

func (s *SveltinFSManager) NewPublicPage(name string, language string) *composer.File {
	return &composer.File{
		Name:       helpers.PublicPageFilename(name, language),
		TemplateId: language,
		TemplateData: &config.TemplateData{
			Name: name,
		},
	}
}

func (s *SveltinFSManager) NewNoPage(name string, projectConfig *config.SiteConfig, resources []string, contents map[string][]string, metadata map[string][]string, pages []string) *composer.File {
	return &composer.File{
		Name:       name + ".xml",
		TemplateId: name,
		TemplateData: &config.TemplateData{
			NoPage: &config.NoPage{
				Config: projectConfig,
				Items:  helpers.NewNoPageItems(resources, contents, metadata, pages),
			},
		},
	}
}

func (s *SveltinFSManager) NewConfigFile(projectName string, name string, cliVersion string) *composer.File {
	filename := strings.ToLower(name) + ".js.ts"
	return &composer.File{
		Name:       filename,
		TemplateId: name,
		TemplateData: &config.TemplateData{
			ProjectName: projectName,
			Name:        filename,
			Misc:        cliVersion,
		},
	}
}

func (s *SveltinFSManager) NewDotEnvFile(projectName string, name string) *composer.File {
	return &composer.File{
		Name:       name,
		TemplateId: "dotenv",
		TemplateData: &config.TemplateData{
			Name: name,
			Misc: "http://localhost:3000",
		},
	}
}

func (s *SveltinFSManager) NewContentFile(name string, template string, resource string) *composer.File {
	return &composer.File{
		Name:       s.maker.GetResourceContentFilename(),
		TemplateId: template,
		TemplateData: &config.TemplateData{
			Name: name,
		},
	}
}
