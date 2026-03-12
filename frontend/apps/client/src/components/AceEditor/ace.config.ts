import ace from 'ace-builds'
import 'ace-builds/src-noconflict/ext-language_tools' // 启用语言工具扩展

// 设置模式的 URL
// JSON
// import modeJsonUrl from 'ace-builds/src-noconflict/mode-json?url'
// ace.config.setModuleUrl('ace/mode/json', modeJsonUrl)

import modeJavascriptUrl from 'ace-builds/src-noconflict/mode-javascript?url'
ace.config.setModuleUrl('ace/mode/javascript', modeJavascriptUrl)

// HTML
// import modeHtmlUrl from 'ace-builds/src-noconflict/mode-html?url'
// ace.config.setModuleUrl('ace/mode/html', modeHtmlUrl)

import modeCppUrl from 'ace-builds/src-noconflict/mode-c_cpp?url'
ace.config.setModuleUrl('ace/mode/c_cpp', modeCppUrl)

// YAML
// import modeYamlUrl from 'ace-builds/src-noconflict/mode-yaml?url'
// ace.config.setModuleUrl('ace/mode/yaml', modeYamlUrl)

// SQL
// import modeSqlUrl from 'ace-builds/src-noconflict/mode-sql?url'
// ace.config.setModuleUrl('ace/mode/sql', modeSqlUrl)

// XML
// import modeXmlUrl from 'ace-builds/src-noconflict/mode-xml?url'
// ace.config.setModuleUrl('ace/mode/xml', modeXmlUrl)

// MARKDOWN
// import modeMarkdownUrl from 'ace-builds/src-noconflict/mode-markdown?url'
// ace.config.setModuleUrl('ace/mode/markdown', modeMarkdownUrl)

import modeTextUrl from 'ace-builds/src-noconflict/mode-text?url'
ace.config.setModuleUrl('ace/mode/text', modeTextUrl)

import modePythonUrl from 'ace-builds/src-noconflict/mode-python?url'
ace.config.setModuleUrl('ace/mode/python', modePythonUrl)

import modeJavaUrl from 'ace-builds/src-noconflict/mode-java?url'
ace.config.setModuleUrl('ace/mode/java', modeJavaUrl)

import modeCsharpUrl from 'ace-builds/src-noconflict/mode-csharp?url'
ace.config.setModuleUrl('ace/mode/csharp', modeCsharpUrl)

// import modeGoUrl from 'ace-builds/src-noconflict/mode-go?url'
// ace.config.setModuleUrl('ace/mode/go', modeGoUrl)

import modeRustUrl from 'ace-builds/src-noconflict/mode-rust?url'
ace.config.setModuleUrl('ace/mode/rust', modeRustUrl)

// 设置主题的 URL
import themeGithubUrl from 'ace-builds/src-noconflict/theme-github?url'
ace.config.setModuleUrl('ace/theme/github', themeGithubUrl)

import themeChromeUrl from 'ace-builds/src-noconflict/theme-chrome?url'
ace.config.setModuleUrl('ace/theme/chrome', themeChromeUrl)

import themeMonokaiUrl from 'ace-builds/src-noconflict/theme-monokai?url'
ace.config.setModuleUrl('ace/theme/monokai', themeMonokaiUrl)

import themeXcodeUrl from 'ace-builds/src-noconflict/theme-xcode?url'
ace.config.setModuleUrl('ace/theme/xcode', themeXcodeUrl)

import themeDraculaUrl from 'ace-builds/src-noconflict/theme-dracula?url'
ace.config.setModuleUrl('ace/theme/dracula', themeDraculaUrl)

import themeCloudsUrl from 'ace-builds/src-noconflict/theme-clouds?url'
ace.config.setModuleUrl('ace/theme/clouds', themeCloudsUrl)

import themeTerminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', themeTerminalUrl)
// dawn	import themeDawnUrl from 'ace-builds/src-noconflict/theme-dawn?url'
// eclipse	import themeEclipseUrl from 'ace-builds/src-noconflict/theme-eclipse?url'
// iplastic	import themeIplasticUrl from 'ace-builds/src-noconflict/theme-iplastic?url'
// solarized_light	import themeSolarizedLightUrl from 'ace-builds/src-noconflict/theme-solarized_light?url'
// textmate	import themeTextmateUrl from 'ace-builds/src-noconflict/theme-textmate?url'
// tomorrow	import themeTomorrowUrl from 'ace-builds/src-noconflict/theme-tomorrow?url'
// monokai	import themeMonokaiUrl from 'ace-builds/src-noconflict/theme-monokai?url'
// ambiance	import themeAmbianceUrl from 'ace-builds/src-noconflict/theme-ambiance?url'
// chaos	import themeChaosUrl from 'ace-builds/src-noconflict/theme-chaos?url'
// clouds_midnight	import themeCloudsMidnightUrl from 'ace-builds/src-noconflict/theme-clouds_midnight?url'
// cobalt	import themeCobaltUrl from 'ace-builds/src-noconflict/theme-cobalt?url'
// dracula	import themeDraculaUrl from 'ace-builds/src-noconflict/theme-dracula?url'
// gob	import themeGobUrl from 'ace-builds/src-noconflict/theme-gob?url'
// gruvbox	import themeGruvboxUrl from 'ace-builds/src-noconflict/theme-gruvbox?url'
// idle_fingers	import themeIdleFingersUrl from 'ace-builds/src-noconflict/theme-idle_fingers?url'
// kr_theme	import themeKrThemeUrl from 'ace-builds/src-noconflict/theme-kr_theme?url'
// merbivore	import themeMerbivoreUrl from 'ace-builds/src-noconflict/theme-merbivore?url'
// merbivore_soft	import themeMerbivoreSoftUrl from 'ace-builds/src-noconflict/theme-merbivore_soft?url'
// mono_industrial	import themeMonoIndustrialUrl from 'ace-builds/src-noconflict/theme-mono_industrial?url'
// nord_dark	import themeNordDarkUrl from 'ace-builds/src-noconflict/theme-nord_dark?url'
// pastel_on_dark	import themePastelOnDarkUrl from 'ace-builds/src-noconflict/theme-pastel_on_dark?url'
// solarized_dark	import themeSolarizedDarkUrl from 'ace-builds/src-noconflict/theme-solarized_dark?url'
// terminal	import themeTerminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
// tomorrow_night	import themeTomorrowNightUrl from 'ace-builds/src-noconflict/theme-tomorrow_night?url'
// tomorrow_night_blue	import themeTomorrowNightBlueUrl from 'ace-builds/src-noconflict/theme-tomorrow_night_blue?url'
// tomorrow_night_bright	import themeTomorrowNightBrightUrl from 'ace-builds/src-noconflict/theme-tomorrow_night_bright?url'
// tomorrow_night_eighties	import themeTomorrowNightEightiesUrl from 'ace-builds/src-noconflict/theme-tomorrow_night_eighties?url'
// twilight	import themeTwilightUrl from 'ace-builds/src-noconflict/theme-twilight?url'
// vibrant_ink	import themeVibrantInkUrl from 'ace-builds/src-noconflict/theme-vibrant_ink?url'

// 设置 worker 的 URL
import workerBaseUrl from 'ace-builds/src-noconflict/worker-base?url'
ace.config.setModuleUrl('ace/worker/base', workerBaseUrl)

import workerJavascriptUrl from 'ace-builds/src-noconflict/worker-javascript?url'
ace.config.setModuleUrl('ace/worker/javascript', workerJavascriptUrl)

// 设置代码片段的 URL
import snippetsHtmlUrl from 'ace-builds/src-noconflict/snippets/html?url'
ace.config.setModuleUrl('ace/snippets/html', snippetsHtmlUrl)

import snippetsJsUrl from 'ace-builds/src-noconflict/snippets/javascript?url'
ace.config.setModuleUrl('ace/snippets/javascript', snippetsJsUrl)

import snippetsSqlUrl from 'ace-builds/src-noconflict/snippets/sql?url'
ace.config.setModuleUrl('ace/snippets/sql', snippetsSqlUrl)

import snippetsCppUrl from 'ace-builds/src-noconflict/snippets/c_cpp?url'
ace.config.setModuleUrl('ace/snippets/c_cpp', snippetsCppUrl)

import snippetsPythonUrl from 'ace-builds/src-noconflict/snippets/python?url'
ace.config.setModuleUrl('ace/snippets/python', snippetsPythonUrl)

import snippetsJavaUrl from 'ace-builds/src-noconflict/snippets/java?url'
ace.config.setModuleUrl('ace/snippets/java', snippetsJavaUrl)

import snippetsCsharpUrl from 'ace-builds/src-noconflict/snippets/csharp?url'
ace.config.setModuleUrl('ace/snippets/csharp', snippetsCsharpUrl)
// 设置语言工具的 URL
import languageToolsUrl from 'ace-builds/src-noconflict/ext-language_tools?url'
ace.config.setModuleUrl('ace/ext/language_tools', languageToolsUrl)
