import uuid

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from python.db.filters import Filter


class PostVote(SQLModel):
    __tablename__ = "post_votes"  # type: ignore
    __table_args__ = {"schema": "forum"}

    thread_id: uuid.UUID = Field(nullable=False)
    user_id: uuid.UUID = Field(nullable=False)
    vote: int = Field(default=0, nullable=False)


class PostVoteModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select_count(self, filters: Filter) -> int | None:
        query = text(
            """
            SELECT CASE
                       WHEN SUM(vote) IS NULL THEN 0
                       ELSE SUM(vote)
                       END AS total_votes
            FROM forum.post_votes
            WHERE (:post_id IS NULL OR post_id = :post_id)
              AND (:user_id IS NULL OR user_id = :user_id);
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query, {"post_id": filters.post_id, "user_id": filters.user_id}
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return int(row.total_votes)
            except Exception as e:
                raise e
