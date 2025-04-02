from dotenv import load_dotenv

load_dotenv()

from fastapi import APIRouter, FastAPI
from app.routers.receipt import receipt
from app.routers.outing import outing
from app.routers.auth import auth
from starlette.middleware.cors import CORSMiddleware
from fastapi.templating import Jinja2Templates

api_router = APIRouter()
api_router.include_router(receipt.router)
api_router.include_router(outing.router)
api_router.include_router(auth.router)

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8081"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


app.include_router(api_router, prefix="/api/v1")
