"""FastAPI 应用入口：lifespan 管理 Milvus 连接，挂载所有路由。"""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.deps import init_milvus_client
from app.routers import analyze_code, chat, generate_tests, guidance

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(_app: FastAPI):
    """应用启动/关闭时的生命周期管理。"""
    # 启动：尝试初始化 Milvus 连接（非阻塞，连接失败不影响其他功能）
    try:
        init_milvus_client()
        print(f"[启动] Milvus 连接已建立 ({settings.MILVUS_HOST}:{settings.MILVUS_PORT})")
    except Exception as e:
        logger.warning("Milvus 连接失败（RAG 功能不可用）: %s", e)
        print(f"[警告] Milvus 连接失败，RAG 功能将不可用: {e}")
    yield
    # 关闭：清理资源
    print("[关闭] 清理完成")


app = FastAPI(
    title="NexusOJ AI Backend",
    version="2.0.0",
    lifespan=lifespan,
)

# CORS 中间件
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(chat.router, tags=["chat"])
app.include_router(analyze_code.router, tags=["analyze-code"])
app.include_router(generate_tests.router, tags=["generate-tests"])
app.include_router(guidance.router, tags=["guidance"])


@app.get("/health")
async def health_check():
    return {"status": "ok", "version": "2.0.0"}
