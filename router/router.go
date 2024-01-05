package router

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"sgnacos/conf"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Http(app *fiber.App) {
	var configDir = conf.BaseConf.DataDir

	g := app.Group("/nacos")

	// 登录
	g.Post("/v1/auth/users/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"accessToken": "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuYWNvcyIsImV4cCI6MTYwNTYyOTE2Nn0.2TogGhhr11_vLEjqKko1HJHUJEmsPuCxkur-CfNojDo",
			"tokenTtl":    18000,
			"globalAdmin": true,
		})
	})

	// 获取配置
	g.Get("/v1/cs/configs", func(c *fiber.Ctx) error {
		dataId, tenant, group := c.Query("dataId"), c.Query("tenant"), c.Query("group")
		file := path.Join(configDir, tenant, dataId)
		logrus.Debugf("获取配置: [%s]-[%s]-[%s]: [%s]", tenant, group, dataId, file)
		return c.SendFile(file)
	})

	// 监听配置
	g.Post("/v1/cs/configs/listener", func(c *fiber.Ctx) error {
		logrus.Debugf("configs listener [%s]", c.FormValue("Listening-Configs"))

		response := ""
		for _, config := range strings.Split(c.FormValue("Listening-Configs"), "\x01") {
			strs := strings.Split(config, "\x02")
			if len(strs) != 4 {
				continue
			}
			dataId := strs[0]
			tenant := strs[3]

			content, err := os.ReadFile(path.Join(configDir, tenant, dataId))
			if os.IsNotExist(err) {
				continue
			}
			md5str := fmt.Sprintf("%x", md5.Sum(content))
			if md5str != strs[2] {
				response = response + dataId + "\x02" + strs[1] + "\x02" + strs[3] + "\x01"
			}
		}

		if response == "" {
			// 没变化，延迟返回
			if timeout, err := strconv.Atoi(c.Get("Long-Pulling-Timeout")); err == nil {
				timeout := time.Duration(timeout) * time.Millisecond
				logrus.Debugf("long pulling timeout: %s", timeout.String())
				time.Sleep(timeout)
			}
		}
		return c.SendString(response)
	})

	// 查询实例列表
	g.Get("/v1/ns/instance/list", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{})
	})

	g.Post("/v1/ns/instance", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	g.Get("/v1/ns/service/list", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"count": 0,
			"doms":  []string{},
		})
	})

	g.Put("/v1/ns/instance/beat", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"clientBeatInterval": 3600,
			"lightBeatEnabled":   true,
		})
	})
}
