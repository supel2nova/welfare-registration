package domain_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNoFrameworkImports(t *testing.T) {
	direct := []string{
		"github.com/gin-gonic/gin",
		"github.com/jackc/pgx",
		"net/http",
		"database/sql",
	}
	out, err := exec.Command("go", "list", "-f", `{{.ImportPath}}|{{join .Imports ","}}`, "./...").Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		pkg, imports, _ := strings.Cut(line, "|")
		for _, imp := range strings.Split(imports, ",") {
			for _, b := range direct {
				if imp == b || strings.HasPrefix(imp, b+"/") {
					t.Errorf("%s import %q ตรงๆ", pkg, imp)
				}
			}
		}
	}

	deps, err := exec.Command("go", "list", "-deps", "./...").Output()
	if err != nil {
		t.Fatalf("go list -deps: %v", err)
	}
	for _, dep := range strings.Fields(string(deps)) {
		for _, b := range []string{"github.com/gin-gonic/gin", "github.com/jackc/pgx"} {
			if dep == b || strings.HasPrefix(dep, b+"/") {
				t.Errorf("domain พึ่ง %q ทางอ้อม (%s)", b, dep)
			}
		}
	}
}
