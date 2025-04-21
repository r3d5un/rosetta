import uuid

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from python.db.filters import Filter


class PostVote(SQLModel):
    __tablename__ = "post_votes"  # type: ignore
    __table_args__ = {"schema": "forum"}

    post_id: uuid.UUID = Field(nullable=False)
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

    def vote(self, vote: PostVote) -> PostVote | None:
        query = text(
            """
            WITH input_data AS (SELECT :post_id AS post_id,
                                       :user_id AS user_id,
                                       :vote    AS vote),
                 delete_if_zero AS (
                     DELETE FROM forum.post_votes
                         WHERE post_id = (SELECT post_id FROM input_data)
                             AND user_id = (SELECT user_id FROM input_data)
                             AND (SELECT vote FROM input_data) = 0)
            INSERT
            INTO forum.post_votes (post_id, user_id, vote)
            SELECT post_id, user_id, vote
            FROM input_data
            WHERE vote != 0
            ON CONFLICT (post_id, user_id) DO UPDATE
                SET vote = EXCLUDED.vote;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "post_id": vote.post_id,
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
