import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from src.python.db.filters import Filter, Metadata


class Post(SQLModel, table=True):
    __tablename__ = "posts"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    thread_id: uuid.UUID = Field()
    reply_to: Optional[uuid.UUID] = Field(default=None)
    author_id: uuid.UUID = Field()
    content: str = Field(nullable=False)
    created_at: datetime = Field(nullable=False, default_factory=datetime.now)
    updated_at: datetime = Field(nullable=False, default_factory=datetime.now)
    deleted: bool = Field(default=False, nullable=False)
    deleted_at: Optional[datetime] = Field(nullable=False, default_factory=datetime.now)


class PostPatch(SQLModel):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    thread_id: uuid.UUID = Field(nullable=False)
    content: Optional[str] = Field(default=None, nullable=False)


class PostModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select(self, id: uuid.UUID) -> Post | None:
        query = text(
            """
            SELECT id,
                   thread_id,
                   reply_to,
                   author_id,
                   content,
                   created_at,
                   updated_at,
                   likes,
                   deleted,
                   deleted_at
            FROM forum.posts
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
                return Post(
                    id=row.id,
                    thread_id=row.thread_id,
                    reply_to=row.reply_to,
                    author_id=row.author_id,
                    content=row.content,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def select_all(
        self, filters: Filter
    ) -> tuple[list[Post] | None, Metadata | None] | None:
        query = text(
            f"""
            SELECT id,
                   thread_id,
                   reply_to,
                   author_id,
                   content,
                   created_at,
                   updated_at,
                   likes,
                   deleted,
                   deleted_at
            FROM forum.posts
            WHERE (:id IS NULL OR id = :id)
              AND (:thread_id IS NULL OR thread_id = :thread_id)
              AND (:author_id IS NULL OR author_id = :author_id)
              AND (:created_at_from IS NULL or created_at >= :created_at_from)
              AND (:created_at_to IS NULL or created_at <= :created_at_to)
              AND (:updated_at_from IS NULL or updated_at >= :updated_at_from)
              AND (:updated_at_to IS NULL or updated_at <= :updated_at_to)
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
                        "thread_id": filters.thread_id,
                        "author_id": filters.author_id,
                        "created_at_from": filters.created_at_from,
                        "created_at_to": filters.created_at_to,
                        "updated_at_from": filters.updated_at_from,
                        "updated_at_to": filters.updated_at_to,
                        "deleted_at_from": filters.deleted_at_from,
                        "deleted_at_to": filters.deleted_at_to,
                        "deleted": filters.deleted,
                        "last_seen": filters.last_seen,
                    },
                ).fetchall()
                forums = [
                    Post(
                        id=row.id,
                        thread_id=row.thread_id,
                        reply_to=row.reply_to,
                        author_id=row.author_id,
                        content=row.content,
                        created_at=row.created_at,
                        updated_at=row.updated_at,
                        deleted=row.deleted,
                        deleted_at=row.deleted_at,
                    )
                    for row in rows
                ]
                length = len(forums)
                metadata = Metadata()
                if length > 0:
                    id = forums[length - 1].id
                    if id is not None:
                        metadata.last_seen = id
                    metadata.next = True
                metadata.response_length = length

                return (forums, metadata)

            except Exception as e:
                raise e

    def insert(self, post: Post) -> Post | None:
        query = text(
            """
            INSERT INTO forum.posts(thread_id, reply_to, content, author_id)
            VALUES (:thread_id,
                    :reply_to,
                    :content,
                    :author_id)
            RETURNING id,
                thread_id,
                reply_to,
                author_id,
                content,
                created_at,
                updated_at,
                likes,
                deleted,
                deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "thread_id": post.thread_id,
                        "reply_to": post.reply_to,
                        "content": post.content,
                        "author_id": post.author_id,
                    },
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Post(
                    id=row.id,
                    thread_id=row.thread_id,
                    reply_to=row.reply_to,
                    author_id=row.author_id,
                    content=row.content,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def update(self, patch: PostPatch) -> Post | None:
        query = text(
            """
            UPDATE forum.posts
            SET content = COALESCE(:content, content)
            WHERE id = :id
              AND thread_id = :thread_id
            RETURNING id,
                thread_id,
                reply_to,
                author_id,
                content,
                created_at,
                updated_at,
                likes,
                deleted,
                deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "id": patch.id,
                        "thread_id": patch.thread_id,
                        "content": patch.content,
                    },
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Post(
                    id=row.id,
                    thread_id=row.thread_id,
                    reply_to=row.reply_to,
                    author_id=row.author_id,
                    content=row.content,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def soft_delete(self, id: uuid.UUID) -> Post | None:
        query = text(
            """
            UPDATE forum.posts
            SET deleted    = TRUE,
                deleted_at = NOW(),
                updated_at = NOW()
            WHERE id = :id
            RETURNING id,
                thread_id,
                reply_to,
                author_id,
                content,
                created_at,
                updated_at,
                likes,
                deleted,
                deleted_at;
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
                return Post(
                    id=row.id,
                    thread_id=row.thread_id,
                    reply_to=row.reply_to,
                    author_id=row.author_id,
                    content=row.content,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e
