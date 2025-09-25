package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Data struct {
	FeatName string   // input full feature (ex: "master/area")
	Parts    []string // split parts ["master","area"]
	Entity   string   // last ("area")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: make gen feat=<feature>  (ex: customer | master/area)")
	}

	feat := os.Args[1]
	parts := strings.Split(feat, "/")
	entity := parts[len(parts)-1]

	d := Data{
		FeatName: feat,
		Parts:    parts,
		Entity:   entity,
	}

	// daftar template yang mau diproses
	files := []struct {
		TplPath   string
		OutDir    string
		OutSuffix string
		TplName   string
	}{
		{
			TplPath:   "tools/templates/model.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "models"),
			OutSuffix: ".model.go",
			TplName:   "model",
		},
		{
			TplPath:   "tools/templates/validation.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "validations"),
			OutSuffix: ".validation.go",
			TplName:   "validation",
		},
		{
			TplPath:   "tools/templates/service.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "services"),
			OutSuffix: ".service.go",
			TplName:   "service",
		},
		{
			TplPath:   "tools/templates/controller.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "controllers"),
			OutSuffix: ".controller.go",
			TplName:   "controller",
		},
		{
			TplPath:   "tools/templates/repository.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "repositories"),
			OutSuffix: ".repository.go",
			TplName:   "repository",
		},
		{
			TplPath:   "tools/templates/dto.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s", "dto"),
			OutSuffix: ".dto.go",
			TplName:   "dto",
		},
		{
			TplPath:   "tools/templates/route.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s"),
			OutSuffix: "",
			TplName:   "route",
		},
		{
			TplPath:   "tools/templates/module.tmpl",
			OutDir:    filepath.Join("internal", "modules", toKebabPath(d.Parts[:len(d.Parts)-1]), toKebab(d.Entity)+"s"),
			OutSuffix: "",
			TplName:   "module",
		},
	}

	for _, file := range files {
		// pastikan template ketemu
		if _, err := os.Stat(file.TplPath); err != nil {
			log.Fatalf("template not found at %s: %v", file.TplPath, err)
		}

		// parse template
		tpl := template.Must(
			template.New(file.TplName).
				Funcs(template.FuncMap{
					"Pascal": toPascalCase,
					"Camel":  toCamelCase,
					"Plural": toPlural,
					"Kebab":  toKebab,
				}).
				ParseFiles(file.TplPath),
		)

		// pastikan folder ada
		if err := os.MkdirAll(file.OutDir, 0o755); err != nil {
			log.Fatalf("make dir: %v", err)
		}

		// nama file output
		var outFile string
		switch file.TplName {
		case "route":
			outFile = filepath.Join(file.OutDir, "route.go")
		case "module":
			outFile = filepath.Join(file.OutDir, "module.go")
		default:
			outFile = filepath.Join(file.OutDir, strings.ToLower(d.Entity)+file.OutSuffix)
		}

		// hindari overwrite
		if _, err := os.Stat(outFile); err == nil {
			log.Fatalf("file already exists: %s", outFile)
		}

		f, err := os.Create(outFile)
		if err != nil {
			log.Fatalf("create file: %v", err)
		}
		defer f.Close()

		if err := tpl.ExecuteTemplate(f, file.TplName, d); err != nil {
			log.Fatal(err)
		}

		log.Println("Generated:", outFile)
	}

	updateMainRoute(d)
}

func updateMainRoute(d Data) {
	routeFile := "internal/route/route.go"
	content, err := os.ReadFile(routeFile)
	if err != nil {
		log.Printf("skip update route.go: %v", err)
		return
	}

	// entity & path
	modPath := filepath.Join(append(toCamelParts(d.Parts[:len(d.Parts)-1]), toCamelCase(d.Entity)+"s")...)
	modName := toCamelCase(d.Entity) + "s"
	pkgName := toPascalCase(d.Entity) + "Module"

	// Inject import
	importLine := fmt.Sprintf("\t%[1]s \"%s/internal/modules/%s\"", modName, "github.com/hafizhproject45/Golang-Boilerplate.git", modPath)
	if !strings.Contains(string(content), importLine) {
		content = []byte(strings.Replace(string(content),
			"// MODULE IMPORTS",
			importLine+"\n\t// MODULE IMPORTS",
			1))
	}

	// Inject registry
	registryLine := fmt.Sprintf("\t\t%[1]s.%[2]s{},", modName, pkgName)
	if !strings.Contains(string(content), registryLine) {
		content = []byte(strings.Replace(string(content),
			"// MODULE REGISTRY",
			registryLine+"\n\t\t// MODULE REGISTRY",
			1))
	}

	if err := os.WriteFile(routeFile, content, 0644); err != nil {
		log.Fatal(err)
	}
	log.Println("Updated:", routeFile)
}

func toPascalCase(s string) string {
	sep := func(r rune) bool { return r == '_' || r == '-' || r == ' ' || r == '/' }
	parts := strings.FieldsFunc(s, sep)
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}

func toCamelCase(s string) string {
	p := toPascalCase(s)
	if p == "" {
		return ""
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// simple pluralizer (cukup untuk kasus umum: tambah 's')
func toPlural(s string) string {
	s = strings.ToLower(s)
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if !(prev == 'a' || prev == 'i' || prev == 'u' || prev == 'e' || prev == 'o') {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

// kebab-case (untuk folder)
func toKebab(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	out := b.String()
	out = strings.ReplaceAll(out, "--", "-")
	return strings.Trim(out, "-")
}

// join multiple parts jadi kebab path
func toKebabPath(parts []string) string {
	return filepath.Join(toKebabParts(parts)...)
}

func toKebabParts(parts []string) []string {
	var out []string
	for _, p := range parts {
		out = append(out, toKebab(p))
	}
	return out
}

// join multiple parts jadi camelCase path
// func toCamelPath(parts []string) string {
// 	return filepath.Join(toCamelParts(parts)...)
// }

func toCamelParts(parts []string) []string {
	var out []string
	for i, p := range parts {
		if i == 0 {
			// part pertama lower-case semua
			out = append(out, toCamelCase(p))
		} else {
			// part berikutnya PascalCase biar tetap nyambung camel
			out = append(out, toPascalCase(p))
		}
	}
	return out
}
