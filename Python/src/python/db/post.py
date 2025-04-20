import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel


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
    thread_id: Optional[uuid.UUID] = Field(default=None, nullable=False)
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
