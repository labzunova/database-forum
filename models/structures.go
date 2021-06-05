package models

import "time"

type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

// Status required:
//      - user
//      - forum
//      - thread
//      - post
type Status struct {
	User   int32 `json:"user"`   // Кол-во пользователей в базе данных
	Forum  int32 `json:"forum"`  // Кол-во разделов в базе данных
	Thread int32 `json:"thread"` // Кол-во веток обсуждения в базе данных
	Post   int64 `json:"post"`   // Кол-во сообщений в базе данных
}

// User Информация о пользователе.
//     required:
//      - fullname
//      - email
type User struct {
	About    string `json:"about"`    // Описание пользователя
	Email    string `json:"email"`    // Почтовый адрес пользователя (уникальное поле)
	FullName string `json:"fullname"` // Полное имя пользователя
	Nickname string `json:"nickname"` //Имя пользователя (уникальное поле)
	//Данное поле допускает только латиницу, цифры и знак подчеркивания.
	//Сравнение имени регистронезависимо
}

// Forum Информация о форуме.
//    required:
//      - title
//      - user
//      - slug
type Forum struct {
	Slug    string `json:"slug"`    // Человекопонятный URL
	Title   string `json:"title"`   // Название форума
	User    string `json:"user"`    // Nickname пользователя, который отвечает за форум
	Posts   int    `json:"posts"`   // Общее кол-во сообщений в данном форуме
	Threads int    `json:"threads"` // Общее кол-во ветвей обсуждения в данном форуме
}

// Thread Ветка обсуждения на форуме.
//     required:
//      - title
//      - author
//      - message
type Thread struct {
	Author  string    `json:"author"`  // Пользователь, создавший данную тему.
	Created time.Time `json:"created"` // Дата создания ветки на форуме.
	Forum   string    `json:"forum"`   // Форум, в котором расположена данная ветка обсуждения.
	ID      int       `json:"id"`      // Идентификатор ветки обсуждения.
	Message string    `json:"message"` // Описание ветки обсуждения.
	Title   string    `json:"title"`   // Заголовок ветки обсуждения.
	Slug    string    `json:"slug"`    // Человекопонятный URL
	Votes   int       `json:"votes"`   // Кол-во голосов непосредственно за данное сообщение форума.
}

// Post Сообщение внутри ветки обсуждения на форуме.
//     required:
//      - author
//      - message
type Post struct { //  Сообщение внутри ветки обсуждения на форуме.
	Author   string    `json:"author"`   // Автор, написавший данное сообщение.
	Created  time.Time `json:"created"`  // Дата создания сообщения на форуме.
	Forum    string    `json:"forum"`    // Идентификатор форума (slug) данного сообещния.
	ID       int       `json:"id"`       // Идентификатор данного сообщения.
	Message  string    `json:"message"`  // Собственно сообщение форума.
	Parent   int       `json:"parent"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	IsEdited bool      `json:"isEdited"` // Истина, если данное сообщение было изменено.
	Thread   int    `json:"thread"`   // Идентификатор форума (slug) данного сообещния.
}

// PostUpdate Сообщение для обновления сообщения внутри ветки на форуме. Пустые параметры остаются без изменений.
type PostUpdate struct {
	Message string `json:"message"` // Собственно сообщение форума.
}

// PostFull Полная информация о сообщении, включая связанные объекты
type PostFull struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}

// Vote Информация о голосовании пользователя.
//     required:
//      - nickname
//      - voice
type Vote struct {
	Nickname string `json:"nickname"` // Идентификатор пользователя.
	Voice    int32  `json:"voice"`    // Отданный голос.
}
