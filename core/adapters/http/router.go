package httpadapter

import (
	"net/http"

	assetshandler "picup/core/adapters/handlers/assets"
	authhandler "picup/core/adapters/handlers/auth"
	confighandler "picup/core/adapters/handlers/config"
	imagehandler "picup/core/adapters/handlers/image"
	webuihandler "picup/core/adapters/handlers/webui"
	"picup/core/domain/repositories"
	"picup/core/usecases/workflow"
)

type RouterDependencies struct {
	Repository  repositories.Repository
	FileStorage workflow.FileStorage
}

func NewRouter(deps RouterDependencies) http.Handler {
	mux := http.NewServeMux()
	authMiddleware := authhandler.NewAuthMiddleware(deps.Repository)

	webuihandler.NewWebUIHandler().RouteInit(mux, defaultHandler)
	confighandler.NewConfigHandler().RouteInit(mux, defaultHandler)
	authhandler.NewAuthHandler(deps.Repository).RouteInit(mux, defaultHandler)
	imagehandler.NewImageHandler(deps.Repository, deps.FileStorage).RouteInit(mux, authMiddleware.RequireAuth)
	assetshandler.NewAssetsHandler().RouteInit(mux, authMiddleware.RequireAuth)

	return CORSMiddleware(LoggingMiddleware(mux))
}
