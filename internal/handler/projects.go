package handler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gorest-api/internal/dto"
	"gorest-api/internal/logs"
	"gorest-api/internal/utils"
)

// @Summary Create a project
// @Security ApiKeyAuth
// @Tags projects
// @Description Creating a project
// @ID project-id
// @Accept  json
// @Produce  json
// @Param input body model.Project true "project information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects [post]
func (h *Handler) createProject(c *fiber.Ctx) error {
	log.Println("Creating a project... ")

	err := h.userIdentity(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	input := dto.CreateProjectDto{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors": true,
			"msg":    "invalid input body",
		})
	}

	id, err := h.services.Projects.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors": true,
			"msg":    utils.ValidatorErrors(err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":    false,
		"msg":       nil,
		"projectId": id,
	})
}

// @Summary Get all projects
// @Security ApiKeyAuth
// @Tags projects
// @Description Get all projects from database
// @ID get-all-projects
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects [get]
func (h *Handler) getAllProjects(c *fiber.Ctx) error {
	logs.Log().Info("Getting all projects... ")

	err := h.userIdentity(c)
	if err != nil {
		return err
	}

	// userId, err := getUserId(c)

	// if err != nil {
	//	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	//		"errors":   true,
	//		"message": utils.ValidatorErrors(err),
	//	})
	// }

	projects, err := h.services.Projects.GetAll(c.UserContext())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(projects),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"projects": projects,
	})
}

// @Summary Get project by title
// @Security ApiKeyAuth
// @Tags projects
// @Description Get project by title from database
// @ID get-project-by-title
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Project
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:title [get]
func (h *Handler) getProjectByTitle(c *fiber.Ctx) error {
	logs.Log().Info("Getting a project by title... ")

	err := h.userIdentity(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	title := c.Params("title")

	project, err := h.services.Projects.GetByTitle(c.UserContext(), userId, title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"project": project,
	})
}

// @Summary Update a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Update a title of project
// @ID update-project-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:id [put]
func (h *Handler) updateProject(c *fiber.Ctx) error {
	log.Println("Updating a project... ")

	err := h.userIdentity(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateProjectDto
	projectId := c.Params("id", "")

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	if err := h.services.Projects.Update(c.UserContext(), userId, projectId, *input.NewTitle); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

// @Summary Delete a project by id
// @Security ApiKeyAuth
// @Tags projects
// @Description Delete a project from database
// @ID delete-project-by-title
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/projects/:id [delete]
func (h *Handler) deleteProject(c *fiber.Ctx) error {
	err := h.userIdentity(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	id := c.Params("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	err = h.services.Projects.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{})
}
