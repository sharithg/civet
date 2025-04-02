import os
from sqlalchemy.ext.asyncio import create_async_engine

DB_URL = os.environ["DATABASE_URL"]
engine = create_async_engine(DB_URL)
