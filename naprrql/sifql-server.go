// sifql-server.go
//
// simple web server to support gql queries & web ui (graphiql)
//
package naprrql

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/playlyfe/go-graphql"
)

var executor *graphql.Executor

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// the core graphql handler routine
//
func graphQLHandler(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}
	// log.Printf("variables: %v\n\n", variables)
	result, err := executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)

}

//
// launches the server
//
func RunQLServer() {

	executor = buildExecutor()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Static("/", "public")
	e.File("/sifql", "public/index.html")

	e.POST("/graphql", graphQLHandler)

	e.Logger.Fatal(e.Start(":1329"))
}
