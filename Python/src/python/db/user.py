import datetime
import uuid
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, Session, SQLModel, select


class UserBase(SQLModel):
    name: str = Field(nullable=False, max_length=256)
    username: str = Field(nullable=False, max_length=256)
    email: str = Field(nullable=False, max_length=256)


class User(UserBase, table=True):
    __tablename__ = "users"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: Optional[uuid.UUID] = Field(default_factory=uuid.uuid4, primary_key=True)
    created_at: Optional[datetime.datetime] = Field(
        nullable=False, default_factory=datetime.datetime.now
    )
    updated_at: Optional[datetime.datetime] = Field(
        nullable=False, default_factory=datetime.datetime.now
    )
    deleted: Optional[bool] = Field(default=False, nullable=False)
    deleted_at: Optional[datetime.datetime] = Field(nullable=True, default=None)


class UserCreate(UserBase):
    pass


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
                results = session.execute(query, {"id": id})
                session.commit()
                users = [
                    User(
                        id=row.id,
                        name=row.name,
                        username=row.username,
                        email=row.email,
                        created_at=row.created_at,
                        updated_at=row.updated_at,
                        deleted=row.deleted,
                        deleted_at=row.deleted_at,
                    )
                    for row in results
                ]

                if len(users) < 1:
                    raise NoResultFound

                return users[0]
            except Exception as e:
                raise e

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
                user.email = user_patch.email
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
