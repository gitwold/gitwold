package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/magomedcoder/repo/internal/config"
	deliveryhttp "github.com/magomedcoder/repo/internal/delivery/http"
	"github.com/magomedcoder/repo/internal/delivery/http/handler"
	deliveryssh "github.com/magomedcoder/repo/internal/delivery/ssh"
	"github.com/magomedcoder/repo/internal/infrastructure/git"
	"github.com/magomedcoder/repo/internal/infrastructure/persistence/sqlite"
	"github.com/magomedcoder/repo/internal/usecase"
	"github.com/magomedcoder/repo/pkg/bcrypt"
	"github.com/magomedcoder/repo/pkg/token"
)

func main() {
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatalf("git binary required for Smart HTTP: %v", err)
	}

	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := sqlite.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
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

	authUC := usecase.NewAuthUseCase(userStore, sessionStore, orgStore, hasher, tokens)
	folderUC := usecase.NewFolderUseCase(folderStore, repoStore, orgStore)
	repoUC := usecase.NewRepositoryUseCase(repoStore, folderStore, userStore, orgStore, gitRepo, activityStore)
	tokenUC := usecase.NewTokenUseCase(tokenStore, tokens)
	sshKeyUC := usecase.NewSSHKeyUseCase(sshKeyStore, userStore)
	gitUC := usecase.NewGitUseCase(repoStore, folderStore, userStore, orgStore, tokenStore, hasher, gitRepo, activityStore)
	issueUC := usecase.NewIssueUseCase(issueStore, repoStore, folderStore, userStore, orgStore, activityStore)
	pullUC := usecase.NewPullRequestUseCase(pullStore, repoStore, folderStore, userStore, orgStore, gitRepo, activityStore)
	profileUC := usecase.NewProfileUseCase(userStore, cfg.AvatarsDir)
	orgUC := usecase.NewOrganizationUseCase(orgStore, userStore, repoStore, folderStore, activityStore)
	activityUC := usecase.NewActivityUseCase(activityStore, userStore, orgStore, repoStore)

	authHandler := handler.NewAuthHandler(authUC)
	folderHandler := handler.NewFolderHandler(folderUC)
	repoHandler := handler.NewRepositoryHandler(repoUC, cfg.ReposRoot)
	tokenHandler := handler.NewTokenHandler(tokenUC)
	sshKeyHandler := handler.NewSSHKeyHandler(sshKeyUC)
	gitHandler := handler.NewGitHandler(gitUC)
	issueHandler := handler.NewIssueHandler(issueUC)
	pullHandler := handler.NewPullRequestHandler(pullUC)
	profileHandler := handler.NewProfileHandler(profileUC)
	orgHandler := handler.NewOrganizationHandler(orgUC)
	activityHandler := handler.NewActivityHandler(activityUC)

	router := deliveryhttp.NewRouter(authHandler, repoHandler, folderHandler, tokenHandler, sshKeyHandler, gitHandler, issueHandler, pullHandler, profileHandler, orgHandler, activityHandler, authUC)

	sshServer := deliveryssh.NewServer(sshKeyUC, gitUC, cfg.SSHDir, cfg.SSHAddr)
	go func() {
		if err := sshServer.ListenAndServe(); err != nil {
			log.Fatalf("ssh: %v", err)
		}
	}()

	log.Printf("listening on %s (ssh %s)", cfg.HTTPAddr, cfg.SSHAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}
