package main

import (
	"fmt"
	"net/http"

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

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests",
	Help: "Requests for RPS metric",
})

func counterMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			requests.Inc()
			handler.ServeHTTP(response, request)
		})
	}
}

func main() {
	prometheus.MustRegister(requests)
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

	//api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api := mux.NewRouter()
	api.Use(counterMiddleware(api))
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	api.HandleFunc("/api/forum/create", forumHandler.ForumCreate).Methods(http.MethodPost)
	api.HandleFunc("/api/forum/{slug}/details", forumHandler.ForumGetOne).Methods(http.MethodGet)
	api.HandleFunc("/api/forum/{slug}/create", forumHandler.ThreadCreate).Methods(http.MethodPost)
	api.HandleFunc("/api/forum/{slug}/users", forumHandler.ForumGetUsers).Methods(http.MethodGet)
	api.HandleFunc("/api/forum/{slug}/threads", forumHandler.ForumGetThreads).Methods(http.MethodGet)
	api.HandleFunc("/api/post/{id}/details", postsHandler.PostGetOne).Methods(http.MethodGet)
	api.HandleFunc("/api/post/{id}/details", postsHandler.PostUpdate).Methods(http.MethodPost)
	api.HandleFunc("/api/service/clear", serviceHandler.Clear).Methods(http.MethodPost)
	api.HandleFunc("/api/service/status", serviceHandler.Status).Methods(http.MethodGet)
	api.HandleFunc("/api/thread/{slug_or_id}/create", postsHandler.PostsCreate).Methods(http.MethodPost)
	api.HandleFunc("/api/thread/{slug_or_id}/details", threadsHandler.ThreadGetOne).Methods(http.MethodGet)
	api.HandleFunc("/api/thread/{slug_or_id}/details", threadsHandler.ThreadUpdate).Methods(http.MethodPost)
	api.HandleFunc("/api/thread/{slug_or_id}/posts", threadsHandler.ThreadGetPosts).Methods(http.MethodGet)
	api.HandleFunc("/api/thread/{slug_or_id}/vote", threadsHandler.ThreadVote).Methods(http.MethodPost)
	api.HandleFunc("/api/user/{nickname}/create", userHandler.UserCreate).Methods(http.MethodPost)
	api.HandleFunc("/api/user/{nickname}/profile", userHandler.UserGetOne).Methods(http.MethodGet)
	api.HandleFunc("/api/user/{nickname}/profile", userHandler.UserUpdate).Methods(http.MethodPost)

	http.ListenAndServe(":5000", api)
}
