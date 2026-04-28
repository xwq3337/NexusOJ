"""Milvus 向量数据库服务：Embedding 生成 + 相似度检索。"""

import logging
from typing import Any

import httpx
from pymilvus import MilvusClient

from app.config import settings
from app.deps import get_milvus_client

logger = logging.getLogger(__name__)

# 嵌入向量维度（智谱 embedding-2 输出维度）
EMBEDDING_DIMENSION = 1024

# 简单内存缓存，避免对同一 query 重复计算 embedding
_embedding_cache: dict[str, list[float]] = {}
_CACHE_MAX_SIZE = 1000


async def get_embedding(text: str) -> list[float] | None:
    """调用智谱 API 获取文本的 embedding 向量，带简单缓存。"""
    if text in _embedding_cache:
        return _embedding_cache[text]

    try:
        async with httpx.AsyncClient(timeout=15) as client:
            resp = await client.post(
                f"{settings.GLM_API_BASE_URL}/embeddings",
                headers={
                    "Content-Type": "application/json",
                    "Authorization": f"Bearer {settings.GLM_API_KEY}",
                },
                json={"model": settings.GLM_EMBEDDING_MODEL, "input": text},
            )
            resp.raise_for_status()
            data = resp.json()
            embedding = data["data"][0]["embedding"]

            # 写入缓存
            if len(_embedding_cache) >= _CACHE_MAX_SIZE:
                # 清理一半缓存（FIFO 简化策略）
                keys_to_remove = list(_embedding_cache.keys())[:_CACHE_MAX_SIZE // 2]
                for k in keys_to_remove:
                    del _embedding_cache[k]
            _embedding_cache[text] = embedding

            return embedding
    except Exception as e:
        logger.error("获取 embedding 失败: %s", e)
        return None


def search_knowledge_base(query_embedding: list[float], top_k: int = 3) -> list[dict[str, Any]]:
    """在 Milvus 中执行向量相似度搜索，返回 top_k 条结果。"""
    client: MilvusClient = get_milvus_client()

    results = client.search(
        collection_name=settings.MILVUS_COLLECTION,
        data=[query_embedding],
        anns_field="embedding",
        search_params={"nprobe": 10},
        limit=top_k,
        output_fields=["file_name", "content"],
    )

    hits = results[0] if results else []
    return [
        {
            "file_name": hit["entity"]["file_name"],
            "content": hit["entity"]["content"],
            "score": hit["distance"],
        }
        for hit in hits
    ]


async def rag_search(query: str, top_k: int = 3) -> list[dict[str, Any]]:
    """完整的 RAG 检索流程：query → embedding → Milvus search。"""
    embedding = await get_embedding(query)
    if embedding is None:
        return []
    return search_knowledge_base(embedding, top_k)


def format_rag_results(results: list[dict[str, Any]], query: str) -> str:
    """将检索结果格式化为 LLM 可读的上下文字符串。"""
    if not results:
        return f"未在知识库中找到与 '{query}' 相关的信息"

    lines = [f"根据知识库中找到的相关信息，关于 '{query}' 的内容如下：\n"]
    for r in results:
        lines.append(f"来源文档: {r['file_name']}")
        lines.append(f"内容摘要: {r['content'][:500]}...")
        lines.append(f"相似度得分: {r['score']:.4f}")
        lines.append("-" * 80)
    return "\n".join(lines)
