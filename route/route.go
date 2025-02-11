package route

import (
	"database/sql"
	controller "signUp/controllers"

	"github.com/ReneKroon/ttlcache"
	"github.com/gorilla/mux"
)

type Route struct {
	Db     *sql.DB
	Router *mux.Router
	Cache  *ttlcache.Cache
}

func (r *Route) Routes() {
	s := controller.SignUpController{
		Db:    r.Db,
		Cache: r.Cache,
	}

	//will be writing all route handlers over here
	r.Router.HandleFunc("/signup", s.UserDetails).Methods("POST")
	r.Router.HandleFunc("/verify-email", s.VerifyEmail).Methods("POST")

	g := controller.GoogleSignUp{
		Db: r.Db,
	}
	//will show the google signin page
	r.Router.HandleFunc("/google/login", g.GoogleLogin).Methods("GET")
	r.Router.HandleFunc("/google/callback", g.GoogleCallback).Methods("GET")
}
