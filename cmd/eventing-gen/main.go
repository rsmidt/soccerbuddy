package main

import (
	"github.com/rsmidt/soccerbuddy/internal/core"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

var (
	eventTypRe = regexp.MustCompile(`^(.*Event)Type$`)
	eventVerRe = regexp.MustCompile(`^(.*Event)Version$`)
)

type event struct {
	eventmeta

	Name    string
	PkgName string
	PkgPath string
	Agg     string
}

type eventmeta struct {
	Typ     string
	Version string
}

func (em *eventmeta) Isset() bool {
	return em.Typ != "" && em.Version != ""
}

type eventlocator struct {
	path string
	name string
}

func main() {
	config := packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Dir:  path.Join(core.Root, "internal"),
	}
	pkgs, err := packages.Load(&config, "./...")
	if err != nil {
		panic(err)
	}
	if len(pkgs) == 0 {
		log.Println("no packages found")
		return
	}

	var events []event

	var fnameToAggTyp = make(map[string]string)
	var locToEventMeta = make(map[eventlocator]eventmeta)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok {
					return true
				}

				switch decl.Tok {
				case token.CONST:
					for _, spec := range decl.Specs {
						vspec := spec.(*ast.ValueSpec)
						if len(vspec.Values) != 1 {
							continue
						}
						name := vspec.Names[0].Name
						filename := pkg.Fset.Position(file.Package).Filename
						if strings.HasSuffix(name, "AggregateType") {
							typ := vspec.Values[0].(*ast.CallExpr).Args[0].(*ast.BasicLit).Value
							fnameToAggTyp[pkg.PkgPath] = strings.Replace(typ, `"`, "", -1)
							continue
						}
						eventTypRes := eventTypRe.FindStringSubmatch(name)
						if len(eventTypRes) == 2 {
							name := eventTypRes[1]
							// Update or create eventmeta.
							loc := eventlocator{path: filename, name: name}
							em, ok := locToEventMeta[loc]
							if !ok {
								em = eventmeta{}
							}
							em.Typ = strings.Replace(vspec.Values[0].(*ast.CallExpr).Args[0].(*ast.BasicLit).Value, `"`, "", -1)
							locToEventMeta[loc] = em
							continue
						}
						eventVerRes := eventVerRe.FindStringSubmatch(name)
						if len(eventVerRes) == 2 {
							name := eventVerRes[1]
							// Update or create eventmeta.
							loc := eventlocator{path: filename, name: name}
							em, ok := locToEventMeta[loc]
							if !ok {
								em = eventmeta{}
							}
							em.Version = strings.Replace(vspec.Values[0].(*ast.CallExpr).Args[0].(*ast.BasicLit).Value, `"`, "", -1)
							locToEventMeta[loc] = em
							continue
						}
					}
					return false
				case token.TYPE:
					for _, spec := range decl.Specs {
						tspec := spec.(*ast.TypeSpec)
						stype, ok := tspec.Type.(*ast.StructType)
						if !ok {
							continue
						}
						for _, field := range stype.Fields.List {
							starxpr, ok := field.Type.(*ast.StarExpr)
							if !ok {
								continue
							}
							selxpr, ok := starxpr.X.(*ast.SelectorExpr)
							if !ok {
								continue
							}
							if selxpr.Sel.Name == "EventBase" {
								name := tspec.Name.Name
								filename := pkg.Fset.Position(file.Package).Filename
								loc := eventlocator{path: filename, name: name}
								em, ok := locToEventMeta[loc]
								if !ok || !em.Isset() {
									log.Fatalf("eventmeta not found for %s", name)
								}
								agg, ok := fnameToAggTyp[pkg.PkgPath]
								if !ok {
									log.Fatalf("aggregate type not found for %s", filename)
								}
								events = append(events, event{
									eventmeta: em,
									Name:      name,
									PkgName:   pkg.Name,
									PkgPath:   pkg.PkgPath,
									Agg:       agg,
								})
								break
							}
						}
					}
					return false
				}

				return true
			})
		}
	}

	// Generate registry.
	var uniqueImports = make(map[string]struct{})
	for _, e := range events {
		uniqueImports[e.PkgPath] = struct{}{}
	}
	var imports []string
	for imp := range uniqueImports {
		imports = append(imports, imp)
	}

	type registryData struct {
		Imports []string
		Events  []event
	}

	tmpl := template.Must(template.New("registry").Parse(registryTmpl))
	data := registryData{
		Imports: imports,
		Events:  events,
	}
	// Create gen directory if not exists.
	if _, err := os.Stat(path.Join(core.Root, "gen", "eventregistry")); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Join(core.Root, "gen", "eventregistry"), 0755); err != nil {
			log.Fatalf("failed to create directory: %v", err)
		}
	}
	destFile := path.Join(core.Root, "gen", "eventregistry", "registry_gen.go")
	f, err := os.Create(destFile)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	if err := tmpl.Execute(f, &data); err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}
}

const registryTmpl = `package eventregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"time"
	{{range .Imports}}
	"{{.}}"
	{{end}}
)

var Default = &genRegistry{}

type genRegistry struct {}

func (r *genRegistry) MapFrom(
	aggregateID      eventing.AggregateID, 
	aggregateType    eventing.AggregateType, 
	eventVersion     eventing.EventVersion, 
	eventType        eventing.EventType,
	eventID          eventing.EventID,
	aggregateVersion eventing.AggregateVersion,
	journalPosition  eventing.JournalPosition,
	insertedAt       time.Time,
	payload          []byte,
) (*eventing.JournalEvent, error) {
	base := eventing.NewEventBase(aggregateID, aggregateType, eventVersion, eventType)
	combinedID := fmt.Sprintf("%s::%s%s", aggregateType, eventType, eventVersion)

	switch combinedID {
	{{range .Events}}
	case "{{.Agg}}::{{.Typ}}{{.Version}}":
		event := &{{.PkgName}}.{{.Name}}{
			EventBase: base,
		}
		if err := json.Unmarshal(payload, event); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return eventing.NewJournalEvent(event, eventID, aggregateVersion, journalPosition, insertedAt), nil
	{{end}}
	default:
		return nil, errors.New("event not registered")
	}
}
`
