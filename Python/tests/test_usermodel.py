from tests.conftest import get_connection


def test_something():
    with get_connection().cursor() as c:
        c.execute("SELECT * FROM forum.users")
        users = c.fetchall()
        print(f"Users: {users}")
