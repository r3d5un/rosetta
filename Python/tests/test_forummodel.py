from python.db.filters import Filter, Metadata
from python.db.forum import Forum, ForumPatch
from python.db.session import Models
from python.db.user import User
from tests.conftest import get_testcontainer_db_engine


def test_insert():
    models = Models(engine=get_testcontainer_db_engine())
    test_user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(test_user)
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


def test_select():
    models = Models(engine=get_testcontainer_db_engine())
    test_user = User(
        name="Saburo Arasaka", username="s.arasaka", email="s.arasaka@arasaka.com"
    )
    try:
        user = models.users.insert(test_user)
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

    selected_forum = models.forums.select(inserted_forum.id)
    if selected_forum is None:
        raise ValueError("no user returnd upon insertion")
    assert selected_forum.name == forum.name


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

    result: tuple[list[Forum] | None, Metadata | None] | None = (
        models.forums.select_all(Filter(page_size=100))
    )
    if result is None:
        raise ValueError("no results returned")
    forums = result[0]
    metadata = result[1]

    if forums is None:
        raise ValueError("forums is None")
    assert len(forums) > 0
    if metadata is None:
        raise ValueError("metadata is None")
    assert metadata.response_length == len(forums)


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

    description = "test description"
    update = ForumPatch(id=inserted_forum.id, description=description)
    try:
        updated_forum = models.forums.update(update)
    except Exception as e:
        raise ValueError(f"unable to update forum: {e}")
    if updated_forum is None:
        raise ValueError("inserted username is None")
    if updated_forum.description is None:
        raise ValueError("inserted username is None")
    assert updated_forum.description == description


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
        forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if forum is None:
        raise ValueError("no forum returnd upon insertion")
    if forum.id is None:
        raise ValueError("inserted forum ID is None")

    try:
        deleted_forum = models.forums.soft_delete(forum.id)
    except Exception as e:
        raise ValueError(f"unable to update forum: {e}")
    if deleted_forum is None:
        raise ValueError("deleted forum is None")
    assert deleted_forum.deleted is True


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
        forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if forum is None:
        raise ValueError("no forum returnd upon insertion")
    if forum.id is None:
        raise ValueError("inserted forum ID is None")

    try:
        forum = models.forums.restore(forum.id)
    except Exception as e:
        raise ValueError(f"unable to restore forum: {e}")
    if forum is None:
        raise ValueError("restored forum is None")
    assert forum.deleted is False


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
        forum = models.forums.insert(forum)
    except Exception as e:
        raise Exception(f"error upon inserting forum: {e}")
    if forum is None:
        raise ValueError("no forum returnd upon insertion")
    if forum.id is None:
        raise ValueError("inserted forum ID is None")

    try:
        forum = models.forums.delete(forum.id)
    except Exception as e:
        raise ValueError(f"unable to delete forum: {e}")
    if forum is None:
        raise ValueError("deleted forum is None")
