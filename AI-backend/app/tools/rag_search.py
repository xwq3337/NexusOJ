"""LangChain Tool: 知识库 RAG 检索。"""

from langchain_core.tools import tool

from app.services.milvus_service import rag_search, format_rag_results


@tool
async def knowledge_search(query: str, top_k: int = 3) -> str:
    """搜索算法知识库。当用户询问算法、数据结构、编程概念类问题时使用此工具。

    Args:
        query: 搜索关键词或完整问题
        top_k: 返回结果数量，默认 3
    """
    results = await rag_search(query, top_k)
    return format_rag_results(results, query)
