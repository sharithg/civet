from datetime import timedelta
from io import BytesIO
import os
from minio import Minio

client = Minio(
    os.environ["MINIIO_HOST"],
    access_key=os.environ["MINIIO_ACCESS_KEY_ID"],
    secret_key=os.environ["MINIIO_SECRET_ACCESS_KEY"],
    secure=False,
)


def upload_image_bytes(
    bucket_name: str,
    object_name: str,
    image_bytes: bytes,
    content_type: str = "image/png",
):
    if not client.bucket_exists(bucket_name):
        client.make_bucket(bucket_name)

    byte_stream = BytesIO(image_bytes)
    byte_stream.seek(0)
    size = len(image_bytes)

    client.put_object(
        bucket_name,
        object_name,
        data=byte_stream,
        length=size,
        content_type=content_type,
    )


def get_presigned_url(
    bucket: str, object_name: str, expires_in_seconds: int = 3600
) -> str:
    """Generate a pre-signed URL to access an object."""
    url = client.presigned_get_object(
        bucket_name=bucket,
        object_name=object_name,
        expires=timedelta(seconds=expires_in_seconds),
    )
    return url
