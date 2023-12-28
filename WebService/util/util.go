package util

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func GetJSONBody(c echo.Context) (map[string]interface{}, error) {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("GetJSONBody failed! error: %w", err)
	}

	return jsonBody, nil
}
