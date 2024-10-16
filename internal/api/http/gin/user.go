package gin

import (
	"net/http"
	"todo-app/domain"
	"todo-app/pkg/clients"
	"todo-app/pkg/tokenprovider"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	Register(data *domain.UserCreate) error
	Login(data *domain.UserLogin) (tokenprovider.Token, error)
	GetAllUser() ([]domain.User, error)
	GetUserByID(id uuid.UUID) (domain.User, error)
	UpdateUser(id uuid.UUID, user *domain.UserUpdate) error
	DeleteUser(id uuid.UUID) error
}

type userHandler struct {
	userService UserService
}

func NewUserHandler(apiVersion *gin.RouterGroup, svc UserService , middlewareAuth func(c *gin.Context), middlewareRateLimit func(c *gin.Context)) {
	userHandler := &userHandler{
		userService: svc,
	}

	users := apiVersion.Group("/users")
	users.POST("/register", userHandler.RegisterUserHandler)
	users.POST("/login", userHandler.LoginHandler)
	users.GET("/", userHandler.GetAllUserHandler)
	users.GET("/:id",middlewareAuth, userHandler.GetUserHandler)
	users.PATCH("/:id",middlewareAuth, userHandler.UpdateUserHandler)
	users.DELETE("/:id", userHandler.DeleteUserHandler)
}

// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param data body domain.UserCreate true "User creation data"
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users/register [post]
func (h *userHandler) RegisterUserHandler(c *gin.Context) {
	var data domain.UserCreate

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.userService.Register(&data); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(data.ID))
}

// LoginHandler godoc
// @Summary Log in a user
// @Description Authenticate a user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param data body domain.UserLogin true "User login data"
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users/login [post]
func (h *userHandler) LoginHandler(c *gin.Context) {
	var data domain.UserLogin

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	token, err := h.userService.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(token))
}

// GetAllUserHandler godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users [get]
func (h *userHandler) GetAllUserHandler(c *gin.Context) {
	users, err := h.userService.GetAllUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(users))
}

// GetUserHandler godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users/{id} [get]
func (h *userHandler) GetUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(user))
}

// UpdateUserHandler godoc
// @Summary Update a user
// @Description Update a user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body domain.UserUpdate true "User update data"
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users/{id} [patch]
func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))
		return
	}

	user := domain.UserUpdate{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.userService.UpdateUser(id, &user); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}

// DeleteUserHandler godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} clients.SuccessRes
// @Failure      400  {object}  clients.AppError
// @Failure      404  {object}  clients.AppError
// @Failure      500  {object}  clients.AppError
// @Router /users/{id} [delete]
func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, clients.ErrInvalidRequest(err))

		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, clients.SimpleSuccessResponse(true))
}
