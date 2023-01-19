package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

func GetScore(c *fiber.Ctx) error {
	users := []models.User{}
	logs := []models.Log{}
	questions := []models.Question{}

	type ScoreUser struct {
		Score    uint
		Username string
	}

	responseScore := []ScoreUser{}

	database.Database.Db.Find(&users)

	for _, user := range users {
		var sum uint = 0
		database.Database.Db.Find(&logs, "status=? AND user_refer=?", true, user.ID)
		for _, log := range logs {
			database.Database.Db.Find(&questions, "id=?", log.QuestionRefer)
			for _, question := range questions {
				sum = question.Score + sum
				fmt.Println(log.Status)
				fmt.Println(sum)
			}
		}
		database.Database.Db.Find(&user, "id=?", user.ID)
		responseScore = append(responseScore, ScoreUser{Score: sum, Username: user.Username})
	}

	for i := 0; i < len(responseScore); i++ {
		var temp ScoreUser
		for j := 0; j < len(responseScore); j++ {
			if responseScore[j].Score <= responseScore[i].Score {
				temp = responseScore[i]
				responseScore[i] = responseScore[j]
				responseScore[j] = temp
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Scores listed succesfully",
		"data":    responseScore,
	})

}
