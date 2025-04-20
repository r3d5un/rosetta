from python.db.filters import Filter, Metadata
from python.db.forum import Forum
from python.db.post import Post, PostPatch
from python.db.session import Models
from python.db.thread import Thread
from python.db.user import User
from tests.conftest import get_testcontainer_db_engine


def test_insert():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")


def test_select():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    try:
        post = models.posts.select(post.id)
    except Exception as e:
        raise ValueError(f"error upon selecting post: {e}")
    assert post is not None


def test_select_all():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    result: tuple[list[Post] | None, Metadata | None] | None = models.posts.select_all(
        Filter(page_size=100)
    )
    if result is None:
        raise ValueError("no results returned")
    posts = result[0]
    metadata = result[1]

    if posts is None:
        raise ValueError("posts is None")
    assert len(posts) > 0
    if metadata is None:
        raise ValueError("metadata is None")
    assert metadata.response_length == len(posts)


def test_update():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    content = "Johnny's car is at the docks"
    patch = PostPatch(thread_id=thread.id, content=content, id=post.id)
    try:
        post = models.posts.update(patch)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")


def test_soft_delete():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    try:
        post = models.posts.soft_delete(post.id)
    except Exception as e:
        raise Exception(f"error upon soft deleting post: {e}")
    if post is None:
        raise ValueError("no post returned upon soft deletion")
    if post.deleted is None:
        raise ValueError("inserted post deleted status is None")
    assert post.deleted is True


def test_restore():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    try:
        post = models.posts.restore(post.id)
    except Exception as e:
        raise Exception(f"error upon restoring post: {e}")
    if post is None:
        raise ValueError("no post returned upon restoring")
    if post.deleted is None:
        raise ValueError("inserted post deleted status is None")
    assert post.deleted is False


def test_delete():
    models = Models(engine=get_testcontainer_db_engine())
    user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if user is None:
        raise ValueError("no user returnd upon insertion")
    if user.id is None:
        raise ValueError("inserted user ID is None")

    forum = Forum(owner_id=user.id, name="Crushing Militech", description="")
    try:
        inserted_forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if inserted_forum is None:
        raise ValueError("no forum returnd upon insertion")
    if inserted_forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=inserted_forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    post = Post(
        thread_id=thread.id, reply_to=None, content="Porche 911", author_id=user.id
    )
    try:
        post = models.posts.insert(post)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
    if post.id is None:
        raise ValueError("inserted post ID is None")

    try:
        post = models.posts.delete(post.id)
    except Exception as e:
        raise Exception(f"error upon inserting post: {e}")
    if post is None:
        raise ValueError("no post returned upon insertion")
