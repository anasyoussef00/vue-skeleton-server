package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/youssef-182/vue-skeleton-server/pkg/controllers"
)

func SetupAuthenticationRoutes(r chi.Router) chi.Router {
	authRoutes := r.Route("/auth", func(r chi.Router) {
		assemblePost(r)
	})
	return authRoutes
}

func assemblePost(r chi.Router) {
	//r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
	//	key := []byte("AAAAAAAAA")
	//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//		"username": "yofs",
	//	})
	//	s, err := token.SignedString(key)
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		w.Write([]byte(fmt.Sprintf("AN ERROR HAS OCCURED: %v", err)))
	//		return
	//	}
	//	w.Write([]byte(s))
	//})
	r.Post("/register", controllers.Register)
	r.Post("/login", controllers.Login)
}
