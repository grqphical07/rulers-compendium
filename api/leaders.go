package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/grqphical07/rulers-compendium/database"

	"github.com/labstack/echo/v4"
)

func getLeadersByCivilization(civilzation string, limit int, db *database.Database) []database.Leader {
	civ_leaders := make([]database.Leader, 0)

	for _, leader := range db.Leaders {
		if limit != 0 && len(civ_leaders) == limit {
			return civ_leaders
		}

		if strings.ToLower(leader.Civilization) == strings.ToLower(civilzation) {
			civ_leaders = append(civ_leaders, leader)
		}
	}

	return civ_leaders
}

// Index godoc
// @Description Get the status of the API Server
// @Tags root
// @Accept */*
// @Produce plain
// @Success 200 {object} string
// @Router /api/v1/ [get]
func (r *Router) Index(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// GetLeaders godoc
// @Description Get's all leaders from civ 6 including their agenda's and abilities
// @Tags leaders
// @Accept */*
// @Produce json
// @Param limit query int false "limits amount of results returned"
// @Param civilization query string false "filters leaders by civilization"
// @Success 200 {object} []database.Leader
// @Failure 400 {object} string "Invalid limit value"
// @Router /api/v1/leaders [get]
func (r *Router) GetLeaders(c echo.Context) error {
	limit := c.QueryParam("limit")
	civilzation := c.QueryParam("civilization")

	if limit == "" {
		if civilzation != "" {
			return c.JSON(http.StatusOK, getLeadersByCivilization(civilzation, 0, r.db))
		}
		leaders := r.db.Leaders
		return c.JSON(http.StatusOK, leaders)
	} else {
		limit, err := strconv.Atoi(limit)
		if err != nil || limit <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid limit value. Must be number above zero")
		}

		if civilzation != "" {
			return c.JSON(http.StatusOK, getLeadersByCivilization(civilzation, limit, r.db))
		}

		leaders := r.db.Leaders[:limit]
		return c.JSON(http.StatusOK, leaders)
	}

}

// GetLeader godoc
// @Description Gets a single leader from civ 6
// @Tags leaders
// @Accept */*
// @Produce json
// @Param name path string true "leader to get"
// @Success 200 {object} database.Leader
// @Router /api/v1/leaders/{name} [get]
func (r *Router) GetLeader(c echo.Context) error {
	name := c.Param("name")
	name = strings.ReplaceAll(name, "%20", " ")

	for _, leader := range r.db.Leaders {
		if leader.Name == name {
			return c.JSON(http.StatusOK, leader)
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "Could not find leader")
}
