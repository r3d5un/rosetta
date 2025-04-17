import os
import typing

import psycopg
import pytest
from sqlalchemy import Engine
from sqlmodel import create_engine
from testcontainers.postgres import PostgresContainer

postgres = PostgresContainer("postgres:17.4")


@pytest.fixture(scope="module", autouse=True)
def setup_database(request):
    postgres.start()

    def remove_container():
        postgres.stop()

    request.addfinalizer(remove_container)


@pytest.fixture(scope="function")
def db_connection():
    return get_connection()


def get_connection() -> psycopg.Connection:
    host = postgres.get_container_host_ip()
    port = str(postgres.get_exposed_port(5432))
    username = postgres.username
    password = postgres.password
    database = postgres.dbname

    return psycopg.connect(
        f"host={host} dbname={database} user={username} password={password} port={port}"
    )


@pytest.fixture(scope="function", autouse=True)
def apply_migrations(db_connection):
    migrations_dir = f"{find_project_root()}/migrations"
    up_migrations = list_up_migrations(migrations_dir)
    for migration in up_migrations:
        print(f"Performing up migration: {migration}")
        with open(migration) as f:
            migration_sql = f.read()
            with db_connection.cursor() as c:
                c.execute(typing.cast(typing.LiteralString, migration_sql))
                db_connection.commit()
    print("Database migration complete")

    yield

    down_migrations = list_down_migrations(migrations_dir)
    for migration in down_migrations:
        print(f"Performing down migration: {migration}")
        with open(migration) as f:
            migration_sql = f.read()
            with db_connection.cursor() as c:
                c.execute(typing.cast(typing.LiteralString, migration_sql))
                db_connection.commit()
    print("Database down migration complete")


def get_testcontainer_db_engine() -> Engine:
    host = postgres.get_container_host_ip()
    port = str(postgres.get_exposed_port(5432))
    username = postgres.username
    password = postgres.password
    database = postgres.dbname

    return create_engine(
        f"postgresql://{username}:{password}@{host}:{port}/{database}?sslmode=disable",
        echo=True,
    )


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


def list_down_migrations(dir_path: str) -> list[str]:
    migrations = []
    try:
        files = os.listdir(dir_path)
    except OSError as e:
        raise e

    for file in files:
        if not os.path.isdir(os.path.join(dir_path, file)) and file.endswith(
            ".down.sql"
        ):
            migrations.insert(0, os.path.join(dir_path, file))
            migrations.sort(reverse=True)

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
