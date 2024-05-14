package handler

import (
	"bash_server_project/pkg/repository"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
	"strconv"
)

func createCommand(c *gin.Context) {
	var commands *repository.Command

	if err := c.ShouldBindJSON(&commands); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := executeBashScript(commands.Script)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO commands (script, result) VALUES ($1, $2) RETURNING id"
	err = repository.Db.QueryRow(query, commands.Script, res).Scan(&commands.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, commands)

}

func executeBashScript(script string) (string, error) {
	cmd := exec.Command("bash", "-c", script)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}

func getCommands(c *gin.Context) {
	var commands []repository.Command
	err := repository.Db.Select(&commands, "SELECT * FROM commands")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commands)
}

func getCommand(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var commands repository.Command
	err := repository.Db.Get(&commands, "SELECT * FROM commands WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commands)
}

func InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/commands", createCommand)
	router.GET("/commands", getCommands)
	router.GET("/commands/:id", getCommand)

	return router
}
