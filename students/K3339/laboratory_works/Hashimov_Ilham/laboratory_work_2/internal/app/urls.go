package app

import "net/http"

func NewRouter(homeworkView *HomeworkView, adminView *AdminView, userView *UserView, auth *Auth) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", auth.Login)
	mux.HandleFunc("POST /login", auth.Login)

	mux.HandleFunc("GET /logout", auth.Logout)

	mux.HandleFunc("GET /users", userView.List)

	mux.HandleFunc("GET /users/create", userView.Create)
	mux.HandleFunc("POST /users/create", userView.Create)

	mux.HandleFunc("GET /homework/{id}", homeworkView.Detail)
	mux.HandleFunc("GET /homeworks", homeworkView.List)

	mux.HandleFunc("GET /homework/{id}/update", auth.RequireAdmin(homeworkView.Update))
	mux.HandleFunc("POST /homework/{id}/update", auth.RequireAdmin(homeworkView.Update))

	mux.HandleFunc("GET /homework/{id}/delete", auth.RequireAdmin(homeworkView.Delete))
	mux.HandleFunc("POST /homework/{id}/delete", auth.RequireAdmin(homeworkView.Delete))

	mux.HandleFunc("GET /homeworks/create", auth.RequireAdmin(homeworkView.Create))
	mux.HandleFunc("POST /homeworks/create", auth.RequireAdmin(homeworkView.Create))

	mux.HandleFunc("GET /admin", auth.RequireAdmin(adminView.Index))

	mux.HandleFunc("GET /admin/users", auth.RequireAdmin(adminView.Users))
	mux.HandleFunc("GET /admin/homeworks", auth.RequireAdmin(adminView.Homeworks))
	mux.HandleFunc("GET /admin/submissions", auth.RequireAdmin(adminView.Submissions))

	mux.HandleFunc("GET /admin/users/create", auth.RequireAdmin(adminView.CreateUser))
	mux.HandleFunc("POST /admin/users/create", auth.RequireAdmin(adminView.CreateUser))

	mux.HandleFunc("GET /admin/homeworks/create", auth.RequireAdmin(adminView.CreateHomework))
	mux.HandleFunc("POST /admin/homeworks/create", auth.RequireAdmin(adminView.CreateHomework))

	mux.HandleFunc("GET /admin/submissions/create", auth.RequireAdmin(adminView.CreateSubmission))
	mux.HandleFunc("POST /admin/submissions/create", auth.RequireAdmin(adminView.CreateSubmission))

	return mux
}
