from sqlalchemy import Engine
from sqlmodel import create_engine

from python.db.forum import ForumModel
from python.db.thread import ThreadModel
from python.db.user import UserModel

connstr = "postgresql://postgres:postgres@localhost?database=rosetta"
engine = create_engine(connstr, echo=True)


class Models:
    def __init__(self, engine: Engine):
        self.engine = engine
        self.users = UserModel(engine)
        self.forums = ForumModel(engine)
        self.threads = ThreadModel(engine)
