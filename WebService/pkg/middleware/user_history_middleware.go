package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

/*
新增使用者紀錄的 Echo Middleware
*/
func UserHistory(databaseRepository repository.DatabaseRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				c.Error(err)
			}

			errorMessage := "Create UserHistory failed! error: %w"

			userId := ""

			if jwtClaims := util.GetJWTClaims(c); jwtClaims != nil {
				userId = jwtClaims.UserId
			}

			path := c.Path()

			// 不用記錄 "瀏覽使用者紀錄" 的操作
			if path != "/api/restricted/user/history" {
				_, err := databaseRepository.CreateUserHistory(
					c.Request().Context(),
					model.UserHistory{
						UserId: userId,
						Method: c.Request().Method,
						Path:   path,
					},
				)
				if err != nil {
					middleErr := fmt.Errorf(errorMessage, err)
					c.Logger().Error(middleErr)
					return middleErr
				}
			}

			return nil
		}
	}
}
