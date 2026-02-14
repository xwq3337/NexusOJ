"""RAG工具，用于检索相关文档并提供给AI模型"""
import os
from typing import Dict, List
import asyncio

import os
from pymilvus import connections, Collection
from zhipuai import ZhipuAI
from loguru import logger
from dotenv import load_dotenv
from langchain.tools import BaseTool
from pydantic import BaseModel, Field
from loguru import logger

# ZhipuAI配置
ZHIPU_API_KEY = os.getenv("ZHIPU_API_KEY")
client = ZhipuAI(api_key=ZHIPU_API_KEY)

# 集合名称
COLLECTION_NAME = "algorithm_knowledge"

logger.add("log/rag_search.log", rotation="10 MB", retention="10 days", level="INFO")
load_dotenv()



class RAGInput(BaseModel):
    query: str = Field(description="用户查询内容")


class RAGTool(BaseTool):
    name: str = "rag_search"
    description: str = "用于检索与用户查询相关的知识库文档，适用于回答关于算法、编程等技术问题"
    args_schema: type = RAGInput

    def _run(self, query: str) -> str:
        """执行RAG搜索"""
        try:
            # 连接到Milvus
            if not connect_milvus():
                return "无法连接到知识库，请稍后重试"
            
            # 搜索相关文档
            results = search(query, top_k=3)
            
            if not results:
                return f"未在知识库中找到与 '{query}' 相关的信息"
            
            # 构建返回结果
            rag_context = f"根据知识库中找到的相关信息，关于 '{query}' 的内容如下：\n\n"
            
            for idx, result in enumerate(results, 1):
                rag_context += f"来源文档: {result['file_name']}\n"
                rag_context += f"内容摘要: {result['content'][:500]}...\n"
                rag_context += f"相似度得分: {result['score']:.4f}\n"
                rag_context += "-" * 80 + "\n"
            
            return rag_context
            
        except Exception as e:
            return f"检索过程中发生错误: {str(e)}"
    
    async def _arun(self, query: str):
        """异步版本"""
        raise NotImplementedError("RAGTool不支持异步操作")


def connect_milvus():
    """连接到Milvus数据库"""
    try:
        # 阿里云Milvus使用host:port的方式连接
        logger.info(f"尝试连接Milvus: {os.getenv('MILVUS_HOST')}:{os.getenv('MILVUS_PORT')}")

        connections.connect(
            alias="default",
            host=os.getenv("MILVUS_HOST", "localhost"),
            port=os.getenv("MILVUS_PORT", "19530"),
            user=os.getenv("MILVUS_USER", ""),
            password=os.getenv("MILVUS_PASSWORD", ""),
        )
        logger.info(f"成功连接到Milvus数据库")
        return True
    except Exception as e:
        logger.error(f"连接Milvus失败: {e}")
        return False


def get_embedding(text):
    """使用ZhipuAI获取文本的向量表示"""
    try:
        response = client.embeddings.create(
            model="embedding-2",
            input=text
        )
        embedding = response.data[0].embedding
        return embedding
    except Exception as e:
        logger.error(f"获取embedding失败: {e}")
        return None


def search(query, top_k=5):
    """在Milvus中搜索相似新闻"""
    try:
        # 获取集合
        collection = Collection(name=COLLECTION_NAME)

        # 确保集合已加载
        collection.load()

        # 获取查询文本的向量
        query_embedding = get_embedding(query)

        if not query_embedding:
            logger.error("无法获取查询文本的向量")
            return []

        # 定义搜索参数
        search_params = {
            "metric_type": "COSINE",
            "params": {"nprobe": 10}
        }

        # 执行搜索
        search_coro = collection.search(
            data=[query_embedding],
            anns_field="embedding",
            param=search_params,
            limit=top_k,
            output_fields=["file_name", "content"]
        )

        # 如果返回的是协程，运行它获取结果
        if asyncio.iscoroutine(search_coro):
            results = asyncio.run(search_coro)
        else:
            results = search_coro

        # 处理结果
        search_results = []
        for hits in results:
            for hit in hits:
                search_results.append({
                    "file_name": hit.entity.get("file_name"),
                    "content": hit.entity.get("content"),
                    "distance": hit.distance,
                    "score": hit.score
                })

        return search_results

    except Exception as e:
        logger.error(f"搜索失败: {e}")
        return []


# 创建RAG工具实例
rag_tool = RAGTool()




