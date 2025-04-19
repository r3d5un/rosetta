import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine
from sqlmodel import Field, SQLModel


class Thread(SQLModel, table=True):
    __tablename__ = "threads"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    forum_id: uuid.UUID = Field(nullable=False)
    title: str = Field(nullable=False, max_length=256)
    forum_id: uuid.UUID = Field(nullable=False)
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
