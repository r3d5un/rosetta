import os
import typing

import psycopg
import psycopg.sql as sql
import pytest
from testcontainers.postgres import PostgresContainer

postgres = PostgresContainer("postgres:17.4")


@pytest.fixture(scope="module", autouse=True)
def setup(request):
    postgres.start()

    def remove_container():
        postgres.stop()

    request.addfinalizer(remove_container)
    os.environ["DB_CONN"] = postgres.get_connection_url()
    os.environ["DB_HOST"] = postgres.get_container_host_ip()
    os.environ["DB_PORT"] = str(postgres.get_exposed_port(5432))
    os.environ["DB_USERNAME"] = postgres.username
    os.environ["DB_PASSWORD"] = postgres.password
    os.environ["DB_NAME"] = postgres.dbname


@pytest.fixture(scope="function", autouse=True)
def up_migrations():
    pass


def list_up_migrations(dir_path: str) -> list[str]:
    migrations = []
    try:
        files = os.listdir(dir_path)
    except OSError as e:
        raise e

    for file in files:
        if not os.path.isdir(os.path.join(dir_path, file)) and file.endswith(".up.sql"):
            migrations.append(os.path.join(dir_path, file))

    return migrations


def find_project_root():
    cwd = os.getcwd()
    marker_file = ".git"

    while True:
        if os.path.exists(os.path.join(cwd, marker_file)):
            print(cwd)
            return cwd
        parent = os.path.dirname(cwd)
        if parent == cwd:
            raise FileNotFoundError("project root not found")
        cwd = parent


def test_something():
    conn = get_connection()

    migrations_dir = f"{find_project_root()}/migrations"
    print(f"Looking for migrations in: {migrations_dir}")
    migrations = list_up_migrations(migrations_dir)

    for migration in migrations:
        print(f"Performing up migration: {migration}")
        with open(migration) as f:
            migration_sql = f.read()
            with conn.cursor() as c:
                c.execute(typing.cast(typing.LiteralString, migration_sql))
                conn.commit()

    with conn.cursor() as c:
        c.execute("SELECT * FROM forum.users")
        users = c.fetchall()
        print(f"Users: {users}")


def get_connection() -> psycopg.Connection:
    host = os.getenv("DB_HOST", "localhost")
    port = os.getenv("DB_PORT", "5432")
    username = os.getenv("DB_USERNAME", "postgres")
    password = os.getenv("DB_PASSWORD", "postgres")
    database = os.getenv("DB_NAME", "postgres")

    return psycopg.connect(
        f"host={host} dbname={database} user={username} password={password} port={port}"
    )
