package login

import "net/http"

const loginPage =
	`<html>
		<head>
			<title>Login</title>
		</head>
		<body></body>
	</html>`

func NewLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(loginPage))
	}
}
