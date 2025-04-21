import uuid

from sqlalchemy import Engine
from sqlmodel import Field, SQLModel


class PostVote(SQLModel, table=True):
    __tablename__ = "post_votes"  # type: ignore
    __table_args__ = {"schema": "forum"}

    thread_id: uuid.UUID = Field(nullable=False)
    user_id: uuid.UUID = Field(nullable=False)
    vote: int = Field(default=0, nullable=False)


class PostVoteModel:
    def __init__(self, engine: Engine):
        self.engine = engine
