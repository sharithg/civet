from fastapi import APIRouter, Request, UploadFile
from fastapi.responses import JSONResponse

from app.cloud_vision.cloud_vision import CloudVision
from app.routers.receipt.model import (
    get_receipt_by_hash,
    get_receipt_by_id,
    get_receipts,
    save_receipt,
)
from app.utils.utils import row_to_dict
from app.receipt.receipt import Extract
from app.routers.deps import SessionDep
from app.storage.storage import get_presigned_url
from app.template import templates

router = APIRouter(prefix="/receipt", tags=["items"])

vision = CloudVision()


@router.post("/upload")
async def read_items(file: UploadFile, session: SessionDep, request: Request):
    if file.content_type.startswith("image"):
        outing_id = request.headers.get("outingId")

        if not outing_id:
            return JSONResponse(
                content={"message": "outing id required"},
                status_code=404,
            )

        image_bytes = await file.read()
        file_name = file.filename

        extractor = Extract(image_bytes, file_name)
        image_hash = extractor.image_hash

        existing = await get_receipt_by_hash(session, image_hash)

        if existing:
            return JSONResponse(
                content={"receipt": row_to_dict(existing), "existing": True},
                status_code=200,
            )

        receipt, text, bucket, key = extractor.run()

        await save_receipt(
            session=session,
            bucket=bucket,
            hash=image_hash,
            key=key,
            name=file_name,
            receipt=receipt,
            text_val=text,
            outing_id=outing_id,
        )

        existing = await get_receipt_by_hash(session, image_hash)

        return JSONResponse(
            content={"receipt": row_to_dict(existing), "existing": False},
            status_code=200,
        )
    else:
        return {"error": "Uploaded file is not an image"}


@router.get("/view/{id}")
async def read_item(request: Request, session: SessionDep, id: str):
    existing = await get_receipt_by_id(session, id)

    if not existing:
        return JSONResponse(
            content={"message": "Not found"},
            status_code=404,
        )

    presigned_url = get_presigned_url(existing["bucket"], existing["key"])

    print(existing)

    return templates.TemplateResponse(
        request=request,
        name="receipt.html",
        context={"id": id, "presigned_url": presigned_url, "data": existing},
    )


@router.get("/all")
async def receipt_list(request: Request, session: SessionDep):
    receipts = await get_receipts(session)
    print(receipts)
    return templates.TemplateResponse(
        request=request,
        name="receipt_list.html",
        context={"id": id, "data": receipts},
    )
