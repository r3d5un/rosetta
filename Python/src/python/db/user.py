import datetime
import uuid
from typing import Optional

from sqlalchemy import Engine
from sqlmodel import Field, Session, SQLModel, select


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

    def select(self, id: uuid.UUID) -> User | None:
        return Session(self.engine).exec(select(User).where(User.id == id)).first()

    def insert(self, user: User) -> User:
        with Session(self.engine) as session:
            session.add(user)
            session.commit()
            session.refresh(user)
            session.close()
            return user

    def update(self, user_patch: User) -> User | None:
        with Session(self.engine) as session:
            user = session.exec(select(User).where(User.id == user_patch.id)).first()
            if user is None:
                return None
            if user_patch.name != "" and user_patch.name is not None:
                user.name = user_patch.name
            if user_patch.username != "" and user_patch.username is not None:
                user.username = user_patch.username
            if user_patch.email != "" and user_patch.email is not None:
                user.username = user_patch.username
            user_patch.updated_at = datetime.datetime.now()

            session.commit()
            session.refresh(user)

            return user

    def soft_delete(self, id: uuid.UUID) -> User | None:
        with Session(self.engine) as session:
            user = session.exec(select(User).where(User.id == id)).first()
            if user is None:
                return None
            user.deleted = True
            session.commit()
            session.refresh(user)
            return user
