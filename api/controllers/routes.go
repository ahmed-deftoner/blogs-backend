package controllers

import "github.com/ahmed-deftoner/blogs-backend/api/middleware"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middleware.MiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middleware.MiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/confirm", middleware.MiddlewareJSON(s.ConfirmEmail)).Methods("GET")

	//Users routes
	s.Router.HandleFunc("/users", middleware.MiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middleware.MiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middleware.MiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middleware.MiddlewareJSON(middleware.MiddlewareAuth(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middleware.MiddlewareAuth(s.DeleteUser)).Methods("DELETE")

	//Posts routes
	s.Router.HandleFunc("/posts", middleware.MiddlewareJSON(s.CreatePost)).Methods("POST")
	s.Router.HandleFunc("/posts", middleware.MiddlewareJSON(s.GetPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middleware.MiddlewareJSON(s.GetPost)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middleware.MiddlewareJSON(middleware.MiddlewareAuth(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/posts/{id}", middleware.MiddlewareAuth(s.DeletePost)).Methods("DELETE")
}
