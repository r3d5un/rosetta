from python.db.user import User, UserModel, UserPatch
from src.python.db.filters import Filter, Metadata
from tests.conftest import get_connection, get_testcontainer_db_engine

test_user = User(
    name="Johnny Silverhand", username="samurai", email="jsilverhand@samurai.com"
)


def test_insert():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    print(inserted_user)


def test_select():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    selected_user = user_model.select(inserted_user.id)
    if selected_user is None:
        raise ValueError("no user returnd upon insertion")
    assert selected_user.name == inserted_user.name
    assert selected_user.username == inserted_user.username
    assert selected_user.email == inserted_user.email


def test_select_all():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    result: tuple[list[User] | None, Metadata | None] | None = user_model.select_all(
        Filter(page_size=100)
    )
    if result is None:
        raise ValueError("no results returned")
    users = result[0]
    metadata = result[1]

    if users is None:
        raise ValueError("users is None")
    assert len(users) > 0
    if metadata is None:
        raise ValueError("metadata is None")
    assert metadata.response_length == len(users)


def test_update():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    update = UserPatch(id=inserted_user.id, username="johnnyboy")
    try:
        updated_user = user_model.update(update)
    except Exception as e:
        raise ValueError(f"unable to update user: {e}")
    if updated_user is None:
        raise ValueError("inserted username is None")
    if updated_user.username is None:
        raise ValueError("inserted username is None")
    assert update.username == updated_user.username
    assert inserted_user.email == updated_user.email
    assert inserted_user.name == updated_user.name


def test_soft_delete():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    try:
        deleted_user = user_model.soft_delete(inserted_user.id)
    except Exception as e:
        raise ValueError(f"unable to update user: {e}")
    if deleted_user is None:
        raise ValueError("deleted user is None")
    assert deleted_user.deleted is True


def test_restore():
    user_model = UserModel(get_testcontainer_db_engine())

    try:
        inserted_user = user_model.insert(test_user)
    except Exception as e:
        raise Exception(f"error upon inserting user: {e}")
    if inserted_user is None:
        raise ValueError("no user returnd upon insertion")
    if inserted_user.id is None:
        raise ValueError("inserted user ID is None")

    try:
        deleted_user = user_model.restore(inserted_user.id)
    except Exception as e:
        raise ValueError(f"unable to update user: {e}")
    if deleted_user is None:
        raise ValueError("deleted user is None")
    assert deleted_user.deleted is False


def test_something():
    with get_connection().cursor() as c:
        c.execute("SELECT * FROM forum.users")
        users = c.fetchall()
        print(f"Users: {users}")
