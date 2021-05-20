DROP TABLE IF EXISTS users CASCADE ;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE users (
                       id SERIAL NOT NULL PRIMARY KEY,
                       nickname TEXT NOT NULL UNIQUE,
                       fullname TEXT,
                       about TEXT,
                       email TEXT
);

CREATE TABLE forums (
                        id SERIAL NOT NULL PRIMARY KEY,
                        title TEXT UNIQUE NOT NULL,
                        "user" TEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                        slug TEXT NOT NULL,
                        posts INT,
                        threads INT
);

CREATE TABLE threads (
                         id SERIAL NOT NULL PRIMARY KEY,
                         title TEXT NOT NULL,
                         "author" INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                         "forum" TEXT REFERENCES forums(title) ON DELETE CASCADE  NOT NULL,
                         message TEXT,
                         votes INT,
                         slug TEXT,
                         created TIMESTAMP
);

CREATE TABLE posts (
                       id SERIAL NOT NULL PRIMARY KEY,
                       parent INTEGER DEFAULT 0 NOT NULL,
                       "author" TEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                       message TEXT NOT NULL,
                       isEdited BOOLEAN NOT NULL,
                       "forum" INTEGER REFERENCES forums(id) ON DELETE CASCADE NOT NULL,
                       "thread" INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,
                       created TIMESTAMP NOT NULL
);

CREATE TABLE forum_users (
    userID  INTEGER REFERENCES users (id),
    forumID INTEGER REFERENCES forums (id)
);
