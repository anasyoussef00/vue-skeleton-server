package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/youssef-182/vue-skeleton-server/pkg/controllers"
	"github.com/youssef-182/vue-skeleton-server/pkg/router/middlewares"
)

func SetupBookRoutes(r chi.Router) chi.Router {
	bookRoutes := r.Route("/book", func(r chi.Router) {
		assembleBookGet(r)
		assembleBookPost(r)
		assembleBookPut(r)
		assembleBookDelete(r)
	})
	return bookRoutes
}

func assembleBookGet(r chi.Router) {
	r.Get("/", controllers.IndexBook)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.BookCtx)
		r.Get("/{bookID}", controllers.ShowBook)
	})
}

func assembleBookPost(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWT)
		r.Post("/store", controllers.StoreBook)
	})
}

func assembleBookPut(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWT)
		r.Use(middlewares.BookCtx)
		r.Put("/{bookID}", controllers.UpdateBook)
	})
}

func assembleBookDelete(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWT)
		r.Use(middlewares.BookCtx)
		r.Delete("/{bookID}", controllers.DeleteBook)
	})
}
