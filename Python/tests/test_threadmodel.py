from python.db.filters import Filter, Metadata
from python.db.forum import Forum
from python.db.session import Models
from python.db.thread import Thread, ThreadPatch
from python.db.threadvote import ThreadVote
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

    new_title = "[Update] Johnny Boy"
    thread_patch = ThreadPatch(id=thread.id, title=new_title)
    try:
        thread = models.threads.update(thread_patch)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.id is None:
        raise ValueError("inserted thread ID is None")
    assert thread.title == new_title


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

    try:
        thread = models.threads.soft_delete(thread.id)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.deleted is None:
        raise ValueError("inserted thread deleted status is None")
    assert thread.deleted is True


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

    try:
        thread = models.threads.restore(thread.id)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")
    if thread.deleted is None:
        raise ValueError("inserted thread deleted status is None")
    assert thread.deleted is False


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

    try:
        thread = models.threads.delete(thread.id)
    except Exception as e:
        raise Exception(f"error upon inserting thread: {e}")
    if thread is None:
        raise ValueError("no thread returned upon insertion")


def test_thread_votes():
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

    try:
        vote = models.thread_votes.vote(
            ThreadVote(thread_id=thread.id, user_id=user.id, vote=1)
        )
        if vote is None:
            raise Exception("vote is None")
    except Exception as e:
        raise Exception(f"error occurred when voting: {e}")

    try:
        vote_sum = models.thread_votes.select_count(
            filters=Filter(thread_id=thread.id, user_id=user.id)
        )
    except Exception as e:
        raise Exception(f"error occurred when getting thread vote sum: {e}")
    if vote_sum is None:
        raise Exception("vote sum is None")
    assert vote_sum >= 1
