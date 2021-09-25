package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/devstack/example/apps/webapp/database"
	"github.com/razorpay/devstack/example/apps/webapp/models"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/helpers"
)

func Status(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func DeletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person models.Person
	span, ctx := helpers.GetChildSpanFromContext(c, "deletePerson", nil)
	defer span.Finish()
	db := database.GetDB(ctx)
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdatePerson(c *gin.Context) {

	var person models.Person
	id := c.Params.ByName("id")
	span, ctx := helpers.GetChildSpanFromGinContext(c, "updatePerson", nil)
	defer span.Finish()
	db := database.GetDB(ctx)
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)

	db.Save(&person)
	c.JSON(200, person)

}

func CreatePerson(c *gin.Context) {

	var person models.Person
	c.BindJSON(&person)
	span, ctx := helpers.GetChildSpanFromGinContext(c, "createPerson", nil)
	defer span.Finish()
	db := database.GetDB(ctx)
	db.Create(&person)
	c.JSON(200, person)
}

func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person models.Person
	span, ctx := helpers.GetChildSpanFromGinContext(c, "getPerson", nil)
	defer span.Finish()
	db := database.GetDB(ctx)
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}
func GetPeople(c *gin.Context) {
	var people []models.Person
	span, ctx := helpers.GetChildSpanFromGinContext(c, "getPeople", nil)
	defer span.Finish()
	db := database.GetDB(ctx)
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}

}
