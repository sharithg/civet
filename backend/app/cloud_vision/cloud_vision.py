import os
import pickle
from google.cloud import vision


class CloudVision:
    def __init__(self, cache_dir="cache/cloud_vision"):
        self.client = vision.ImageAnnotatorClient()
        self.cache_dir = cache_dir
        os.makedirs(self.cache_dir, exist_ok=True)

    def detect_text(self, content: bytes, image_hash: str):
        cache_path = os.path.join(self.cache_dir, f"{image_hash}.pkl")

        if os.path.exists(cache_path):
            with open(cache_path, "rb") as f:
                print(f"Loaded cached response for hash {image_hash}")
                return pickle.load(f)

        print(f"Calling API for hash {image_hash}")
        image = vision.Image(content=content)
        response = self.client.text_detection(image=image)
        with open(cache_path, "wb") as f:
            pickle.dump(response, f)

        return response
