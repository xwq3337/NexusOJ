"""FastAPI 依赖注入：JWT 解码、Milvus 客户端单例、LLM 实例。"""

import jwt
from fastapi import Depends, Header, HTTPException
from pymilvus import MilvusClient

from app.config import settings


# ==================== JWT 解码 ====================

def get_current_user_id(authorization: str = Header(..., alias="Authorization")) -> int:
    """从 Authorization header 解析 JWT，提取 userID。所有 AI 请求必须带 JWT。"""
    if not authorization:
        raise HTTPException(status_code=401, detail="缺少 Authorization header")

    parts = authorization.split(" ")
    if len(parts) != 2 or parts[0].lower() != "bearer":
        raise HTTPException(status_code=401, detail="Authorization header 格式错误，应为 Bearer <token>")

    token = parts[1]
    try:
        # 与 Go 后端使用相同的签名密钥和算法
        payload = jwt.decode(token, settings.JWT_SIGN_KEY, algorithms=["HS256"], options={"verify_exp": False})
        user_id = payload.get("userID", 0)
        if user_id == 0:
            raise HTTPException(status_code=401, detail="无效的用户 ID")
        return int(user_id)
    except jwt.InvalidTokenError as e:
        raise HTTPException(status_code=401, detail=f"JWT 解码失败: {e}")


# 可选的 JWT 依赖（不强制要求登录）
def get_optional_user_id(authorization: str = Header(default="", alias="Authorization")) -> int:
    """尝试解析 JWT，失败返回 0。"""
    if not authorization:
        return 0
    parts = authorization.split(" ")
    if len(parts) != 2 or parts[0].lower() != "bearer":
        return 0
    try:
        payload = jwt.decode(parts[1], settings.JWT_SIGN_KEY, algorithms=["HS256"], options={"verify_exp": False})
        return int(payload.get("userID", 0))
    except jwt.InvalidTokenError:
        return 0


# 从请求体 JSON 中提取 user_id（由 Go Backend 注入）
def get_user_id_from_body(request) -> int:
    """从请求体 JSON 中提取 user_id 字段（Go Backend 代理时注入）。"""
    try:
        body = request.json()
        return int(body.get("user_id", 0))
    except Exception:
        return 0


# ==================== Milvus 客户端单例 ====================

_milvus_client: MilvusClient | None = None


def get_milvus_client() -> MilvusClient:
    """获取 Milvus 客户端单例。应用启动时通过 lifespan 初始化。"""
    global _milvus_client
    if _milvus_client is None:
        uri = f"http://{settings.MILVUS_HOST}:{settings.MILVUS_PORT}"
        _milvus_client = MilvusClient(uri=uri)
    return _milvus_client


def init_milvus_client() -> MilvusClient:
    """应用 lifespan 调用，初始化 Milvus 连接。"""
    global _milvus_client
    uri = f"http://{settings.MILVUS_HOST}:{settings.MILVUS_PORT}"
    _milvus_client = MilvusClient(uri=uri)
    return _milvus_client
