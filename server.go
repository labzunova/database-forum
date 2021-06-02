package main

import (
	"DBproject/internal/forum"
	http2 "DBproject/internal/forum/delivery/http"
	repository2 "DBproject/internal/forum/repository"
	usecase5 "DBproject/internal/forum/usecase"
	"DBproject/internal/posts"
	http3 "DBproject/internal/posts/delivery/http"
	repository3 "DBproject/internal/posts/repository"
	usecase2 "DBproject/internal/posts/usecase"
	"DBproject/internal/service"
	http5 "DBproject/internal/service/delivery/http"
	repository5 "DBproject/internal/service/repository"
	usecase4 "DBproject/internal/service/usecase"
	"DBproject/internal/threads"
	http4 "DBproject/internal/threads/delivery/http"
	repository4 "DBproject/internal/threads/repository"
	usecase3 "DBproject/internal/threads/usecase"
	"DBproject/internal/user"
	"DBproject/internal/user/delivery/http"
	"DBproject/internal/user/repository"
	"DBproject/internal/user/usecase"
	"database/sql"
	"github.com/labstack/echo/v4"
	"log"
)

func router(e *echo.Echo, user user.UserHandler, forum forum.ForumHandler, posts posts.PostsHandler,
	threads threads.ThreadsHandler, service service.ServiceHandler) {
	e.POST("/forum/create", forum.ForumCreate)
	e.GET("/forum/:slug/details", forum.ForumGetOne)
	e.POST("/forum/:slug/create", forum.ThreadCreate)
	e.GET("/forum/:slug/users", forum.ForumGetUsers)
	e.GET(" /forum/:slug/threads", forum.ForumGetThreads)
	e.GET("/post/:id/details", posts.PostGetOne) // todo
	e.POST("/post/:id/details", posts.PostUpdate)
	e.POST("/service/clear", service.Clear)
	e.GET("/service/status", service.Status)
	e.POST("/thread/:slug_or_id/create", posts.PostsCreate)
	e.GET("/thread/:slug_or_id/details", threads.ThreadGetOne)
	e.POST("/thread/:slug_or_id/details", threads.ThreadUpdate)
	e.GET("/thread/:slug_or_id/posts", threads.ThreadGetPosts) // todo
	e.POST("/thread/:slug_or_id/vote", threads.ThreadVote)
	e.POST("/user/:nickname/create", user.UserCreate)
	e.GET("/user/:nickname/profile", user.UserGetOne)
	e.POST("/user/:nickname/profile", user.UserUpdate)
}

func main() {
	e := echo.New()

	dsn := "jdbc:postgresql://localhost:5432/postgres?user=labzunova&password=1111" // TODO
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUsersRepo(db)
	userUcase := usecase.NewUserUsecase(userRepo)
	userHandler := http.NewUserHandler(userUcase)

	forumRepo := repository2.NewForumRepo(db)
	forumUcase := usecase5.NewForumUsecase(forumRepo)
	forumHandler := http2.NewForumHandler(forumUcase)

	postsRepo := repository3.NewPostsRepo(db)
	postsUcase := usecase2.NewPostsUsecase(postsRepo)
	postsHandler := http3.NewPostsHandler(postsUcase)

	threadsRepo := repository4.NewThreadsRepo(db)
	threadsUcase := usecase3.NewThreadsUsecase(threadsRepo)
	threadsHandler := http4.NewThreadsHandler(threadsUcase)

	serviceRepo := repository5.NewServiceRepo(db)
	serviceUcase := usecase4.NewServiceUsecase(serviceRepo)
	serviceHandler := http5.NewServiceHandler(serviceUcase)

	router(e, userHandler, forumHandler, postsHandler, threadsHandler, serviceHandler)
	e.Logger.Fatal(e.Start(":5000"))
}
