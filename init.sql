CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;

-- USERS --
CREATE UNLOGGED TABLE users (
                                id SERIAL NOT NULL PRIMARY KEY ,
                                nickname CITEXT COLLATE "POSIX" NOT NULL UNIQUE,
                                fullname TEXT,
                                about TEXT,
                                email CITEXT UNIQUE
);
CREATE INDEX IF NOT EXISTS users_full ON users (nickname, fullname, about, email);

-- FORUMS --
CREATE UNLOGGED TABLE forums (
                                 id SERIAL NOT NULL PRIMARY KEY,
                                 title TEXT NOT NULL,
    -- "user" CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
                                 "user" CITEXT NOT NULL,
                                 FOREIGN KEY ("user") REFERENCES Users (nickname),
                                 slug CITEXT UNIQUE NOT NULL,
                                 created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                 threads_count INTEGER DEFAULT 0,
                                 posts_count INTEGER DEFAULT 0
);

-- THREADS --
CREATE UNLOGGED TABLE threads (
                                  id SERIAL NOT NULL PRIMARY KEY,
                                  title TEXT NOT NULL,
                                  slug CITEXT unique,
                                  "author" CITEXT  NOT NULL,
                                  FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
    -- "forum" CITEXT REFERENCES forums(slug) ON DELETE CASCADE  NOT NULL,
                                  forum CITEXT NOT NULL,
                                  FOREIGN KEY (forum) REFERENCES forums (slug),
                                  message TEXT NOT NULL,
                                  votes INT DEFAULT 0 NOT NULL,
                                  created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  NOT NULL
);

CREATE INDEX IF NOT EXISTS thread_forum_and_created ON threads (forum, created); -- для get forum threads

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
                                parent INTEGER,
                                "author" CITEXT NOT NULL,
                                FOREIGN KEY ("author") REFERENCES users(nickname) ON DELETE CASCADE,
                                message TEXT NOT NULL,
                                isEdited BOOLEAN NOT NULL DEFAULT FALSE,
                                "forum" CITEXT NOT NULL,
                                FOREIGN KEY ("forum") REFERENCES forums(slug) ON DELETE CASCADE,
                                "thread" INTEGER NOT NULL,-- ??? надо бы slug
                                FOREIGN KEY ("thread") REFERENCES threads(id) ON DELETE CASCADE,
                                created TIMESTAMP(3) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                path INTEGER[] NOT NULL
);

CREATE INDEX IF NOT EXISTS post_id_and_thread ON posts (thread, id);
CREATE INDEX IF NOT EXISTS post_path_and_thread ON posts (thread, path); -- запрос для сортировка
CREATE INDEX IF NOT EXISTS post_path_first_and_thread ON posts (thread, path); -- запрос для сортировка
CREATE INDEX IF NOT EXISTS post_path_parent_and_thread ON posts (thread, path, parent); -- для parentTree сортирвки

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
                                      userNickname CITEXT REFERENCES users (nickname),
                                      FOREIGN KEY (userNickname) REFERENCES users (nickname),
                                      forumSlug CITEXT REFERENCES forums (slug), -- изменила из-за GetUsers
                                      FOREIGN KEY (forumSlug) REFERENCES forums (slug),

                                      unique (userNickname, forumSlug)
);

DROP FUNCTION IF EXISTS new_forum_user_added() CASCADE;
CREATE OR REPLACE FUNCTION new_forum_user_added() RETURNS TRIGGER AS
$new_forum_user_added$
begin
    insert into forum_users (userNickname, forumSlug)
    values (new.author, new.forum) on conflict do nothing;

    return null;
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

    return null;
end;
$change_vote$ LANGUAGE plpgsql;
create trigger change_vote
    AFTER update on votes for each row
execute procedure change_vote();
