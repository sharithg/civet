import base64
from datetime import datetime
import hashlib
from typing import Dict
import uuid
from fastapi.responses import JSONResponse
from jose import jwt


def parse_cookie_header(cookie_header: str) -> Dict[str, Dict[str, str]]:
    cookies: Dict[str, Dict[str, str]] = {}

    for cookie in cookie_header.split(";"):
        trimmed = cookie.strip()

        if "=" in trimmed:
            key, value = trimmed.split("=", 1)
            name = key.strip()
            if name not in cookies:
                cookies[name] = {"value": value}
            else:
                cookies[name]["value"] = value
        elif trimmed.lower() == "httponly":
            last_cookie = list(cookies.keys())[-1]
            cookies[last_cookie]["httpOnly"] = "true"
        elif trimmed.lower().startswith("expires="):
            last_cookie = list(cookies.keys())[-1]
            cookies[last_cookie]["expires"] = trimmed[8:]
        elif trimmed.lower().startswith("max-age="):
            last_cookie = list(cookies.keys())[-1]
            cookies[last_cookie]["maxAge"] = trimmed[8:]

    return cookies


def build_cookie_string(name: str, value: str, options: dict) -> str:
    parts = [
        f"{name}={value}",
        f"Max-Age={options['max_age']}",
        f"Path={options['path']}",
        f"SameSite={options['samesite']}",
    ]
    if options.get("httponly"):
        parts.append("HttpOnly")
    if options.get("secure"):
        parts.append("Secure")
    return "; ".join(parts)


def compute_at_hash(access_token: str, alg: str = "RS256") -> str:
    # Determine hash algorithm from alg
    if alg not in ("RS256", "HS256"):
        raise ValueError(f"Unsupported alg: {alg}")

    hash_fn = hashlib.sha256  # RS256/HS256 → SHA-256

    # Step 1: Hash the ASCII representation of the access_token
    digest = hash_fn(access_token.encode("ascii")).digest()

    # Step 2: Take the left-most half of the hash
    half_digest = digest[: len(digest) // 2]

    # Step 3: Base64url encode (without padding)
    token_hash_computed = (
        base64.urlsafe_b64encode(half_digest).rstrip(b"=").decode("ascii")
    )
    return token_hash_computed


def generate_tokens(
    sub: str,
    user_info: dict,
    jwt_exp_s: int,
    refresh_exp_s: int,
    jwt_secret: str,
    jwt_alg: str,
):
    issued_at = int(datetime.utcnow().timestamp())
    jti = str(uuid.uuid4())

    payload_base = {
        **{k: v for k, v in user_info.items() if k not in ["exp", "at_hash"]},
        "sub": sub,
        "iat": issued_at,
    }

    access_token = jwt.encode(
        {
            **payload_base,
            "exp": issued_at + jwt_exp_s,
        },
        jwt_secret,
        algorithm=jwt_alg,
    )

    refresh_token = jwt.encode(
        {
            **user_info,
            "sub": sub,
            "type": "refresh",
            "jti": jti,
            "iat": issued_at,
            "exp": issued_at + refresh_exp_s,
        },
        jwt_secret,
        algorithm=jwt_alg,
    )

    return access_token, refresh_token, issued_at


def respond_with_cookies(
    access_token: str,
    refresh_token: str,
    issued_at: int,
    expires_at: int,
    cookie_name: str,
    cookie_options: dict,
    refresh_cookie_name: str,
    refresh_cookie_options: dict,
):
    res = JSONResponse(
        {
            "success": True,
            "issuedAt": issued_at,
            "expiresAt": issued_at + expires_at,
        }
    )
    res.headers.append(
        "Set-Cookie", build_cookie_string(cookie_name, access_token, cookie_options)
    )
    res.headers.append(
        "Set-Cookie",
        build_cookie_string(refresh_cookie_name, refresh_token, refresh_cookie_options),
    )
    return res
