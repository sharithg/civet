from typing import Annotated
from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer
from app.core.db import engine
from typing import AsyncGenerator
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker


AsyncSessionLocal = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with AsyncSessionLocal() as session:
        yield session


SessionDep = Annotated[AsyncSession, Depends(get_db)]
OAuth2Dep = Depends(oauth2_scheme)
