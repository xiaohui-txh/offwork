package handler

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

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

func NearbyOffwork2(c *gin.Context) {
	lngStr := c.Query("lng")
	latStr := c.Query("lat")

	lng, _ := strconv.ParseFloat(lngStr, 64)
	lat, _ := strconv.ParseFloat(latStr, 64)

	// 1. 最近一小时全量下班人数
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

	// 2. 查询 nearby 数据
	queryNearby := func() ([][2]float64, error) {
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
			return nil, err
		}
		defer rows.Close()

		var points [][2]float64
		for rows.Next() {
			var lo, la float64
			rows.Scan(&lo, &la)
			points = append(points, [2]float64{lo, la})
		}
		return points, nil
	}

	points, err := queryNearby()
	if err != nil {
		c.JSON(500, gin.H{"code": 2, "msg": err.Error()})
		return
	}

	// 3. 如果不足 5 条，自动造 5 条数据
	if len(points) < 5 {
		log.Printf("addr[%s:%s] points:%d < 10", lngStr, latStr, len(points))
		tx, _ := db.DB.Begin()

		for i := 0; i < 5; i++ {
			rlng, rlat := randomPointIn3Km(lng, lat)
			_, _ = tx.Exec(
				`INSERT INTO offwork_checkin (lng, lat, created_at)
				 VALUES (?, ?, NOW())`,
				rlng, rlat,
			)
		}
		tx.Commit()

		// 重新查询
		points, _ = queryNearby()
	}

	// 4. 聚合（100m）
	grid := make(map[string]*model.NearbyGridItem)

	for _, p := range points {
		glng := grid100m(p[0])
		glat := grid100m(p[1])
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

// 在 center 点 3km 内生成随机点
func randomPointIn3Km(lng, lat float64) (float64, float64) {
	const earthRadius = 6371.0 // km

	// 随机半径（0~3km）
	r := 3 * math.Sqrt(rand.Float64())
	// 随机角度
	theta := rand.Float64() * 2 * math.Pi

	dLat := r * math.Cos(theta) / earthRadius
	dLng := r * math.Sin(theta) / (earthRadius * math.Cos(lat*math.Pi/180))

	newLat := lat + dLat*180/math.Pi
	newLng := lng + dLng*180/math.Pi

	return newLng, newLat
}

// 100米网格
func grid100m(v float64) float64 {
	return math.Floor(v*1000) / 1000
}
