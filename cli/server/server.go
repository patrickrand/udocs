package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/ultimatesoftware/udocs/cli/config"
	"github.com/ultimatesoftware/udocs/cli/storage"
	"github.com/ultimatesoftware/udocs/cli/udocs"
	"golang.org/x/net/context"
)

type Server struct {
	treeMux    *httptreemux.TreeMux
	fileServer http.Handler
	settings   config.Settings
	dao        storage.Dao
	tmpl       *udocs.Template
	scheme     string
	host       string
}

var BaseDirs = []string{
	udocs.ArchivePath(),
	udocs.BuildPath(),
	udocs.DeployPath(),
}

func New(settings *config.Settings, dao storage.Dao) *Server {
	if err := createBaseDirs(); err != nil {
		log.Fatalf("server.New: failed to create base directories: %v", err)
	}

	tmpl := udocs.MustParseTemplate(defaultTemplateParams(*settings), udocs.DefaultTemplateFiles()...)

	if err := os.Symlink(udocs.StaticPath(), filepath.Join(udocs.DeployPath(), "static")); err != nil && os.IsNotExist(err) {
		log.Fatalf("server.New: failed to symlink static directory: %v", err)
	}

	scheme, host := parseHostURL(settings.EntryPoint)

	if settings.RootRoute == "" {
		settings.RootRoute = "index.html"
	}

	s := &Server{
		treeMux:    httptreemux.New(),
		fileServer: http.FileServer(http.Dir(udocs.DeployPath())),
		settings:   *settings,
		dao:        dao,
		tmpl:       tmpl,
		scheme:     scheme,
		host:       host,
	}

	s.registerEndpoints()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// use file-server to serve static pages (fonts, stylesheets, scripts, etc.)
	if strings.HasPrefix(r.URL.Path, "/static") {
		s.fileServer.ServeHTTP(w, r)
		return
	}

	// otherwise, use default router
	s.treeMux.ServeHTTP(w, r)
}

func (s *Server) Handle(method, path string, h ContextHandlerFunc) {
	s.treeMux.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)

		ctx := context.Background()
		if params != nil {
			for k, v := range params {
				ctx = context.WithValue(ctx, k, v)
			}
		}

		h(ctx, w, r)
	})
}

func (s *Server) registerEndpoints() {
	s.Handle(http.MethodGet, "/", s.reverseProxyHandler)
	s.Handle(http.MethodGet, "/:route", s.pageHandler)
	s.Handle(http.MethodGet, "/:route/*", s.pageHandler)
	s.Handle(http.MethodPost, "/api/:route", s.updateHandler)
	s.Handle(http.MethodDelete, "/api/:route", s.destroyHandler)
	s.Handle(http.MethodPost, "/api/:project/:repo", s.repoHandler)
	s.Handle(http.MethodGet, "/search", s.searchHandler)
}

func (s *Server) Seed() error {
	for _, resource := range s.settings.Seed {
		tokens := strings.Split(resource, "/")
		if len(tokens) != 2 {
			return fmt.Errorf("server.seedDocs: failed to parse seed %q", resource)
		}
		if err := udocs.GitArchive(tokens[0], tokens[1], s.dao); err != nil {
			log.Printf("server.seedDocs: failed to pull %q: %v\nMake sure that the repository exists and that you added the ssh key to it.\n", resource, err)
		}
	}

	return nil
}

func parseHostURL(url string) (string, string) {
	if host := strings.TrimPrefix(url, "http://"); len(host) != len(url) {
		return "http", host
	}
	if host := strings.TrimPrefix(url, "https://"); len(host) != len(url) {
		return "https", host
	}

	return "", url
}

func updateSidebar(sidebar []udocs.Summary, summary udocs.Summary, dao storage.Dao) error {
	var found bool
	for i, item := range sidebar {
		if item.Route == summary.Route {
			sidebar[i] = summary
			found = true
		}
	}

	// when running udocs-serve locally, we may have arbitrary routes that are not predefined
	if !found {
		sidebar = append(sidebar, summary)
	}

	globalSummaryData, err := json.Marshal(sidebar)
	if err != nil {
		return fmt.Errorf("api.UpdateSummary failed to marshal global SUMMARY.json: %v", err)
	}

	if err := dao.Insert(udocs.SIDEBAR_JSON, globalSummaryData); err != nil {
		return fmt.Errorf("api.UpdateSummary failed to insert global SUMMARY.json: %v", err)
	}

	return nil
}

func createBaseDirs() error {
	for _, dir := range BaseDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("server.createBaseDirs failed creating required base directory: %v\n", err)
		}
	}
	return nil
}

func defaultTemplateParams(settings config.Settings) map[string]interface{} {
	m := make(map[string]interface{})
	m["entrypoint"] = settings.EntryPoint
	m["organization"] = settings.Organization
	m["email"] = settings.Email
	m["search_placeholder"] = settings.SearchPlaceholder
	m["logo"] = settings.LogoURL
	return m
}
