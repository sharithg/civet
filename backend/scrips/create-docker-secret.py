#!/usr/bin/env python3

import os
import subprocess
from pathlib import Path
from dotenv import load_dotenv

# Define the mapping of environment variables to secret names
SECRET_MAPPING = {
    "DATABASE_URL": "db_url",
    "MINIIO_SECRET_ACCESS_KEY": "minio_secret_access_key",
    "JWT_SECRET": "jwt_secret",
    "OPENAI_API_KEY": "openai_api_key",
    "GOOGLE_CLIENT_SECRET": "google_client_secret",
}


def create_docker_secret(name: str, value: str) -> None:
    """Create a Docker secret with the given name and value."""
    try:
        process = subprocess.Popen(
            ["docker", "secret", "create", name, "-"], stdin=subprocess.PIPE, text=True
        )
        process.communicate(input=value)
        print(f"Created secret: {name}")
    except subprocess.CalledProcessError as e:
        print(f"Error creating secret {name}: {e}")


def main():
    # Get the path to the .env file (one directory up from the script)
    env_path = Path(__file__).parent.parent / ".env"

    if not env_path.exists():
        print("Error: .env file not found")
        return 1

    # Load environment variables from .env file
    load_dotenv(env_path)

    # Create Docker secrets for mapped environment variables
    for env_var, secret_name in SECRET_MAPPING.items():
        value = os.getenv(env_var)
        if value:
            create_docker_secret(secret_name, value)

    print("Docker secrets creation completed")
    return 0


if __name__ == "__main__":
    exit(main())
