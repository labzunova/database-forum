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
                        title TEXT NOT NULL,
                        "user" TEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                        slug TEXT UNIQUE NOT NULL,
                        posts INT DEFAULT 0,
                        threads INT DEFAULT 0,
                        created TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
);

CREATE TABLE threads (
                         id SERIAL NOT NULL PRIMARY KEY,
                         title TEXT NOT NULL,
                         slug TEXT,
                         "author" INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
                         "forum" TEXT REFERENCES forums(slug) ON DELETE CASCADE  NOT NULL,
                         message TEXT UNIQUE,
                         votes INT DEFAULT 0,
                         created TIMESTAMP
);

CREATE TABLE posts (
                       id SERIAL NOT NULL PRIMARY KEY,
                       parent INTEGER DEFAULT 0 NOT NULL,
                       "author" TEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                       message TEXT NOT NULL,
                       isEdited BOOLEAN NOT NULL,
                       "forum" TEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
                       "thread" TEXT REFERENCES threads(slug) ON DELETE CASCADE NOT NULL,
                       created TIMESTAMP NOT NULL
);

CREATE TABLE forum_users (
    userID  INTEGER REFERENCES users (id),
    forumSlug TEXT REFERENCES forums (slug) -- изменила из-за GetUsers
);
