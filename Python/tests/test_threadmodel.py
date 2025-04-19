from python.db.filters import Filter, Metadata
from python.db.forum import Forum
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
        forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if forum is None:
        raise ValueError("no forum returnd upon insertion")
    if forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    thread = models.threads.select(thread.id)
    if thread is None:
        raise ValueError("no user returnd upon insertion")
    assert thread.title == thread.title


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
        forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if forum is None:
        raise ValueError("no forum returnd upon insertion")
    if forum.id is None:
        raise ValueError("inserted forum ID is None")

    thread = Thread(forum_id=forum.id, author_id=user.id, title="Johnny Boy")
    try:
        thread = models.threads.insert(thread)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")

    result: tuple[list[Thread] | None, Metadata | None] | None = (
        models.threads.select_all(Filter(page_size=100))
    )
    if result is None:
        raise ValueError("no results returned")
    threads = result[0]
    metadata = result[1]

    if threads is None:
        raise ValueError("threads is None")
    assert len(threads) > 0
    if metadata is None:
        raise ValueError("metadata is None")
    assert metadata.response_length == len(threads)
