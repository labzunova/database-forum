package main

import (
	"DBproject/internal/forum/delivery"
	repository2 "DBproject/internal/forum/repository"
	http3 "DBproject/internal/posts/delivery/http"
	repository3 "DBproject/internal/posts/repository"
	http5 "DBproject/internal/service/delivery/http"
	repository5 "DBproject/internal/service/repository"
	http4 "DBproject/internal/threads/delivery/http"
	repository4 "DBproject/internal/threads/repository"
	http0 "DBproject/internal/user/delivery/http"
	"DBproject/internal/user/repository"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	connectionString := "postgres://postgres:1111@3.22.112.0/forums?sslmode=disable"
	//connectionString := "postgres://lbznv:1111@localhost/forums?sslmode=disable"
	config, err := pgx.ParseURI(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: 100,
			AfterConnect:   nil,
			AcquireTimeout: 0,
		})

	userRepo := repository.NewUsersRepo(db)
	//userUcase := usecase.NewUserUsecase(userRepo)
	userHandler := http0.NewUserHandler(userRepo)

	forumRepo := repository2.NewForumRepo(db)
	//forumUcase := usecase5.NewForumUsecase(forumRepo)
	forumHandler := delivery.NewForumHandler(forumRepo)

	postsRepo := repository3.NewPostsRepo(db)
	//postsUcase := usecase2.NewPostsUsecase(postsRepo)
	postsHandler := http3.NewPostsHandler(postsRepo)

	threadsRepo := repository4.NewThreadsRepo(db)
	//threadsUcase := usecase3.NewThreadsUsecase(threadsRepo)
	threadsHandler := http4.NewThreadsHandler(threadsRepo)

	serviceRepo := repository5.NewServiceRepo(db)
	//serviceUcase := usecase4.NewServiceUsecase(serviceRepo)
	serviceHandler := http5.NewServiceHandler(serviceRepo)

	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api.HandleFunc("/forum/create", forumHandler.ForumCreate).Methods(http.MethodPost)
	api.HandleFunc("/forum/{slug}/details", forumHandler.ForumGetOne).Methods(http.MethodGet)
	api.HandleFunc("/forum/{slug}/create", forumHandler.ThreadCreate).Methods(http.MethodPost)
	api.HandleFunc("/forum/{slug}/users", forumHandler.ForumGetUsers).Methods(http.MethodGet)
	api.HandleFunc("/forum/{slug}/threads", forumHandler.ForumGetThreads).Methods(http.MethodGet)
	api.HandleFunc("/post/{id}/details", postsHandler.PostGetOne).Methods(http.MethodGet)
	api.HandleFunc("/post/{id}/details", postsHandler.PostUpdate).Methods(http.MethodPost)
	api.HandleFunc("/service/clear", serviceHandler.Clear).Methods(http.MethodPost)
	api.HandleFunc("/service/status", serviceHandler.Status).Methods(http.MethodGet)
	api.HandleFunc("/thread/{slug_or_id}/create", postsHandler.PostsCreate).Methods(http.MethodPost)
	api.HandleFunc("/thread/{slug_or_id}/details", threadsHandler.ThreadGetOne).Methods(http.MethodGet)
	api.HandleFunc("/thread/{slug_or_id}/details", threadsHandler.ThreadUpdate).Methods(http.MethodPost)
	api.HandleFunc("/thread/{slug_or_id}/posts", threadsHandler.ThreadGetPosts).Methods(http.MethodGet)
	api.HandleFunc("/thread/{slug_or_id}/vote", threadsHandler.ThreadVote).Methods(http.MethodPost)
	api.HandleFunc("/user/{nickname}/create", userHandler.UserCreate).Methods(http.MethodPost)
	api.HandleFunc("/user/{nickname}/profile", userHandler.UserGetOne).Methods(http.MethodGet)
	api.HandleFunc("/user/{nickname}/profile", userHandler.UserUpdate).Methods(http.MethodPost)

	http.ListenAndServe(":5000", api)
}
