import os
from dateutil import parser
import hashlib
from app.cloud_vision.cloud_vision import CloudVision
from app.openai.openai import OpenAIClient
from app.receipt.utils import group_text_by_lines
from app.receipt.schema import Receipt, receipt_schema
from app.storage.storage import upload_image_bytes
import pytz

PROMPT = "Convert the given text of a receipt into a structured output format"

vision = CloudVision()
openai = OpenAIClient()


def hash_image_bytes(image_bytes: bytes) -> str:
    return hashlib.sha256(image_bytes).hexdigest()


class Extract:
    def __init__(self, image_bytes: bytes, fname: str):
        self.image_bytes = image_bytes
        self.fname = str
        self.image_hash = hash_image_bytes(image_bytes)
        self.file_extension = os.path.splitext(fname)[1][1:]

    def upload(self):
        obj_name = f"{self.image_hash}.{self.file_extension}"
        bucket = "receipts"
        upload_image_bytes(
            "receipts", obj_name, self.image_bytes, f"image/{self.file_extension}"
        )
        return bucket, obj_name

    def extract_text(self):
        response = vision.detect_text(self.image_bytes, self.image_hash)
        lines = group_text_by_lines(response.text_annotations)
        return "\n".join(lines)

    def structured_outut(self, input):
        return openai.text_to_json(PROMPT, input, "receipt_info", receipt_schema)

    def run(self):
        bucket, key = self.upload()
        text = self.extract_text()
        output = self.structured_outut(text)
        return self.to_model(output), text, bucket, key

    def to_model(self, output):
        dt = None
        date_str = output["opened"]
        try:
            dt = parser.parse(date_str)
            if dt.tzinfo is not None:
                dt = dt.astimezone(tz=None).replace(tzinfo=None)
        except Exception as e:
            print(f"[WARN] Unable to parse date string: {date_str}", e)
        output["opened"] = dt
        return Receipt(**output)
