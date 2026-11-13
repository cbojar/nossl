package handlers

import "github.com/gin-gonic/gin"

func RootHandler(context *gin.Context) {
	context.Header("Content-Type", "text/html")
	context.String(200, `<!doctype html>
	<html>
		<head>
			<title>NoSSL Proxy</title>
			<meta charset="utf-8" />
			<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
			<meta name="viewport" content="width=device-width" />
		</head>
		<body>
				<h1>NoSSL Proxy</h1>
				<form action="/proxy" method="GET">
				<label for="url">URL:</label> <input type="text" name="url" id="url"/>
				<input type="submit"/>
				</form>
		</body>
	</html>`)
}
