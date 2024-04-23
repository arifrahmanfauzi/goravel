package middleware

import (
	"fmt"
	"io"
	_http "net/http"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Verify() http.Middleware {
	return func(ctx http.Context) {
		token := ctx.Request().Header("Authorization", "")
		if token == "" {
			ctx.Request().AbortWithStatusJson(401, http.Json{
				"status":  "Error",
				"message": "Empty Authorization Header!",
			})
			return
		}
		url := facades.Config().GetString("VERIFY_URL", "goravel")
		client := &_http.Client{}

		// Create a new HTTP client
		req, err := _http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
				"status":  "Error",
				"message": "Cannot send authorize server!",
			})
			return
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{
				"status":  "Error",
				"message": "error sending to server!",
			})
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// Check the response status code
		if resp.StatusCode != http.StatusOK {

			fmt.Println("Error verify: ", err)
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"status":  "Error",
				"message": "Unauthorized!",
			})
			return
		}

		ctx.Request().Next()
	}
}
