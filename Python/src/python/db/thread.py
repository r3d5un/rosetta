import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel


class Thread(SQLModel, table=True):
    __tablename__ = "threads"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    forum_id: uuid.UUID = Field(nullable=False)
    title: str = Field(nullable=False, max_length=256)
    author_id: uuid.UUID = Field(nullable=False)
    created_at: datetime = Field(nullable=False, default_factory=datetime.now)
    updated_at: datetime = Field(nullable=False, default_factory=datetime.now)
    is_locked: bool = Field(default=False, nullable=False)
    deleted: bool = Field(default=False, nullable=False)
    deleted_at: Optional[datetime] = Field(nullable=False, default_factory=datetime.now)
    likes: Optional[int] = Field(default=0, nullable=False)


class ThreadPatch(SQLModel):
    __tablename__ = "threads"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    forum_id: Optional[uuid.UUID] = Field(nullable=False)
    title: Optional[str] = Field(nullable=False, max_length=256)
    author_id: Optional[uuid.UUID] = Field()


class ThreadModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select(self, id: uuid.UUID) -> Thread | None:
        query = text(
            """
            SELECT id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at, likes
            FROM forum.threads
            WHERE id = :id;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "id": id,
                    },
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Thread(
                    id=row.id,
                    forum_id=row.forum_id,
                    title=row.title,
                    author_id=row.author_id,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    is_locked=row.is_locked,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                    likes=row.likes,
                )
            except Exception as e:
                raise e

    def insert(self, thread: Thread):
        query = text(
            """
            INSERT INTO forum.threads(forum_id, title, author_id)
            VALUES(:forum_id, :title, :author_id)
            RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at, likes;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "forum_id": thread.forum_id,
                        "title": thread.title,
                        "author_id": thread.author_id,
                    },
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Thread(
                    id=row.id,
                    forum_id=row.forum_id,
                    title=row.title,
                    author_id=row.author_id,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    is_locked=row.is_locked,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                    likes=row.likes,
                )
            except Exception as e:
                raise e
