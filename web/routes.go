package web

// ConfigureRoute list of routes
func (s *Server) ConfigureRoute() {
	s.Get("/", homeView)
	s.Get("/ns/{articleID}/{userID}", articleRedirectView)
	s.Get("/article/{id}", articleView)
	s.Get("/search", searchAPI)
	// add api urls
	s.Get("/a/v1/facebook_users", facebookUsersView)
}
