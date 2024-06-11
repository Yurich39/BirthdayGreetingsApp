package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lib/pq"
	"golang.org/x/exp/slog"

	"BirthdayGreetingsService/internal/entity"
	"BirthdayGreetingsService/internal/usecase"
	"BirthdayGreetingsService/pkg/logger"
)

type employeeHandler struct {
	t usecase.Employee
	l logger.Interface
}

func newEmployeeHandler(t usecase.Employee, l logger.Interface) *employeeHandler {
	return &employeeHandler{t: t, l: l}
}

func (h *employeeHandler) save(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.EmployeeData
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.EmployeeData", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.EmployeeData"))

		return
	}

	h.l.Info("request body decoded to entity.EmployeeData successfully", slog.Any("request", data))

	res, err := h.t.Save(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		if err == sql.ErrNoRows {
			render.JSON(w, r, Error("employee already exists"))
			return
		} else {
			switch err := err.(type) {
			case *pq.Error:
				render.JSON(w, r, Error("provided data is invalid"))
				return
			default:
				render.JSON(w, r, err.Error())
				return
			}
		}
	}

	render.JSON(w, r,
		EmployeeResponse{
			Status:   StatusOk,
			Employee: &res,
		})
}

func (h *employeeHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updates entity.Employee
	err := render.DecodeJSON(r.Body, &updates)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Employee", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Employee"))

		return
	}

	h.l.Info("request body decoded to entity.Employee successfully", slog.Any("request", updates))

	res, err := h.t.Update(ctx, updates)

	if err != nil {
		h.l.Debug("Failed to update data in DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to update employee data in DB"))

		return
	}

	render.JSON(w, r,
		EmployeeResponse{
			Status:   StatusOk,
			Employee: &res,
		})
}

func (h *employeeHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	h.l.Info("request body decoded to int successfully", slog.Any("request", id))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			EmployeeResponse{
				Status: "Wrong id. Id should be > 0",
			})

		return
	}

	res, err := h.t.Delete(ctx, id)

	if err != nil {
		h.l.Debug("Failed to delete data from DB", h.l.Err(err))

		render.JSON(w, r, Error("Unable to delete data from DB"))

		return
	}

	render.JSON(w, r,
		EmployeeResponse{
			Status:   StatusOk,
			Employee: &res,
		})
}
