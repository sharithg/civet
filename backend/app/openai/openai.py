import os
import json
import hashlib
from openai import OpenAI  # Adjust this import if needed


class OpenAIClient:
    def __init__(self):
        self.client = OpenAI(api_key=os.environ["OPENAI_API_KEY"])
        self.cache_dir = "cache/openai"
        os.makedirs(self.cache_dir, exist_ok=True)

    def _compute_cache_key(self, prompt, input, schema_name, schema):
        key_data = json.dumps(
            {
                "prompt": prompt,
                "input": input,
                "schema_name": schema_name,
                "schema": schema,
            },
            sort_keys=True,
        ).encode("utf-8")
        return hashlib.sha256(key_data).hexdigest()

    def _get_cache_path(self, cache_key):
        return os.path.join(self.cache_dir, f"{cache_key}.json")

    def text_to_json(self, prompt, input, schema_name, schema):
        cache_key = self._compute_cache_key(prompt, input, schema_name, schema)
        cache_path = self._get_cache_path(cache_key)

        if os.path.exists(cache_path):
            with open(cache_path, "r") as f:
                return json.load(f)

        response = self.client.responses.create(
            model="gpt-4o-mini-2024-07-18",
            input=[
                {"role": "system", "content": prompt},
                {"role": "user", "content": input},
            ],
            text={
                "format": {
                    "type": "json_schema",
                    "name": schema_name,
                    "schema": schema,
                    "strict": True,
                }
            },
        )

        event = json.loads(response.output_text)

        with open(cache_path, "w") as f:
            json.dump(event, f, indent=2)

        return event
