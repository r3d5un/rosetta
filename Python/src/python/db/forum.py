import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine
from sqlmodel import Field, SQLModel


class Forum(SQLModel, table=True):
    __tablename__ = "forums"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: uuid.UUID = Field(primary_key=True)
    owner_id: uuid.UUID = Field()
    name: str = Field(nullable=False, max_length=256)
    description: str = Field(nullable=False)
    created_at: datetime = Field(nullable=False, default_factory=datetime.now)
    updated_at: datetime = Field(nullable=False, default_factory=datetime.now)
    deleted: bool = Field(default=False, nullable=False)
    deleted_at: Optional[datetime] = Field(nullable=False, default_factory=datetime.now)


class ForumModel:
    def __init__(self, engine: Engine):
        self.engine = engine
