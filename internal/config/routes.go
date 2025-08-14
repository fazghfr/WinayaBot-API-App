package config

import (
	"Discord_API_DB_v1/internal/handler"
	"Discord_API_DB_v1/internal/repository"
	"Discord_API_DB_v1/internal/service"
)

// login routes
func RegisterRoutes(r *HttpServer) {

	engine := r.gin_object
	db := InitDB()
	UserRepository := repository.InitUserRepo(db)
	UserService := service.InitUserService(UserRepository)
	Userhandler := handler.InitUserHandler(UserService)

	{
		User := engine.Group("/api/user")
		{
			User.POST("/init", Userhandler.InitRegistration)
		}
	}

	TaskRepository := repository.InitTaskRepository(db)
	TaskService := service.InitTaskService(TaskRepository, UserRepository)
	TaskHandler := handler.InitTaskHandler(TaskService)

	{
		Task := engine.Group("/api/task")
		{
			Task.POST("/create", TaskHandler.CreateNewTask)
			Task.GET("/user", TaskHandler.GetTasksByUser)
			Task.PUT("/edit/:task_id", TaskHandler.EditTaskByTaskID)
			Task.DELETE("/delete/:task_id", TaskHandler.DeleteTaskByID)
		}
	}
}
