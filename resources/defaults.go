/*
Copyright © 2021 Mirco Veltri <github@mircoveltri.me>

Use of this source code is governed by Apache 2.0 license
that can be found in the LICENSE file.
*/
package resources

import "embed"

const sveltinAsciiArt = `
                _ _   _
               | | | (_)
  _____   _____| | |_ _ _ __
 / __\ \ / / _ \ | __| | '_ \
 \__ \\ V /  __/ | |_| | | | |
 |___/ \_/ \___|_|\__|_|_| |_|

`

func GetAsciiArt() string {
	return sveltinAsciiArt
}

//go:embed internal/templates/*
var SveltinFS embed.FS

var SveltinProjectFS = map[string]string{
	"defaults":     "internal/templates/site/defaults.js.gotxt",
	"externals":    "internal/templates/site/externals.js.gotxt",
	"website":      "internal/templates/site/website.js.gotxt",
	"init_menu":    "internal/templates/site/init_menu.js.gotxt",
	"menu":         "internal/templates/site/menu.js.gotxt",
	"dotenv":       "internal/templates/misc/env.gotxt",
	"readme":       "internal/templates/misc/README.md.gotxt",
	"license":      "internal/templates/misc/LICENSE.gotxt",
	"theme_config": "internal/templates/themes/theme.config.js.gotxt",
}

var SveltinResourceFS = map[string]string{
	"api":   "internal/templates/resource/api.gotxt",
	"lib":   "internal/templates/resource/lib.gotxt",
	"index": "internal/templates/resource/index.gotxt",
	"slug":  "internal/templates/resource/slug.gotxt",
}

var SveltinMetadataFS = map[string]string{
	"api_single": "internal/templates/resource/metadata/apiSingle.gotxt",
	"api_list":   "internal/templates/resource/metadata/apiList.gotxt",
	"lib_single": "internal/templates/resource/metadata/libSingle.gotxt",
	"lib_list":   "internal/templates/resource/metadata/libList.gotxt",
	"index":      "internal/templates/resource/metadata/index.gotxt",
	"slug":       "internal/templates/resource/metadata/slug.gotxt",
}

var SveltinPageFS = map[string]string{
	"svelte":   "internal/templates/page/page.svelte.gotxt",
	"markdown": "internal/templates/page/page.svx.gotxt",
}

var SveltinContentFS = map[string]string{
	"blank":  "internal/templates/content/blank.svx.gotxt",
	"sample": "internal/templates/content/sample.svx.gotxt",
}

var SveltinXMLFS = map[string]string{
	"sitemap_static": "internal/templates/xml/sitemap.xml.gotxt",
	"rss_static":     "internal/templates/xml/rss.xml.gotxt",
	"sitemap_ssr":    "internal/templates/xml/ssr_sitemap.xml.ts.gotxt",
	"rss_ssr":        "internal/templates/xml/ssr_rss.xml.ts.gotxt",
}

var SveltinVanillaCSSThemeFS = map[string]string{
	"package_json": "internal/templates/themes/vanillacss/package.json.gotxt",
	"app_css":      "internal/templates/themes/vanillacss/app.css",
}

var SveltinTailwindCSSThemeFS = map[string]string{
	"package_json":        "internal/templates/themes/tailwindcss/package.json.gotxt",
	"tailwind_css_config": "internal/templates/themes/tailwindcss/tailwind.config.cjs",
	"app_css":             "internal/templates/themes/tailwindcss/app.css",
	"postcss":             "internal/templates/themes/tailwindcss/postcss.config.cjs",
}
