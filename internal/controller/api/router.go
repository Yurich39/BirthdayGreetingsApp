package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"BirthdayGreetingsService/config"
	"BirthdayGreetingsService/internal/controller/middleware/filter"
	"BirthdayGreetingsService/internal/controller/middleware/pagination"

	// "BirthdayGreetingsService/internal/controller/middleware/sort"
	"BirthdayGreetingsService/internal/usecase"
	"BirthdayGreetingsService/pkg/logger"
)

func NewRouter(cfg *config.Config, router *chi.Mux, l logger.Interface, a usecase.Employee, m usecase.Subscriber) {
	// Middleware для общего использования
	commonMiddleware := chi.Chain(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middleware.URLFormat,
	)

	// Middleware для авторизации администратора
	adminAuthMiddleware := middleware.BasicAuth("BirthdayGreetingsService", map[string]string{
		cfg.HTTPServer.User: cfg.HTTPServer.Pass,
	})

	employee := newEmployeeHandler(a, l)
	subscriber := newSubscriberHandler(m, l)

	router.Route("/employee", func(r chi.Router) {
		r.Use(commonMiddleware.Handler)
		r.With(adminAuthMiddleware).Post("/save", employee.save)
		r.With(adminAuthMiddleware).Put("/update", employee.update)
		r.With(adminAuthMiddleware).Delete("/delete/{id}", employee.delete)
	})

	router.Route("/employees", func(r chi.Router) {
		r.Route("/list", func(r chi.Router) {
			r.Use(commonMiddleware.Handler)
			r.With(filter.Middleware).Get("/", employee.list)
			r.With(filter.Middleware, pagination.Middleware).Get("/next", employee.next)
		})
	})

	router.Route("/subscribers", func(r chi.Router) {
		r.Use(commonMiddleware.Handler)
		r.Post("/add", subscriber.add)
		r.Delete("/remove/{employee_id}&{subscriber_id}", subscriber.remove)
		r.Get("/find/{id}", subscriber.listAll)
		r.Get("/find_employees_with_happy_birthdays/{id}", subscriber.listHappyBirthdays)
	})
}
