package http

import (
	"GolangBackendDiploma26/internal/service"
	"net/http"
	"strconv"
)

type BatteryHandler struct {
	batteryService *service.BatteryService
}

func NewBatteryHandler(batteryService *service.BatteryService) *BatteryHandler {
	return &BatteryHandler{batteryService: batteryService}
}

// List возвращает список аккумуляторов
// @Summary      Список аккумуляторов
// @Description  Возвращает список всех доступных аккумуляторов с пагинацией.
// @Description  По умолчанию возвращает 20 записей на странице, максимум 100.
// @Tags         batteries
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Номер страницы (начиная с 1)"  default(1)  minimum(1)
// @Param        limit  query     int  false  "Количество записей на странице"  default(20)  minimum(1)  maximum(100)
// @Success      200    {object}  service.BatteryListResponse  "Список аккумуляторов с пагинацией"
// @Failure      500    {object}  map[string]string            "Внутренняя ошибка сервера"
// @Router       /batteries [get]
func (h *BatteryHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}

	resp, err := h.batteryService.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
