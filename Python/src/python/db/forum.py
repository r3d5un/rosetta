import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Engine, text
from sqlalchemy.exc import NoResultFound
from sqlalchemy.orm import sessionmaker
from sqlmodel import Field, SQLModel

from src.python.db.filters import Filter, Metadata


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


class ForumPatch(SQLModel):
    __tablename__ = "users"  # type: ignore
    __table_args__ = {"schema": "forum"}

    id: uuid.UUID = Field(primary_key=True)
    owner_id: uuid.UUID = Field()
    name: str = Field(nullable=False, max_length=256)
    description: str = Field(nullable=False)


class ForumModel:
    def __init__(self, engine: Engine):
        self.engine = engine

    def select(self, id: uuid.UUID) -> Forum | None:
        query = text(
            """
            SELECT id, owner_id, name, description, created_at, updated_at, deleted, deleted_at
            FROM forum.forums
            WHERE id = :id;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(query, {"id": id}).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def select_all(
        self, filters: Filter
    ) -> tuple[list[Forum] | None, Metadata | None] | None:
        query = text(
            f"""
            SELECT id, owner_id, name, description, created_at, updated_at, deleted, deleted_at
            FROM forum.forums
            WHERE (:id IS NULL OR id = :id)
              AND (:owner_id IS NULL OR owner_id = :owner_id)
              AND (:name IS NULL or name = :name)
              AND (:created_at_from IS NULL or created_at >= :created_at_from)
              AND (:created_at_to IS NULL or created_at <= :created_at_to)
              AND (:updated_at_from IS NULL or updated_at >= :updated_at_from)
              AND (:updated_at_to IS NULL or updated_at <= :updated_at_to)
              AND (:deleted IS NULL or deleted = :deleted)
              AND (:deleted_at_from IS NULL or deleted_at >= :deleted_at_from)
              AND (:delted_at_to IS NULL or deleted_at <= :delted_at_to)
              AND id > :last_seen
            {filters.create_order_by_clause()}
            LIMIT :page_size
            """
        )
        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                rows = session.execute(
                    query,
                    {
                        "page_size": filters.page_size,
                        "id": filters.id,
                        "name": filters.name,
                        "username": filters.username,
                        "email": filters.email,
                        "created_at_from": filters.created_at_from,
                        "created_at_to": filters.created_at_to,
                        "updated_at_from": filters.updated_at_from,
                        "updated_at_to": filters.updated_at_to,
                        "deleted_at_from": filters.deleted_at_from,
                        "deleted_at_to": filters.deleted_at_from,
                        "last_seen": filters.last_seen,
                    },
                ).fetchall()
                users = [
                    Forum(
                        id=row.id,
                        owner_id=row.owner_id,
                        name=row.name,
                        description=row.description,
                        created_at=row.created_at,
                        updated_at=row.updated_at,
                        deleted=row.deleted,
                        deleted_at=row.deleted_at,
                    )
                    for row in rows
                ]
                length = len(users)
                metadata = Metadata()
                if length > 0:
                    id = users[length - 1].id
                    if id is not None:
                        metadata.last_seen = id
                    metadata.next = True
                metadata.response_length = length

                return (users, metadata)

            except Exception as e:
                raise e

    def insert(self, forum: Forum) -> Forum | None:
        query = text(
            """
            INSERT INTO forum.forums(owner_id, name, description)
            VALUES(:owner_id, :name, :description)
            RETURNING id, owner_id, name, description, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "owner_id": forum.owner_id,
                        "name": forum.name,
                        "description": forum.description,
                    },
                ).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def update(self, forum_patch: ForumPatch) -> Forum | None:
        query = text(
            """
            UPDATE forum.forums
            SET name = COALESCE(:name, name),
                owner_id = COALESCE(:owner_id, owner_id),
                description = COALESCE(:description, description),
                updated_at = NOW()
            WHERE id = :id
            RETURNING id, owner_id, name, description, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(
                    query,
                    {
                        "id": forum_patch.id,
                        "owner_id": forum_patch.name,
                        "description": forum_patch.description,
                        "name": forum_patch.name,
                    },
                ).first()
                session.commit()
                if row is None:
                    raise NoResultFound
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def soft_delete(self, id: uuid.UUID) -> Forum | None:
        query = text(
            """
            UPDATE forum.forums
            SET deleted    = TRUE,
                deleted_at = NOW(),
                updated_at = NOW()
            WHERE id = :id
            RETURNING id, owner_id, name, description, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(query, {"id": id}).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def restore(self, id: uuid.UUID) -> Forum | None:
        query = text(
            """
            UPDATE forum.forums
            SET deleted    = FALSE,
                deleted_at = NULL,
                updated_at = NOW()
            WHERE id = :id
            RETURNING id, owner_id, name, description, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(query, {"id": id}).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e

    def delete(self, id: uuid.UUID) -> Forum | None:
        query = text(
            """
            DELETE
            FROM forum.forums
            WHERE id = :id
            RETURNING id, owner_id, name, description, created_at, updated_at, deleted, deleted_at;
            """
        )

        session = sessionmaker(bind=self.engine)()
        with session:
            try:
                row = session.execute(query, {"id": id}).first()
                if row is None:
                    raise NoResultFound
                session.commit()
                return Forum(
                    id=row.id,
                    owner_id=row.owner_id,
                    name=row.name,
                    description=row.description,
                    created_at=row.created_at,
                    updated_at=row.updated_at,
                    deleted=row.deleted,
                    deleted_at=row.deleted_at,
                )
            except Exception as e:
                raise e
