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
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func router(e *echo.Echo, user user.UserHandler, forum forum.ForumHandler, posts posts.PostsHandler,
	threads threads.ThreadsHandler, service service.ServiceHandler) {
	e.POST("/api/forum/create", forum.ForumCreate)
	e.GET("/api/forum/:slug/details", forum.ForumGetOne)
	e.POST("/api/forum/:slug/create", forum.ThreadCreate)
	e.GET("/api/forum/:slug/users", forum.ForumGetUsers)
	e.GET("/api/forum/:slug/threads", forum.ForumGetThreads)
	e.GET("/api/post/:id/details", posts.PostGetOne)
	e.POST("/api/post/:id/details", posts.PostUpdate)
	e.POST("/api/service/clear", service.Clear)
	e.GET("/api/service/status", service.Status)
	e.POST("/api/thread/:slug_or_id/create", posts.PostsCreate)
	e.GET("/api/thread/:slug_or_id/details", threads.ThreadGetOne)
	e.POST("/api/thread/:slug_or_id/details", threads.ThreadUpdate)
	e.GET("/api/thread/:slug_or_id/posts", threads.ThreadGetPosts) // todo
	e.POST("/api/thread/:slug_or_id/vote", threads.ThreadVote)
	e.POST("/api/user/:nickname/create", user.UserCreate)
	e.GET("/api/user/:nickname/profile", user.UserGetOne)
	e.POST("/api/user/:nickname/profile", user.UserUpdate)
}

func main() {
	e := echo.New()

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", "lbznv", "1111", "forums")
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
