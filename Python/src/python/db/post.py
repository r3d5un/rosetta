import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine
from sqlmodel import Field, SQLModel


class Post(SQLModel, table=True):
    __tablename__ = "forums"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    thread_id: uuid.UUID = Field()
    reply_to: Optional[uuid.UUID] = Field(default=None)
    authro_id: uuid.UUID = Field()
    content: str = Field(nullable=False)
    created_at: datetime = Field(nullable=False, default_factory=datetime.now)
    updated_at: datetime = Field(nullable=False, default_factory=datetime.now)
    deleted: bool = Field(default=False, nullable=False)
    deleted_at: Optional[datetime] = Field(nullable=False, default_factory=datetime.now)


class PostPatch(SQLModel, table=True):
    __tablename__ = "forums"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    thread_id: Optional[uuid.UUID] = Field(default=None, nullable=False)
    content: Optional[str] = Field(default=None, nullable=False)


class PostModel:
    def __init__(self, engine: Engine):
        self.engine = engine
