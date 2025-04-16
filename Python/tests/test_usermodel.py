from python.db.user import User, UserModel
from tests.conftest import get_connection, get_testcontainer_db_engine

test_user = User(
    name="Johnny Silverhand", username="samurai", email="jsilverhand@samurai.com"
)


def test_insert():
    user_model = UserModel(get_testcontainer_db_engine())
    print(f"inserting user: {test_user}")
    inserted_user = user_model.insert(test_user)
    print(f"inserted user: {inserted_user}")
    assert inserted_user is not None


def test_something():
    with get_connection().cursor() as c:
        c.execute("SELECT * FROM forum.users")
        users = c.fetchall()
        print(f"Users: {users}")
