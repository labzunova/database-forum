-- CREATE EXTENSION IF NOT EXISTS CITEXT;
--
-- DROP TABLE IF EXISTS votes CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS posts;
-- DROP TABLE IF EXISTS forums CASCADE;
-- DROP TABLE IF EXISTS threads CASCADE;
-- DROP TABLE IF EXISTS forum_users CASCADE;
--
-- -- USERS --
-- CREATE UNLOGGED TABLE users (
--                                 id SERIAL NOT NULL PRIMARY KEY ,
--                                 nickname CITEXT COLLATE ucs_basic NOT NULL UNIQUE,
--                                 fullname TEXT,
--                                 about TEXT,
--                                 email CITEXT UNIQUE
-- );
-- CREATE INDEX IF NOT EXISTS users_full_but_email ON users (nickname, fullname, about, email); -- GetPostAuthor
-- CREATE INDEX IF NOT EXISTS users_nick_id ON users (nickname, id); -- get id ny nick
-- --CREATE INDEX IF NOT EXISTS user_nickname on users using hash(nickname); -- NEW
--
-- -- FORUMS --
-- CREATE UNLOGGED TABLE forums (
--                                  id SERIAL NOT NULL PRIMARY KEY,
--                                  title TEXT NOT NULL,
--     -- "user" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
--                                  "user" CITEXT NOT NULL,
--                                  FOREIGN KEY ("user") REFERENCES Users (nickname),
--                                  slug CITEXT UNIQUE NOT NULL,
--                                  created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
--                                  threads_count INTEGER DEFAULT 0,
--                                  posts_count INTEGER DEFAULT 0
-- );
-- CREATE INDEX IF NOT EXISTS forums_slug on forums using hash (slug); -- NEW
--
-- -- THREADS --
-- CREATE UNLOGGED TABLE threads (
--                                   id SERIAL NOT NULL PRIMARY KEY,
--                                   title TEXT NOT NULL,
--                                   slug CITEXT unique,
--                                   "author" CITEXT  NOT NULL,
--                                   FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
--     -- "forum" CITEXT REFERENCES forums(slug) ON DELETE CASCADE  NOT NULL,
--                                   forum CITEXT NOT NULL,
--                                   FOREIGN KEY (forum) REFERENCES forums (slug),
--                                   message TEXT NOT NULL,
--                                   votes INT DEFAULT 0 NOT NULL,
--                                   created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  NOT NULL
-- );
--
-- CREATE INDEX IF NOT EXISTS thread_forum_and_created ON threads (forum, created); -- для get forum threads
-- CREATE INDEX IF NOT EXISTS thread_slug ON threads (slug);
-- --CREATE INDEX IF NOT EXISTS thread_slug ON threads (lower(slug));
-- CREATE INDEX IF NOT EXISTS thread_forum ON threads USING hash (forum);
-- CREATE INDEX IF NOT EXISTS thread_created ON threads (created);
-- --CREATE INDEX IF NOT EXISTS thread_id_forum on threads(id, forum); -- NEW
--
-- --CREATE INDEX IF NOT EXISTS thread_slug_forum on threads(slug, forum); -- NEW
--
-- CREATE OR REPLACE FUNCTION new_thread_added() RETURNS TRIGGER AS
-- $new_thread_added$
-- begin
--     update forums set threads_count = threads_count + 1
--     where slug = new.forum;
--
--     return new;
-- end;
-- $new_thread_added$ LANGUAGE plpgsql;
-- create trigger new_thread_added
--     before insert on threads for each row
-- execute procedure new_thread_added();
--
-- -- POSTS --
-- CREATE UNLOGGED TABLE posts (
--                                 id SERIAL NOT NULL PRIMARY KEY,
--                                 parent INTEGER,
--                                 "author" CITEXT NOT NULL,
--                                 FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
--                                 message TEXT NOT NULL,
--                                 isEdited BOOLEAN NOT NULL DEFAULT FALSE,
--                                 "forum" CITEXT NOT NULL,
--                                 FOREIGN KEY ("forum") REFERENCES forums(slug) ON DELETE CASCADE,
--                                 "thread" INTEGER NOT NULL,-- ??? надо бы slug
--                                 FOREIGN KEY ("thread") REFERENCES threads(id) ON DELETE CASCADE,
--                                 created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
--                                 path INTEGER[] NOT NULL
-- );
-- CREATE INDEX IF NOT EXISTS post_thread ON posts (thread);
-- CREATE INDEX IF NOT EXISTS post_id_and_path ON posts (id, path);
-- CREATE INDEX IF NOT EXISTS post_id_and_thread ON posts (thread, id);
-- CREATE INDEX IF NOT EXISTS post_parent_and_id ON posts (parent, id);
-- CREATE INDEX IF NOT EXISTS post_path_and_thread ON posts (thread, path); -- запрос для сортировки
-- --CREATE INDEX IF NOT EXISTS post_id_and_thread ON posts (thread, id); -- !запрос для flat сортировки
-- CREATE INDEX IF NOT EXISTS post_path_parent_and_thread ON posts (thread, parent, path); -- для parentTree сортирвки
-- CREATE INDEX IF NOT EXISTS post_pathFirst_parent_and_thread ON posts (thread, (path[1]), id); -- для parentTree сортирвки
-- --CREATE INDEX IF NOT EXISTS post_id_path ON posts (id, path); -- !запрос для flat сортировки
--
-- CREATE OR REPLACE FUNCTION new_post_added() RETURNS TRIGGER AS
-- $new_post_added$
-- begin
--     update forums set posts_count = posts_count + 1
--     where slug = new.forum;
--
--     return new;
-- end;
-- $new_post_added$ LANGUAGE plpgsql;
-- create trigger new_post_added
--     before insert on posts for each row
-- execute procedure new_post_added();
--
--
-- CREATE OR REPLACE FUNCTION add_path() RETURNS TRIGGER AS
-- $add_path$
-- declare
--     parents INTEGER[];
-- begin
--     if (new.parent is null) then
--         new.path := new.path || new.id;
--     else
--         select path from posts where id = new.parent and thread = new.thread
--         into parents;
--
--         if (coalesce(array_length(parents, 1), 0) = 0) then
--             raise exception 'parent post not exists';
--         end if;
--
--         new.path := new.path || parents || new.id;
--     end if;
--     return new;
-- end;
-- $add_path$ LANGUAGE plpgsql;
-- create trigger add_path
--     before insert on posts for each row
-- execute procedure add_path();
--
--
-- -- FORUM USERS --
-- CREATE UNLOGGED TABLE forum_users (
--                                       fullname TEXT,
--     --  fullname TEXT,
--                                       about TEXT,
--                                       email CITEXT,
--                                       userNickname CITEXT COLLATE ucs_basic,
--
--     -- userNickname CITEXT REFERENCES users (nickname),
--                                       FOREIGN KEY (userNickname) REFERENCES users (nickname),
--                                       forumSlug CITEXT, -- изменила из-за GetUsers
--                                       FOREIGN KEY (forumSlug) REFERENCES forums (slug),
--
--                                       unique (userNickname, forumSlug)
-- );
-- --CREATE INDEX IF NOT EXISTS forum_users ON forum_users (forumSlug, userNickname);
-- CREATE INDEX IF NOT EXISTS forum_users_full ON forum_users (forumSlug, userNickname, fullname, about, email);
-- --CREATE INDEX IF NOT EXISTS forum_users_slug ON forum_users (forumSlug);
-- --CREATE INDEX IF NOT EXISTS forum_users_slug ON forum_users (forumSlug);
--
-- DROP FUNCTION IF EXISTS new_forum_user_added() CASCADE;
-- CREATE OR REPLACE FUNCTION new_forum_user_added() RETURNS TRIGGER AS
-- $new_forum_user_added$
-- begin
--     declare
--         nickAuthor citext;
--         fullnameAuthor text;
--         emailAuthor citext;
--         aboutAuthor text;
--     begin
--         select nickname, fullname, about, email
--         from users where nickname = new.author
--         into nickAuthor, fullnameAuthor, aboutAuthor, emailAuthor;
--
--         insert into forum_users(fullname, about, email, userNickname, forumSlug)
--         values (fullnameAuthor, aboutAuthor, emailAuthor, nickAuthor, new.forum)
--         on conflict do nothing;
--
--         return null;
--     end;
-- end;
-- $new_forum_user_added$ LANGUAGE plpgsql;
-- drop trigger if exists new_forum_user_added on posts;
-- create trigger new_forum_user_added
--     AFTER insert on posts for each row
-- execute procedure new_forum_user_added();
-- drop trigger if exists new_forum_user_added on threads;
-- create trigger new_forum_user_added
--     AFTER insert on threads for each row
-- execute procedure new_forum_user_added();
--
-- -- VOTES --
-- CREATE UNLOGGED TABLE votes (
--                                 "user" CITEXT,
--                                 FOREIGN KEY ("user") REFERENCES users (nickname),
--                                 thread integer,
--                                 FOREIGN KEY (thread)  REFERENCES threads(id),
--                                 vote INTEGER,
--                                 UNIQUE (thread, "user")
-- );
-- --CREATE INDEX IF NOT EXISTS votes_full ON votes (thread, "user", vote);
--
-- DROP FUNCTION IF EXISTS new_vote() CASCADE;
-- CREATE OR REPLACE FUNCTION new_vote() RETURNS TRIGGER AS
-- $new_vote$
-- begin
--     update threads set votes = votes + new.vote
--     where id = new.thread;
--
--     return null;
-- end;
-- $new_vote$ LANGUAGE plpgsql;
-- create trigger new_vote
--     AFTER insert on votes for each row
-- execute procedure new_vote();
--
-- DROP FUNCTION IF EXISTS change_vote() CASCADE;
-- CREATE OR REPLACE FUNCTION change_vote() RETURNS TRIGGER AS
-- $change_vote$
-- begin
--     update threads set votes = votes - old.vote + new.vote
--     where id = new.thread;
--
--     return null;
-- end;
-- $change_vote$ LANGUAGE plpgsql;
-- create trigger change_vote
--     AFTER update on votes for each row
-- execute procedure change_vote();
--
-- ANALYZE;
-- VACUUM ANALYZE;

CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

-- USERS --
CREATE UNLOGGED TABLE users (
                                nickname CITEXT COLLATE ucs_basic primary key NOT NULL UNIQUE,
                                fullname TEXT,
                                about TEXT,
                                email CITEXT UNIQUE
);
CREATE INDEX IF NOT EXISTS users_full_but_id ON users (nickname, fullname, about, email); -- GetPostAuthor
--CREATE INDEX IF NOT EXISTS users_nick_id ON users (nickname, id); -- get id ny nick
-- --CREATE INDEX IF NOT EXISTS user_nickname on users using hash(nickname); -- NEW
CREATE INDEX IF NOT EXISTS user_nickname on users using hash(nickname); -- NEW

-- FORUMS --
CREATE UNLOGGED TABLE forums (
                                 title TEXT NOT NULL,
                                 "user" CITEXT NOT NULL,
                                 slug CITEXT primary key UNIQUE NOT NULL,
    -- "user" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                                 threads_count INTEGER DEFAULT 0,
                                 posts_count INTEGER DEFAULT 0,
                                 FOREIGN KEY ("user") REFERENCES Users (nickname),
                                 created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
--CREATE INDEX IF NOT EXISTS forums_slug on forums using hash (slug); -- NEW

-- THREADS --
CREATE UNLOGGED TABLE threads (
                                  id SERIAL NOT NULL PRIMARY KEY,
                                  slug CITEXT unique,
                                  forum CITEXT NOT NULL,
                                  "author" CITEXT  NOT NULL,
                                  title TEXT NOT NULL,
                                  message TEXT NOT NULL,
                                  votes INT DEFAULT 0 NOT NULL,
    -- "forum" CITEXT REFERENCES forums(slug) ON DELETE CASCADE  NOT NULL,
                                  created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  NOT NULL,
                                  FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
                                  FOREIGN KEY (forum) REFERENCES forums (slug)
);
CREATE INDEX IF NOT EXISTS thread_forum_and_created ON threads (forum, created); -- для get forum threads
CREATE INDEX IF NOT EXISTS thread_forum ON threads (forum); -- для get forum threads
CREATE INDEX IF NOT EXISTS thread_slug ON threads (slug) where slug != '';
--CREATE INDEX IF NOT EXISTS thread_slug ON threads (lower(slug));
--CREATE INDEX IF NOT EXISTS thread_forum ON threads USING hash (forum);
CREATE INDEX IF NOT EXISTS thread_created ON threads (created);
--CREATE INDEX IF NOT EXISTS threads_full ON threads (forum, slug, created, title, author, message, votes);
--CREATE INDEX IF NOT EXISTS thread_id_forum on threads(id, forum); -- NEW
--CREATE INDEX IF NOT EXISTS thread_slug_forum on threads(slug, forum); -- NEW

CREATE OR REPLACE FUNCTION new_thread_added() RETURNS TRIGGER AS
$new_thread_added$
begin
    update forums set threads_count = threads_count + 1
    where slug = new.forum;

    return new;
end;
$new_thread_added$ LANGUAGE plpgsql;
create trigger new_thread_added
    before insert on threads for each row
execute procedure new_thread_added();

-- POSTS --
CREATE UNLOGGED TABLE posts (
                                id SERIAL NOT NULL PRIMARY KEY,
                                "thread" INTEGER NOT NULL,
                                "author" CITEXT NOT NULL,
                                "forum" CITEXT NOT NULL,
                                isEdited BOOLEAN NOT NULL DEFAULT FALSE,
                                message TEXT NOT NULL,
                                parent INTEGER,

                                created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                path INTEGER[] NOT NULL,
                                FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
                                FOREIGN KEY ("forum") REFERENCES forums(slug) ON DELETE CASCADE,
                                FOREIGN KEY ("thread") REFERENCES threads(id) ON DELETE CASCADE
);
-- CREATE INDEX IF NOT EXISTS post_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS post_id_created_thread ON posts (id, created, thread);
-- --CREATE INDEX IF NOT EXISTS post_id_and_path ON posts ((path[1]), id);
CREATE INDEX IF NOT EXISTS post_parent_and_id ON posts (parent, id);
CREATE INDEX IF NOT EXISTS post_id_path ON posts (id, path);
-- -- CREATE INDEX IF NOT EXISTS post_thread_parent ON posts (thread, parent); --
-- CREATE INDEX IF NOT EXISTS post_thread_path1 ON posts (thread, (path[1])); --SUPERNEW MY
--CREATE INDEX IF NOT EXISTS post_thread_path ON posts (thread, path); --SUPERNEW MY
-- CREATE INDEX IF NOT EXISTS post_parent_thread_path1 ON posts (parent, thread, (path[1])); --SUPERNEW MY
CREATE INDEX IF NOT EXISTS post_id_and_thread ON posts (thread, id); -- !запрос для flat сортировки
CREATE INDEX IF NOT EXISTS post_thread_and_path ON posts (thread, path);
CREATE INDEX IF NOT EXISTS post_path_parent_and_thread ON posts (thread, parent, path); -- для parentTree сортирвки
CREATE INDEX IF NOT EXISTS post_pathFirst_parent_and_thread ON posts (thread, (path[1]), id); -- для parentTree сортирвки
-- -- --CREATE INDEX IF NOT EXISTS post_id_path ON posts (id, path); -- !запрос для flat сортировки

CREATE OR REPLACE FUNCTION new_post_added() RETURNS TRIGGER AS
$new_post_added$
begin
    update forums set posts_count = posts_count + 1
    where slug = new.forum;

    return new;
end;
$new_post_added$ LANGUAGE plpgsql;
create trigger new_post_added
    before insert on posts for each row
execute procedure new_post_added();


CREATE OR REPLACE FUNCTION add_path() RETURNS TRIGGER AS
$add_path$
declare
    parents INTEGER[];
begin
    if (new.parent is null) then
        new.path := new.path || new.id;
    else
        select path from posts where id = new.parent and thread = new.thread
        into parents;

        if (coalesce(array_length(parents, 1), 0) = 0) then
            raise exception 'parent post not exists';
        end if;

        new.path := new.path || parents || new.id;
    end if;
    return new;
end;
$add_path$ LANGUAGE plpgsql;
create trigger add_path
    before insert on posts for each row
execute procedure add_path();


-- FORUM USERS --
CREATE UNLOGGED TABLE forum_users (
                                      fullname TEXT,
    --  fullname TEXT,
                                      about TEXT,
                                      email CITEXT,
                                      userNickname CITEXT COLLATE ucs_basic,

    -- userNickname CITEXT REFERENCES users (nickname),
                                      FOREIGN KEY (userNickname) REFERENCES users (nickname),
                                      forumSlug CITEXT, -- изменила из-за GetUsers
                                      FOREIGN KEY (forumSlug) REFERENCES forums (slug),

                                      unique (userNickname, forumSlug)
);
-- КАК ОПТИМИЗИРОВАТЬ ЭТУ ТАБЛИЦУ АЛО
--CREATE INDEX IF NOT EXISTS forum_users_nickname_forumslug ON forum_users using hash(forumSlug);
-- CREATE INDEX IF NOT EXISTS forum_users ON forum_users (forumSlug, userNickname);
-- CREATE INDEX IF NOT EXISTS forum_users_full ON forum_users (forumSlug, userNickname, fullname, about, email);

DROP FUNCTION IF EXISTS new_forum_user_added() CASCADE;
CREATE OR REPLACE FUNCTION new_forum_user_added() RETURNS TRIGGER AS
$new_forum_user_added$
begin
    declare
        nickAuthor citext;
        fullnameAuthor text;
        emailAuthor citext;
        aboutAuthor text;
    begin
        select nickname, fullname, about, email
        from users where nickname = new.author
        into nickAuthor, fullnameAuthor, aboutAuthor, emailAuthor;

        insert into forum_users(fullname, about, email, userNickname, forumSlug)
        values (fullnameAuthor, aboutAuthor, emailAuthor, nickAuthor, new.forum)
        on conflict do nothing;

        return null;
    end;
end;
$new_forum_user_added$ LANGUAGE plpgsql;
drop trigger if exists new_forum_user_added on posts;
create trigger new_forum_user_added
    AFTER insert on posts for each row
execute procedure new_forum_user_added();
drop trigger if exists new_forum_user_added on threads;
create trigger new_forum_user_added
    AFTER insert on threads for each row
execute procedure new_forum_user_added();

-- VOTES --
CREATE UNLOGGED TABLE votes (
                                "user" CITEXT,
                                FOREIGN KEY ("user") REFERENCES users (nickname),
                                thread integer,
                                FOREIGN KEY (thread)  REFERENCES threads(id),
                                vote INTEGER,
                                UNIQUE (thread, "user")
);
--CREATE INDEX IF NOT EXISTS votes_full ON votes (thread, "user", vote);
CREATE INDEX IF NOT EXISTS votes_full ON votes (thread, "user");

DROP FUNCTION IF EXISTS new_vote() CASCADE;
CREATE OR REPLACE FUNCTION new_vote() RETURNS TRIGGER AS
$new_vote$
begin
    update threads set votes = votes + new.vote
    where id = new.thread;

    return null;
end;
$new_vote$ LANGUAGE plpgsql;
create trigger new_vote
    AFTER insert on votes for each row
execute procedure new_vote();

DROP FUNCTION IF EXISTS change_vote() CASCADE;
CREATE OR REPLACE FUNCTION change_vote() RETURNS TRIGGER AS
$change_vote$
begin
    update threads set votes = votes - old.vote + new.vote
    where id = new.thread;

    return null;a
end;
$change_vote$ LANGUAGE plpgsql;
create trigger change_vote
    AFTER update on votes for each row
execute procedure change_vote();


ANALYZE;
VACUUM ANALYZE;
-- VACUUM;
-- VACUUM ANALYSE;