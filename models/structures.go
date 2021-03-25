package models

import "time"

type Error struct {
	Message string `json:"message"`
}

//     required:
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

// Информация о пользователе.
//     required:
//      - fullname
//      - email
type User struct {
	Nickname string `json:"nickname"` //Имя пользователя (уникальное поле)
	//Данное поле допускает только латиницу, цифры и знак подчеркивания.
	//Сравнение имени регистронезависимо
	FullName string `json:"fullname"` // Полное имя пользователя
	About    string `json:"about"`    // Описание пользователя
	Email    string `json:"email"`    // Почтовый адрес пользователя (уникальное поле)
}

// Информация о форуме.
//    required:
//      - title
//      - user
//      - slug
type Forum struct {
	Title   string `json:"title"`   // Название форума
	User    string `json:"user"`    // Nickname пользователя, который отвечает за форум
	Slug    string `json:"slug"`    // Человекопонятный URL
	Posts   int    `json:"posts"`   // Общее кол-во сообщений в данном форуме
	Threads int    `json:"threads"` // Общее кол-во ветвей обсуждения в данном форуме
}

// Ветка обсуждения на форуме.
//     required:
//      - title
//      - author
//      - message
type Thread struct {
	ID      int       `json:"id"`      // Идентификатор ветки обсуждения.
	Title   string    `json:"title"`   // Заголовок ветки обсуждения.
	Author  string    `json:"author"`  // Пользователь, создавший данную тему.
	Forum   string    `json:"forum"`   // Форум, в котором расположена данная ветка обсуждения.
	Message string    `json:"message"` // Описание ветки обсуждения.
	Votes   int       `json:"votes"`   // Кол-во голосов непосредственно за данное сообщение форума.
	Slug    string    `json:"slug"`    // Человекопонятный URL
	Created time.Time `json:"created"` // Дата создания ветки на форуме.
}

//  Сообщение внутри ветки обсуждения на форуме.
//     required:
//      - author
//      - message
type Post struct { //  Сообщение внутри ветки обсуждения на форуме.
	ID       int       `json:"id"`       // Идентификатор данного сообщения.
	Parent   int       `json:"parent"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	Author   string    `json:"author"`   // Автор, написавший данное сообщение.
	Message  string    `json:"message"`  // Собственно сообщение форума.
	IsEdited bool      `json:"isEdited"` // Истина, если данное сообщение было изменено.
	Forum    string    `json:"forum"`    // Идентификатор форума (slug) данного сообещния.
	Thread   string    `json:"thread"`   // Идентификатор форума (slug) данного сообещния.
	Created  time.Time `json:"created"`  // Дата создания сообщения на форуме.
}

//  Сообщение для обновления сообщения внутри ветки на форуме. Пустые параметры остаются без изменений.
type PostUpdate struct {
	Message string `json:"message"` // Собственно сообщение форума.
}

//  Полная информация о сообщении, включая связанные объекты
type PostFull struct {
	Post   Post   `json:"post"`
	Author User   `json:"author"`
	Thread Thread `json:"thread"`
	Forum  Forum  `json:"forum"`
}

// Информация о голосовании пользователя.
//     required:
//      - nickname
//      - voice
type Vote struct {
	Nickname string `json:"nickname"` // Идентификатор пользователя.
	Voice    int32  `json:"voice"`    // Отданный голос.
}
