from python.db.user import User, UserModel
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


def test_something():
    with get_connection().cursor() as c:
        c.execute("SELECT * FROM forum.users")
        users = c.fetchall()
        print(f"Users: {users}")
