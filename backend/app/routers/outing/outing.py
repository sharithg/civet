from fastapi import APIRouter

from app.routers.deps import SessionDep
from app.routers.outing.dto import CreateOuting
from app.routers.outing.model import (
    create_new_outing,
    get_outings,
    get_receipts_for_outing,
)

router = APIRouter(prefix="/outing", tags=["items"])


@router.post("")
async def read_items(request: CreateOuting, session: SessionDep):
    outing_id = await create_new_outing(session, request.name)
    return {"id": outing_id}


@router.get("")
async def get_outing_list(session: SessionDep):
    outings = await get_outings(session)
    return outings


@router.get("/{outing_id}/receipts")
async def receipts(outing_id: str, session: SessionDep):
    receipts = await get_receipts_for_outing(session, outing_id)
    return receipts
