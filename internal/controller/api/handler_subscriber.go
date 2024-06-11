package api

import (
	"fmt"
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

type subscriberHandler struct {
	t usecase.Subscriber
	l logger.Interface
}

func newSubscriberHandler(t usecase.Subscriber, l logger.Interface) *subscriberHandler {
	return &subscriberHandler{t: t, l: l}
}

func (h *subscriberHandler) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var data entity.Subscriber
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		h.l.Debug("Failed to decode request body to entity.Subscriber", h.l.Err(err))

		render.JSON(w, r, Error("failed to decode request body to entity.Subscriber"))

		return
	}

	h.l.Info("request body decoded to entity.Subscriber successfully", slog.Any("request", data))

	err = h.t.Add(ctx, data)

	if err != nil {
		h.l.Debug("Failed to save data in DB", h.l.Err(err))

		switch err := err.(type) {
		case *pq.Error:
			render.JSON(w, r, Error("provided data is invalid or employee_id to subscribed_to_id assignment already exists"))
			return
		default:
			render.JSON(w, r, err.Error())
			return
		}
	}

	render.JSON(w, r,
		SubscriberResponse{
			Status:     StatusOk,
			Subscriber: &data,
		})
}

func (h *subscriberHandler) remove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employee_id, err := strconv.Atoi(chi.URLParam(r, "employee_id"))

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	subscriber_id, err := strconv.Atoi(chi.URLParam(r, "subscriber_id"))

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	h.l.Info("request body decoded to int successfully", slog.Any("request", employee_id))

	if employee_id == 0 {

		h.l.Info("No data available for employee_id = 0")

		render.JSON(w, r,
			EmployeeResponse{
				Status: "Wrong employee_id. Employee_id should be > 0",
			})

		return
	}

	res, err := h.t.Remove(ctx, employee_id, subscriber_id)

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

func (h *subscriberHandler) listAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			SubscribersResponse{
				Status: "Wrong id. Id should be > 0",
			})

		return
	}

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	res, err := h.t.ListAll(ctx, id)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error(fmt.Sprintf("Database has NO employee with id = %d", id)))

		return
	}

	render.JSON(w, r,
		SubscribersResponse{
			Status:    StatusOk,
			Employees: res,
		})
}

func (h *subscriberHandler) listHappyBirthdays(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if id == 0 {

		h.l.Info("No data available for id = 0")

		render.JSON(w, r,
			SubscribersResponse{
				Status: "Wrong id. Id should be > 0",
			})

		return
	}

	if err != nil {

		h.l.Debug("id parameter in URL is not integer or empty", h.l.Err(err))

		render.JSON(w, r, Error("Unable to retrieve id from URL"))

		return
	}

	res, err := h.t.ListHappyBirthdays(ctx, id)

	if err != nil {
		h.l.Debug("Failed to get data from DB", h.l.Err(err))

		render.JSON(w, r, Error(fmt.Sprintf("Database has NO actor with id = %d", id)))

		return
	}

	if len(res) == 0 {
		render.JSON(w, r,
			SubscribersResponse{
				Status: StatusEmpty,
			})

		return
	}

	render.JSON(w, r,
		SubscribersResponse{
			Status:    StatusOk,
			Employees: res,
		})
}
