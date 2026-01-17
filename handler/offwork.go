package handler

import (
	"fmt"
	"math"
	"net/http"

	"offwork-backend/db"
	"offwork-backend/model"

	"github.com/gin-gonic/gin"
)

// POST /api/v1/offwork/checkin
func OffworkCheckin(c *gin.Context) {
	var req model.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	_, err := db.DB.Exec(
		`INSERT INTO offwork_checkin (lng, lat) VALUES (?, ?)`,
		req.Lng, req.Lat,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2,
			"msg":  err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /api/v1/offwork/nearby
func NearbyOffwork(c *gin.Context) {
	lng := c.Query("lng")
	lat := c.Query("lat")

	// 1. 查询最近一小时的全量下班人数
	var allCount int
	err := db.DB.QueryRow(`
SELECT COUNT(*)
FROM offwork_checkin
WHERE created_at >= NOW() - INTERVAL 1 HOUR
`).Scan(&allCount)
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	rows, err := db.DB.Query(`
SELECT lng, lat
FROM offwork_checkin
WHERE
  created_at >= NOW() - INTERVAL 1 HOUR
  AND
  6371 * acos(
    cos(radians(?)) * cos(radians(lat)) *
    cos(radians(lng) - radians(?)) +
    sin(radians(?)) * sin(radians(lat))
  ) <= 3
`, lat, lng, lat)
	if err != nil {
		c.JSON(500, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	defer rows.Close()

	grid := make(map[string]*model.NearbyGridItem)
	for rows.Next() {
		var lngF, latF float64
		_ = rows.Scan(&lngF, &latF)

		glng := grid100m(lngF)
		glat := grid100m(latF)
		key := fmt.Sprintf("%.3f_%.3f", glng, glat)

		if _, ok := grid[key]; !ok {
			grid[key] = &model.NearbyGridItem{
				Lng:   glng,
				Lat:   glat,
				Count: 0,
			}
		}
		grid[key].Count++
	}

	resp := model.NearbyResponse{
		AllCount: allCount,
	}

	for _, v := range grid {
		resp.OffworkData = append(resp.OffworkData, *v)
	}

	c.JSON(http.StatusOK, resp)
}

// 100米网格
func grid100m(v float64) float64 {
	return math.Floor(v*1000) / 1000
}
