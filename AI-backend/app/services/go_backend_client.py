"""Go 主后端 API 客户端：httpx 异步调用，自动附加 JWT。"""

import logging
from typing import Any

import httpx

from app.config import settings

logger = logging.getLogger(__name__)

# 复用连接池
_client: httpx.AsyncClient | None = None


async def get_client() -> httpx.AsyncClient:
    """获取 httpx 异步客户端单例。"""
    global _client
    if _client is None:
        _client = httpx.AsyncClient(
            base_url=settings.GO_BACKEND_URL,
            timeout=10.0,
        )
    return _client


async def _request(method: str, path: str, token: str, **kwargs) -> dict[str, Any] | None:
    """通用请求方法，自动附加 Authorization header。"""
    client = await get_client()
    headers = {"Authorization": f"Bearer {token}"}

    try:
        resp = await client.request(method, path, headers=headers, **kwargs)
        resp.raise_for_status()
        data = resp.json()
        # Go 后端统一响应格式：{"code": 200, "data": {...}}
        if data.get("code") == 200:
            return data.get("data") or data.get("info")
        logger.warning("Go 后端返回非 200: %s", data)
        return None
    except httpx.TimeoutException:
        logger.error("Go 后端请求超时: %s %s", method, path)
        return None
    except httpx.HTTPStatusError as e:
        logger.error("Go 后端 HTTP 错误: %s %s -> %s", method, path, e.response.status_code)
        return None
    except Exception as e:
        logger.error("Go 后端请求异常: %s %s -> %s", method, path, e)
        return None


# ==================== 用户画像 ====================

async def get_user_ability(token: str) -> dict[str, Any] | None:
    """获取用户能力分析（各标签掌握度）。"""
    return await _request("GET", "/recommend/ability", token)


async def get_user_profile(token: str) -> dict[str, Any] | None:
    """获取用户完整画像。"""
    return await _request("GET", "/recommend/profile", token)


async def get_user_recommendations(token: str, page: int = 1, page_size: int = 10) -> dict[str, Any] | None:
    """获取推荐题目。"""
    return await _request("GET", f"/recommend/problems?page={page}&page_size={page_size}", token)


# ==================== 题目相关 ====================

async def get_problem_detail(problem_id: int, token: str) -> dict[str, Any] | None:
    """获取题目详情。"""
    return await _request("GET", f"/problem/{problem_id}", token)
