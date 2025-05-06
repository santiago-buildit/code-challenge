// @title Library Management API
// @version 1.0
// @description This is a code-challenge API to manage books in a library.
// @host d21meifd8clvjr.cloudfront.net
// @BasePath /api
package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/santiago-buildit/code-challenge/backend/internal/routes"
	"log"
	"strings"
)

var ginLambda *ginadapter.GinLambda

// Cold start code here
func init() {
	log.Println("Initializing application...")
	ginLambda = ginadapter.New(routes.SetupRouter())
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Strip /api prefix if present
	if strings.HasPrefix(req.Path, "/api/") {
		req.Path = strings.TrimPrefix(req.Path, "/api")
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {

	// Start handler
	lambda.Start(handler)
}
