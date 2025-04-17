import datetime
import uuid
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, Session, SQLModel, select


class User(SQLModel, table=True):
    __tablename__ = "users"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    name: str = Field(nullable=False, max_length=256)
    username: str = Field(nullable=False, max_length=256)
    email: str = Field(nullable=False, max_length=256)
    created_at: Optional[datetime.datetime] = Field(
        nullable=False, default_factory=datetime.datetime.now
    )
    updated_at: Optional[datetime.datetime] = Field(
        nullable=False, default_factory=datetime.datetime.now
    )
    deleted: Optional[bool] = Field(default=False, nullable=False)
    deleted_at: Optional[datetime.datetime] = Field(nullable=True, default=None)


class UserPatch(SQLModel):
    __tablename__ = "users"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    name: Optional[str] = Field(default=None, nullable=False, max_length=256)
    username: Optional[str] = Field(default=None, nullable=False, max_length=256)
    email: Optional[str] = Field(default=None, nullable=False, max_length=256)


class UserModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select(self, id: uuid.UUID) -> User | None:
        query = text(
            """
            SELECT id, name, username, email, created_at, updated_at, deleted, deleted_at
            FROM forum.users
            WHERE id = :id;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(query, {"id": id}).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return User(
                    id=row.id,
                    name=row.name,
                    username=row.username,
                    email=row.email,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def insert(self, user: User) -> User:
        query = text(
            """
            INSERT INTO forum.users(name, username, email)
            VALUES (:name, :username, :email)
            RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {"name": user.name, "username": user.username, "email": user.email},
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return User(
                    id=row.id,
                    name=row.name,
                    username=row.username,
                    email=row.email,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def update(self, user_patch: UserPatch) -> User | None:
        query = text(
            """
            UPDATE forum.users
            SET name       = COALESCE(:name, name),
                username   = COALESCE(:username, username),
                email      = COALESCE(:email, email),
                updated_at = NOW()
            WHERE id = :id
            RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "id": user_patch.id,
                        "name": user_patch.name,
                        "username": user_patch.username,
                        "email": user_patch.email,
                    },
                ).first()
                session.commit()
                if row is None:
                    raise NoResultFound
                return User(
                    id=row.id,
                    name=row.name,
                    username=row.username,
                    email=row.email,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                )
            except Exception as e:
                raise e

    def soft_delete(self, id: uuid.UUID) -> User | None:
        with Session(self.engine) as session:
            user = session.exec(select(User).where(User.id == id)).first()
            if user is None:
                return None
            user.deleted = True
            session.commit()
            session.refresh(user)
            return user

    def restore(self, id: uuid.UUID) -> User | None:
        with Session(self.engine) as session:
            user = session.exec(select(User).where(User.id == id)).first()
            if user is None:
                return None
            user.deleted = False
            session.commit()
            session.refresh(user)
            return user

    def delete(self, id: uuid.UUID) -> User | None:
        with Session(self.engine) as session:
            user = session.exec(select(User).where(User.id == id)).first()
            if user is None:
                return None
            session.delete(user)
            session.commit()

            return user
