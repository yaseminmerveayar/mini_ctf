package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

type Log struct {
	UserId     uuid.UUID
	QuestionId uuid.UUID
	Status     bool      `json:"status"`
	Date       time.Time `json:"date"`
}

func CreateResponseLog(log models.Log) Log {
	return Log{UserId: log.UserRefer, QuestionId: log.QuestionRefer, Status: log.Status, Date: log.Date}
}

func AnswerQuestion(c *fiber.Ctx) error {
	id := c.Params("id")

	var question models.Question
	var log models.Log

	database.Database.Db.Find(&question, "id=?", id)
	if question.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Question not found",
			"data":    nil,
		})
	}

	if question.IsHidden {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"data":    nil,
		})
	}

	type AnswerQuestion struct {
		Flag string `json:"flag"`
	}

	var data AnswerQuestion

	if err := c.BodyParser(&data); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Review your input",
			"data":    err,
		})
	}

	user_id, err := ConvertSessId(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "uuid cant convert",
		})
	}

	if CheckLog(true, user_id, question.ID, &log) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Already answer question successfully",
		})
	}

	log.UserRefer = user_id
	log.QuestionRefer = question.ID
	log.Date = time.Now()

	if question.Flag == data.Flag {
		log.Status = true
		database.Database.Db.Create(&log)
		responseLog := CreateResponseLog(log)
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Question answered succesfully",
			"data":    responseLog,
		})
	} else {
		log.Status = false
		database.Database.Db.Create(&log)
		responseLog := CreateResponseLog(log)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Answer is wrong",
			"data":    responseLog,
		})
	}

}
func GetLogs(c *fiber.Ctx) error {
	user_id, _ := ConvertSessId(c)
	fmt.Println(user_id)
	if !CheckIfAdmin(user_id) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	logs := []models.Log{}

	database.Database.Db.Find(&logs)
	responseLogs := []Log{}
	for _, log := range logs {
		responseLog := CreateResponseLog(log)
		responseLogs = append(responseLogs, responseLog)
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Logs listed succesfully",
		"data":    responseLogs,
	})
}
