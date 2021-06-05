--CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

CREATE TABLE users (
                       id SERIAL NOT NULL PRIMARY KEY,
                       nickname CITEXT NOT NULL UNIQUE,
                       fullname TEXT,
                       about TEXT,
                       email CITEXT UNIQUE
);

CREATE TABLE forums (
                        id SERIAL NOT NULL PRIMARY KEY,
                        title TEXT NOT NULL,
    -- "user" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                        "user" CITEXT NOT NULL,
                        FOREIGN KEY ("user") REFERENCES Users (nickname),
                        slug CITEXT UNIQUE NOT NULL,
                        posts INT DEFAULT 0,
                        threads INT DEFAULT 0,
                        created TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);

CREATE TABLE threads (
                         id SERIAL NOT NULL PRIMARY KEY,
                         title TEXT NOT NULL,
                         slug CITEXT unique, -- ????????????????
                         "author" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    -- "forum" CITEXT REFERENCES forums(slug) ON DELETE CASCADE  NOT NULL,
                         forum CITEXT NOT NULL,
                         FOREIGN KEY (forum) REFERENCES forums (slug),
                         message TEXT NOT NULL,
                         votes INT DEFAULT 0 NOT NULL,
                         created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE posts (
                       id SERIAL NOT NULL PRIMARY KEY,
                       parent INTEGER DEFAULT 0 NOT NULL,
                       "author" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                       message TEXT NOT NULL,
                       isEdited BOOLEAN NOT NULL DEFAULT FALSE,
                       "forum" CITEXT REFERENCES forums(slug) ON DELETE CASCADE NOT NULL,
                       "thread" INTEGER REFERENCES threads(id) ON DELETE CASCADE NOT NULL,-- ??? надо бы slug
                       created TIMESTAMP NOT NULL
);

CREATE TABLE forum_users (
                             userID  INTEGER REFERENCES users (id),
                             forumSlug CITEXT REFERENCES forums (slug) -- изменила из-за GetUsers
);

CREATE TABLE votes (
                       "user" CITEXT REFERENCES users(nickname), -- nickname?
                       thread INTEGER REFERENCES threads(id), -- ??? надо бы slug
                       vote INTEGER,
                       UNIQUE (thread, "user")
);
