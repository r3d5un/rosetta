import datetime
import uuid
from typing import Optional

from sqlalchemy import Engine
from sqlmodel import Field, Session, SQLModel


class User(SQLModel, table=True):
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
    deleted: Optional[bool] = Field(default=True, nullable=False)
    deleted_at: Optional[datetime.datetime] = Field(nullable=True, default=None)


class UserModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def insert(self, user: User) -> User:
        session = Session(self.engine)
        session.add(user)
        session.refresh(user)
        session.close()
        return user
