package graph

import (
	"fio-service/iface"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

type GqlHandler struct {
	handler    *handler.Server
	playground *http.HandlerFunc
}

func NewGqlHandler(service iface.PersonService, logger iface.Ilogger) *GqlHandler {
	handler := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{service, logger}}))
	h := playground.Handler("GraphQL", "/query")
	return &GqlHandler{handler, &h}
}

func (gh *GqlHandler) Query(ctx *gin.Context) {
	gh.handler.ServeHTTP(ctx.Writer, ctx.Request)
}

func (gh *GqlHandler) Playground(ctx *gin.Context) {
	gh.playground.ServeHTTP(ctx.Writer, ctx.Request)
}
