from pydantic import BaseModel


class CreateOuting(BaseModel):
    name: str
