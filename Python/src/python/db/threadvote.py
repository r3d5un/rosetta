import uuid

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from python.db.filters import Filter


class ThreadVote(SQLModel):
    __tablename__ = "thread_votes"  # type: ignore
    __table_args__ = {"schema": "forum"}

    thread_id: uuid.UUID = Field(nullable=False)
    user_id: uuid.UUID = Field(nullable=False)
    vote: int = Field(default=0, nullable=False)


class ThreadVoteModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select_count(self, filters: Filter) -> int | None:
        query = text(
            """
            SELECT CASE
                       WHEN SUM(vote) IS NULL THEN 0
                       ELSE SUM(vote)
                       END AS total_votes
            FROM forum.thread_votes
            WHERE (:thread_id IS NULL OR thread_id = :thread_id)
              AND (:user_id IS NULL OR user_id = :user_id);
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query, {"thread_id": filters.thread_id, "user_id": filters.user_id}
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return int(row.total_votes)
            except Exception as e:
                raise e

    def vote(self, vote: ThreadVote) -> ThreadVote | None:
        query = text(
            """
            WITH input_data AS (SELECT :thread_id AS thread_id,
                                       :user_id   AS user_id,
                                       :vote      AS vote),
                 delete_if_zero AS (
                     DELETE FROM forum.thread_votes
                         WHERE thread_id = (SELECT thread_id FROM input_data)
                             AND user_id = (SELECT user_id FROM input_data)
                             AND (SELECT vote FROM input_data) = 0)
            INSERT
            INTO forum.thread_votes (thread_id, user_id, vote)
            SELECT thread_id, user_id, vote
            FROM input_data
            WHERE vote != 0
            ON CONFLICT (thread_id, user_id) DO UPDATE
                SET vote = EXCLUDED.vote;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "thread_id": vote.thread_id,
                        "user_id": vote.user_id,
                        "vote": vote.vote,
                    },
                )
                if row is None:
                    raise NoResultFound
                session.commit()
                return vote
            except Exception as e:
                raise e
