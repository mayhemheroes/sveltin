# sveltin

## 0.2.10

### Patch Changes

- fix: `generate menu` command used `js` instead of `ts` as file extension causing errors on loading
- fix: with tailwindcss used as css lib, the typography plugin's prose class was not rendered correctly
- fix: postcss and its config file provided for tailwindcss only
- fix: with vanillacss a github markdown theme added as default to render markdown content

## 0.2.9

### Patch Changes

- SvelteKit 1.0.0-next.244 fixed [#3473](https://github.com/sveltejs/kit/issues/3473) and [#3521](https://github.com/sveltejs/kit/pull/3521). `clone()` on fetch response as workaround to avoid '_body used already_' error when building the project removed

## 0.2.8

### Patch Changes

- string utility functions added to get valid page names and contents names
- variable names fixed on page templates

  **Full Changelog**: https://github.com/sveltinio/sveltin/compare/v0.2.7...v0.2.8

## 0.2.7

### Patch Changes

- fix: image path on seo components
- fix: seo components added to pages