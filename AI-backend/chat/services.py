import os
import json
import requests
from pymilvus import MilvusClient

COLLECTION_NAME = "algorithm_knowledge"
GLM_API_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"

_milvus_client = None


def get_milvus_client():
    global _milvus_client
    if _milvus_client is None:
        uri = f"http://{os.getenv('MILVUS_HOST')}:{os.getenv('MILVUS_PORT')}"
        _milvus_client = MilvusClient(uri=uri)
    return _milvus_client


def get_embedding(text):
    try:
        resp = requests.post(
            "https://open.bigmodel.cn/api/paas/v4/embeddings",
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {os.getenv('GLM_API_KEY')}",
            },
            json={"model": "embedding-2", "input": text},
            timeout=10,
        )
        resp.raise_for_status()
        data = resp.json()
        return data["data"][0]["embedding"]
    except Exception as e:
        print(f"获取embedding失败: {e}")
        return None


def rag_search(query, top_k=3):
    client = get_milvus_client()
    query_embedding = get_embedding(query)
    if query_embedding is None:
        return []

    results = client.search(
        collection_name=COLLECTION_NAME,
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


def format_rag_context(query):
    try:
        results = rag_search(query)
        if not results:
            return f"未在知识库中找到与 '{query}' 相关的信息"

        lines = [f"根据知识库中找到的相关信息，关于 '{query}' 的内容如下：\n"]
        for r in results:
            lines.append(f"来源文档: {r['file_name']}")
            lines.append(f"内容摘要: {r['content'][:500]}...")
            lines.append(f"相似度得分: {r['score']:.4f}")
            lines.append("-" * 80)

        return "\n".join(lines)
    except Exception as e:
        return f"检索过程中发生错误: {e}"


def stream_chat(messages):
    """生成器：逐块 yield ZhipuAI 流式 API 的内容文本。"""
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {os.getenv('GLM_API_KEY')}",
    }
    payload = {
        "model": os.getenv("GLM_MODEL_NAME", "GLM-4.7-FlashX"),
        "messages": messages,
        "temperature": 0.5,
        "top_p": 0.9,
        "stream": True,
    }

    resp = requests.post(
        GLM_API_URL, headers=headers, json=payload, stream=True, timeout=120
    )
    resp.raise_for_status()

    for line in resp.iter_lines(decode_unicode=True):
        if not line:
            continue
        if line.startswith("data: "):
            data_str = line[6:]
            if data_str.strip() == "[DONE]":
                break
            try:
                chunk = json.loads(data_str)
                delta = chunk.get("choices", [{}])[0].get("delta", {})
                content = delta.get("content", "")
                if content:
                    yield content
            except json.JSONDecodeError:
                continue
