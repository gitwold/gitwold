package http_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	deliveryhttp "github.com/magomedcoder/repo/internal/delivery/http"
	"github.com/magomedcoder/repo/internal/delivery/http/handler"
	"github.com/magomedcoder/repo/internal/infrastructure/git"
	"github.com/magomedcoder/repo/internal/infrastructure/persistence/sqlite"
	"github.com/magomedcoder/repo/internal/usecase"
	"github.com/magomedcoder/repo/pkg/bcrypt"
	"github.com/magomedcoder/repo/pkg/token"
)

func newTestServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o755); err != nil {
		t.Fatal(err)
	}

	db, err := sqlite.NewDB(filepath.Join(dir, "data", "db.sqlite"))
	if err != nil {
		t.Fatal(err)
	}

	userStore := sqlite.NewUserStore(db)
	sessionStore := sqlite.NewSessionStore(db)
	folderStore := sqlite.NewFolderStore(db)
	repoStore := sqlite.NewRepositoryStore(db)
	tokenStore := sqlite.NewAccessTokenStore(db)
	sshKeyStore := sqlite.NewSSHKeyStore(db)
	issueStore := sqlite.NewIssueStore(db)
	pullStore := sqlite.NewPullRequestStore(db)
	orgStore := sqlite.NewOrganizationStore(db)
	activityStore := sqlite.NewActivityStore(db)
	gitRepo := git.NewRepository()
	hasher := bcrypt.NewHasher()
	tokens := token.NewGenerator()
	avatarsDir := filepath.Join(dir, "data", "avatars")
	reposRoot := filepath.Join(dir, "data", "repos")

	authUC := usecase.NewAuthUseCase(userStore, sessionStore, orgStore, hasher, tokens)
	folderUC := usecase.NewFolderUseCase(folderStore, repoStore, orgStore)
	repoUC := usecase.NewRepositoryUseCase(repoStore, folderStore, userStore, orgStore, gitRepo, activityStore)
	tokenUC := usecase.NewTokenUseCase(tokenStore, tokens)
	sshKeyUC := usecase.NewSSHKeyUseCase(sshKeyStore, userStore)
	gitUC := usecase.NewGitUseCase(repoStore, folderStore, userStore, orgStore, tokenStore, hasher, gitRepo, activityStore)
	issueUC := usecase.NewIssueUseCase(issueStore, repoStore, folderStore, userStore, orgStore, activityStore)
	pullUC := usecase.NewPullRequestUseCase(pullStore, repoStore, folderStore, userStore, orgStore, gitRepo, activityStore)
	profileUC := usecase.NewProfileUseCase(userStore, avatarsDir)
	orgUC := usecase.NewOrganizationUseCase(orgStore, userStore, repoStore, folderStore, activityStore)
	activityUC := usecase.NewActivityUseCase(activityStore, userStore, orgStore, repoStore)

	router := deliveryhttp.NewRouter(
		handler.NewAuthHandler(authUC),
		handler.NewRepositoryHandler(repoUC, reposRoot),
		handler.NewFolderHandler(folderUC),
		handler.NewTokenHandler(tokenUC),
		handler.NewSSHKeyHandler(sshKeyUC),
		handler.NewGitHandler(gitUC),
		handler.NewIssueHandler(issueUC),
		handler.NewPullRequestHandler(pullUC),
		handler.NewProfileHandler(profileUC),
		handler.NewOrganizationHandler(orgUC),
		handler.NewActivityHandler(activityUC),
		authUC,
	)
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return srv, dir
}

type jsonResp struct {
	StatusCode int
	Body       string
	Cookies    []*http.Cookie
}

func doJSON(t *testing.T, srv *httptest.Server, method, path, body string, cookies []*http.Cookie) jsonResp {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = bytes.NewReader([]byte(body))
	}

	req, err := http.NewRequest(method, srv.URL+path, reader)
	if err != nil {
		t.Fatal(err)
	}

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	return jsonResp{
		StatusCode: res.StatusCode,
		Body:       string(raw),
		Cookies:    res.Cookies(),
	}
}
