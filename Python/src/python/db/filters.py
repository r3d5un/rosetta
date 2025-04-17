import uuid
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional


@dataclass
class Filter:
    id: Optional[uuid.UUID] = None
    owner_id: Optional[uuid.UUID] = None
    user_id: Optional[uuid.UUID] = None
    post_id: Optional[uuid.UUID] = None
    thread_id: Optional[uuid.UUID] = None
    forum_id: Optional[uuid.UUID] = None
    author_id: Optional[uuid.UUID] = None
    name: Optional[str] = None
    title: Optional[str] = None
    username: Optional[str] = None
    email: Optional[str] = None
    created_at_from: Optional[datetime] = None
    created_at_to: Optional[datetime] = None
    updated_at_from: Optional[datetime] = None
    updated_at_to: Optional[datetime] = None
    deleted_at_from: Optional[datetime] = None
    deleted_at_to: Optional[datetime] = None
    deleted: Optional[bool] = None
    is_locked: Optional[bool] = None

    order_by: list[str] = field(default_factory=list)
    order_by_safelist: list[str] = field(default_factory=list)
    last_seen: uuid.UUID = uuid.UUID("00000000-0000-0000-0000-000000000000")
    page_size: int = 1

    def create_order_by_clause(self) -> str:
        if len(self.order_by) < 1:
            return "ORDER BY id"

        clauses: list[str] = []
        for clause in self.order_by:
            if clause.startswith("-"):
                clauses.append(clause.strip("-") + " DESC")
            else:
                clauses.append(clause + " ASC")

        clauses.append("id ASC")

        return "ORDER BY ".join(clauses)


@dataclass
class Metadata:
    last_seen: uuid.UUID = uuid.UUID("00000000-0000-0000-0000-000000000000")
    next: bool = False
    response_length: int = 0
