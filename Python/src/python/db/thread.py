import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from src.python.db.filters import Filter, Metadata


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
    forum_id: Optional[uuid.UUID] = Field(default=None, nullable=True)
    title: Optional[str] = Field(default=None, nullable=True, max_length=256)
    author_id: Optional[uuid.UUID] = Field(default=None, nullable=True)


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

    def select_all(
        self, filters: Filter
    ) -> tuple[list[Thread] | None, Metadata | None] | None:
        query = text(
            f"""
            SELECT id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at, likes
            FROM forum.threads
            WHERE (:id IS NULL OR id = :id)
              AND (:forum_id IS NULL OR forum_id = :forum_id)
              AND (:title IS NULL OR title = :title)
              AND (:author_id IS NULL OR author_id = :author_id)
              AND (:created_at_from IS NULL or created_at >= :created_at_from)
              AND (:created_at_to IS NULL or created_at <= :created_at_to)
              AND (:updated_at_from IS NULL or updated_at >= :updated_at_from)
              AND (:updated_at_to IS NULL or updated_at <= :updated_at_to)
              AND (:is_locked IS NULL or is_locked = :is_locked)
              AND (:deleted IS NULL or deleted = :deleted)
              AND (:deleted_at_from IS NULL or deleted_at >= :deleted_at_from)
              AND (:deleted_at_to IS NULL or deleted_at <= :deleted_at_to)
              AND id > :last_seen
            {filters.create_order_by_clause()}
            LIMIT :page_size;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                rows = session.execute(
                    query,
                    {
                        "page_size": filters.page_size,
                        "id": filters.id,
                        "author_id": filters.author_id,
                        "title": filters.title,
                        "forum_id": filters.forum_id,
                        "name": filters.name,
                        "username": filters.username,
                        "email": filters.email,
                        "created_at_from": filters.created_at_from,
                        "created_at_to": filters.created_at_to,
                        "updated_at_from": filters.updated_at_from,
                        "updated_at_to": filters.updated_at_to,
                        "deleted_at_from": filters.deleted_at_from,
                        "deleted_at_to": filters.deleted_at_to,
                        "deleted": filters.deleted,
                        "last_seen": filters.last_seen,
                        "is_locked": filters.is_locked,
                    },
                ).fetchall()
                threads = [
                    Thread(
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
                    for row in rows
                ]
                length = len(threads)
                metadata = Metadata()
                if length > 0:
                    id = threads[length - 1].id
                    if id is not None:
                        metadata.last_seen = id
                    metadata.next = True
                metadata.response_length = length

                return (threads, metadata)

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

    def update(self, thread_patch: ThreadPatch) -> Thread | None:
        query = text(
            """
            UPDATE forum.threads
            SET forum_id = COALESCE(:forum_id, forum_id),
                title = COALESCE(:title, title),
                author_id = COALESCE(:author_id, author_id)
            WHERE id = :id
            RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at, likes;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "id": thread_patch.id,
                        "forum_id": thread_patch.forum_id,
                        "title": thread_patch.title,
                        "author_id": thread_patch.author_id,
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
