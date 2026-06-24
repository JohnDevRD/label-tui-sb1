package core

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func ListTemplates() ([]string, error) {
	dir, err := TemplatesDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".zpl") {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func LoadTemplate(name string) (string, error) {
	dir, err := TemplatesDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var replacers = map[string]func(Article, int) string{
	"CODE": func(a Article, _ int) string {
		return a.Code
	},
	"DESCRIPTION": func(a Article, _ int) string {
		return a.Description
	},
	"BARCODE": func(a Article, _ int) string {
		return a.Barcode
	},
	"PRICE": func(a Article, _ int) string {
		return fmt.Sprintf("%.2f", a.Price)
	},
	"QTY": func(_ Article, qty int) string {
		return fmt.Sprintf("%d", qty)
	},
}

func RenderTemplate(rawZPL string, article Article, qty int) string {
	result := rawZPL
	for key, fn := range replacers {
		result = strings.ReplaceAll(result, "{{"+key+"}}", fn(article, qty))
	}
	return result
}

func RenderTemplateBatch(rawZPL string, article Article, qty int) string {
	var sb strings.Builder
	for i := 0; i < qty; i++ {
		rendered := RenderTemplate(rawZPL, article, qty)
		if i > 0 {
			rendered = strings.Replace(rendered, "^XA", "^XA\n^PQ1,0,1,Y", 1)
		}
		rendered = strings.Replace(rendered, "^PQ", fmt.Sprintf("^PQ%d", int(math.Max(1, float64(qty)))), 1)
		sb.WriteString(rendered)
	}
	return sb.String()
}
