package add

import (
	"cvgo/console/internal/console"
	"cvgo/console/internal/consolepath"
	"cvgo/kit/filekit"
	"cvgo/provider"
	"cvgo/provider/clog"
	"cvgo/provider/config"
	"path/filepath"
)

var httpPort string
var routingFile string

// 添加 i18n 支持
func addI18n() {
	routingFile = filepath.Join(pwd, "internal", "routing", "routing.go")
	cfg := provider.Services.NewSingle(config.Name).(config.Service)
	httpPort = cfg.GetHttpPort()
	kv := console.NewKvStorage(filekit.GetParentDir(3))
	webFrameworkKey := "port" + httpPort + "." + "webframework"
	i18nStorageKey := "port" + httpPort + "." + "i18n"
	if val, _ := kv.GetBool(i18nStorageKey); val {
		log.Info("i18n 已经添加过了，无法重复执行。")
		return
	}

	// instance.go 中声明变量
	instanceFile := filepath.Join(filekit.GetParentDir(2), "instance.go")
	content := "\n" + `import "cvgo/provider/i18n"`
	filekit.AddContentUnderLine(instanceFile, "package app", content)

	content = "\n" + `var I18n i18n.Service`
	err := filekit.FileAppendContent(instanceFile, content)
	if err != nil {
		log.Error(err)
	}

	// boot/init.go 中获取实例
	initFile := filepath.Join(pwd, "internal", "boot", "init.go")
	content = "\n" + `
	"cvgo/provider/i18n"
`
	filekit.AddContentUnderLine(initFile, "import (", content)

	content = "    app.I18n = provider.Services.NewSingle(i18n.Name).(i18n.Service)"
	filekit.AddContentUnderLine(initFile, "func init() {", content)

	// 创建语言包 json 文件
	content = `{
  "hello": "你好",
  "world": {
    "china": "中国"
  }
}`
	filekit.FilePutContents(filepath.Join(pwd, "i18n", "zh.json"), content)
	content = `{
  "hello": "hello",
  "world": {
    "china": "china"
  }
}`
	filekit.FilePutContents(filepath.Join(pwd, "i18n", "en.json"), content)

	// 启用中间件
	frameworkType, _ := kv.GetString(webFrameworkKey)
	switch frameworkType {
	case "cvgo":
		useCvgoI18nMiddleware()
	case "fiber":
		useFiberI18nMiddleware()
	}

	// 标识已添加 i18n
	kv.Set(i18nStorageKey, true)
	clog.GreenPrintln("添加 i18n 成功")
}

// cvgo 框架启用 i18n 中间件
func useCvgoI18nMiddleware() {
	content := `    "cvgo/provider/httpserver/middleware"
`
	err := filekit.AddContentUnderLine(routingFile, "import (", content)
	if err != nil {
		log.Error("修改 routing.go 失败", err)
		return
	}
	content = `
    engine.UseMiddleware(middleware.I18n())
`
	filekit.AddContentUnderLine(routingFile, "func Routes(engine *httpserver.Engine) {", content)
}

// fiber 框架添加 i18n 中间件
func useFiberI18nMiddleware() {
	// 拷贝自定义中间件模板
	src := filepath.Join(consolepath.FiberTplForModule(), "middleware", "i18n.go")
	dest := filepath.Join(pwd, "internal", "middleware", "i18n.go")
	filekit.CopyFile(src, dest)
	// 在路由中启用中间件
	content := `    app.Use(middleware.I18n)`
	filekit.AddContentUnderLine(routingFile, "func Routes(app *fiber.App) {", content)
}
