from datetime import datetime, time
import os
from typing import Optional
from urllib.parse import urlencode
import uuid
from fastapi import APIRouter, Cookie, Form, HTTPException, Header, Request, Response
import httpx
import requests
from jose import ExpiredSignatureError, JWTError, jwt, jws

from app.routers.auth.utils import (
    build_cookie_string,
    compute_at_hash,
    generate_tokens,
    parse_cookie_header,
    respond_with_cookies,
)
from app.routers.deps import OAuth2Dep
from fastapi.responses import JSONResponse, RedirectResponse


router = APIRouter(prefix="/auth", tags=["items"])

GOOGLE_CLIENT_ID = os.environ["GOOGLE_CLIENT_ID"]
GOOGLE_CLIENT_SECRET = os.environ["GOOGLE_CLIENT_SECRET"]
GOOGLE_REDIRECT_URI = os.environ["GOOGLE_REDIRECT_URI"]
JWT_SECRET = os.environ["JWT_SECRET"]
JWT_ALGORITHM = "HS256"
JWT_EXPIRATION_TIME_SECONDS = 60 * 15  # 15 minutes
REFRESH_EXPIRY_SECONDS = 60 * 60 * 24 * 7  # 7 days
COOKIE_NAME = "auth_token"
REFRESH_COOKIE_NAME = "refresh_token"
GOOGLE_AUTH_URL = "https://accounts.google.com/o/oauth2/auth"
BASE_URL = os.getenv("BASE_URL", "http://localhost:8081")
SERVER_URL = os.getenv("SERVER_URL", "http://localhost:8000")
APP_SCHEME = os.getenv("APP_SCHEME", "myapp")

COOKIE_OPTIONS = {
    "max_age": JWT_EXPIRATION_TIME_SECONDS,
    "path": "/",
    "httponly": True,
    "secure": False,  # Set to True in production with HTTPS
    "samesite": "Lax",
}

REFRESH_COOKIE_OPTIONS = {
    "max_age": REFRESH_EXPIRY_SECONDS,
    "path": "/",
    "httponly": True,
    "secure": False,
    "samesite": "Lax",
}


@router.get("/login/google")
async def login_google(request: Request):
    if not GOOGLE_CLIENT_ID:
        return JSONResponse(
            {"error": "Missing GOOGLE_CLIENT_ID environment variable"},
            status_code=500,
        )

    query_params = dict(request.query_params)
    internal_client = query_params.get("client_id")
    redirect_uri = query_params.get("redirect_uri")
    requested_scope = query_params.get("scope", "identity")
    state_param = query_params.get("state")

    print("Red: ", redirect_uri)

    # Determine platform
    if "expo" in redirect_uri:
        platform = "mobile"
    elif redirect_uri == BASE_URL:
        platform = "web"
    else:
        return JSONResponse({"error": "Invalid redirect_uri"}, status_code=400)

    # Validate internal client
    if internal_client != "google":
        return JSONResponse({"error": "Invalid client"}, status_code=400)

    # Validate state
    if not state_param:
        return JSONResponse({"error": "Invalid state"}, status_code=400)

    # Construct state with platform for round-trip
    state = f"{platform}|{state_param}"

    # Build final redirect URL
    query = urlencode(
        {
            "client_id": GOOGLE_CLIENT_ID,
            "redirect_uri": f"{SERVER_URL}/api/v1/auth/callback",
            "response_type": "code",
            "scope": requested_scope,
            "state": state,
            "prompt": "select_account",
        }
    )

    return RedirectResponse(url=f"{GOOGLE_AUTH_URL}?{query}")


@router.get("/callback")
async def auth_callback(request: Request):
    query_params = dict(request.query_params)

    combined_state = query_params.get("state")
    if not combined_state:
        return JSONResponse({"error": "Invalid state"}, status_code=400)

    try:
        platform, original_state = combined_state.split("|", 1)
    except ValueError:
        return JSONResponse({"error": "Malformed state"}, status_code=400)

    code = query_params.get("code", "")

    redirect_url = f"{BASE_URL if platform == 'web' else APP_SCHEME}?{urlencode({'code': code, 'state': original_state})}"

    return RedirectResponse(url=redirect_url)


@router.get("/google")
async def auth_google(code: str):
    token_url = "https://accounts.google.com/o/oauth2/token"
    data = {
        "code": code,
        "client_id": GOOGLE_CLIENT_ID,
        "client_secret": GOOGLE_CLIENT_SECRET,
        "redirect_uri": GOOGLE_REDIRECT_URI,
        "grant_type": "authorization_code",
    }
    response = requests.post(token_url, data=data)
    access_token = response.json().get("access_token")
    user_info = requests.get(
        "https://www.googleapis.com/oauth2/v1/userinfo",
        headers={"Authorization": f"Bearer {access_token}"},
    )
    return data


@router.post("/token")
async def google_auth_token(code: str = Form(...), platform: str = Form("native")):
    if not code:
        return JSONResponse({"error": "Missing authorization code"}, status_code=400)

    async with httpx.AsyncClient() as client:
        token_response = await client.post(
            "https://oauth2.googleapis.com/token",
            headers={"Content-Type": "application/x-www-form-urlencoded"},
            data={
                "client_id": GOOGLE_CLIENT_ID,
                "client_secret": GOOGLE_CLIENT_SECRET,
                "redirect_uri": GOOGLE_REDIRECT_URI,
                "grant_type": "authorization_code",
                "code": code,
            },
        )

    token_data = token_response.json()
    if "error" in token_data:
        return JSONResponse(
            {
                "error": token_data["error"],
                "error_description": token_data.get("error_description"),
                "message": "OAuth validation error - please ensure the app complies with Google's OAuth 2.0 policy",
            },
            status_code=400,
        )

    if "id_token" not in token_data:
        return JSONResponse({"error": "Missing required parameters"}, status_code=400)

    user_info = jwt.get_unverified_claims(token_data["id_token"])

    sub = user_info.get("sub")
    if not sub:
        return JSONResponse({"error": "Missing sub in ID token"}, status_code=400)

    access_token, refresh_token, issued_at = generate_tokens(
        sub=sub,
        user_info=user_info,
        jwt_alg=JWT_ALGORITHM,
        jwt_exp_s=JWT_EXPIRATION_TIME_SECONDS,
        jwt_secret=JWT_SECRET,
        refresh_exp_s=REFRESH_EXPIRY_SECONDS,
    )

    if platform == "web":
        # Build and return the response with cookies
        res = JSONResponse(
            {
                "success": True,
                "issuedAt": issued_at,
                "expiresAt": issued_at + COOKIE_OPTIONS["max_age"],
            }
        )
        res.headers.append(
            "Set-Cookie",
            build_cookie_string(COOKIE_NAME, access_token, COOKIE_OPTIONS),
        )
        res.headers.append(
            "Set-Cookie",
            build_cookie_string(
                REFRESH_COOKIE_NAME, refresh_token, REFRESH_COOKIE_OPTIONS
            ),
        )
        return res

    # For native clients
    return {
        "accessToken": access_token,
        "refreshToken": refresh_token,
    }


@router.get("/session")
async def get_session(request: Request):
    cookie_header = request.headers.get("cookie")
    if not cookie_header:
        raise HTTPException(status_code=401, detail="Not authenticated")

    cookies = parse_cookie_header(cookie_header)

    if COOKIE_NAME not in cookies or "value" not in cookies[COOKIE_NAME]:
        raise HTTPException(status_code=401, detail="Not authenticated")

    token = cookies[COOKIE_NAME]["value"]

    try:
        verified = jwt.decode(
            token, JWT_SECRET, algorithms=[JWT_ALGORITHM], audience=GOOGLE_CLIENT_ID
        )
    except JWTError as e:
        print("ERROR JWT: ", e)
        raise HTTPException(status_code=401, detail="Invalid token")

    cookie_expiration = None
    if "maxAge" in cookies[COOKIE_NAME]:
        try:
            max_age = int(cookies[COOKIE_NAME]["maxAge"])
            issued_at = verified.get(
                "iat", int(datetime.now(datetime.timezone.utc).timestamp())
            )
            cookie_expiration = issued_at + max_age
        except ValueError:
            pass

    return {**verified, "cookieExpiration": cookie_expiration}


@router.post("/refresh")
async def refresh_token(
    request: Request,
    response: Response,
    refresh_token_body: Optional[str] = Form(None),
    platform: Optional[str] = Form("native"),
    authorization: Optional[str] = Header(None),
    cookie_refresh: Optional[str] = Cookie(None, alias=REFRESH_COOKIE_NAME),
):
    # Determine platform from query or form
    platform = request.query_params.get("platform", platform or "native")
    refresh_token = refresh_token_body if platform == "native" else cookie_refresh

    # Fallback: use access token from Authorization header
    if not refresh_token and authorization and authorization.startswith("Bearer "):
        access_token = authorization.split(" ")[1]
        try:
            decoded = jwt.decode(access_token, JWT_SECRET, algorithms=[JWT_ALGORITHM])
            sub = decoded.get("sub")
            issued_at = int(datetime.utcnow().timestamp())

            new_access_token = jwt.encode(
                {
                    **decoded,
                    "exp": issued_at + JWT_EXPIRATION_TIME_SECONDS,
                    "iat": issued_at,
                },
                JWT_SECRET,
                algorithm=JWT_ALGORITHM,
            )

            if platform == "web":
                response.headers["Set-Cookie"] = build_cookie_string(
                    COOKIE_NAME, new_access_token, COOKIE_OPTIONS
                )
                return {"success": True, "warning": "Access token used as fallback"}

            return {"accessToken": new_access_token, "warning": "Access token fallback"}

        except JWTError:
            raise HTTPException(status_code=401, detail="Invalid or expired token")

    # Validate refresh token
    try:
        decoded = jwt.decode(refresh_token, JWT_SECRET, algorithms=[JWT_ALGORITHM])
    except ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Refresh token expired")
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid refresh token")
    except Exception:
        raise HTTPException(status_code=401, detail="Unknown error decoding token")

    if decoded.get("type") != "refresh":
        raise HTTPException(status_code=401, detail="Invalid token type")

    sub = decoded.get("sub")
    if not sub:
        raise HTTPException(status_code=401, detail="Missing user subject")

    user_info = {
        "sub": sub,
        "name": decoded.get("name") or f"User {sub[:6]}",
        "email": decoded.get("email") or f"user-{sub[:6]}@example.com",
        "picture": decoded.get("picture") or "https://ui-avatars.com/api/?name=User",
    }

    access_token, refresh_token, issued_at = generate_tokens(sub, user_info)

    if platform == "web":
        return respond_with_cookies(access_token, refresh_token, issued_at)

    return {
        "accessToken": access_token,
        "refreshToken": refresh_token,
    }
