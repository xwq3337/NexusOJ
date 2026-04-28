"""LLM 服务：封装 LangChain ChatOpenAI 实例和流式调用工具。"""

import logging
from typing import AsyncGenerator

from langchain_openai import ChatOpenAI

from app.config import settings

logger = logging.getLogger(__name__)

# LLM 单例
_llm: ChatOpenAI | None = None
_llm_for_agent: ChatOpenAI | None = None


def get_llm() -> ChatOpenAI:
    """获取基础 LLM 实例（用于直接调用，不含 tool 绑定）。"""
    global _llm
    if _llm is None:
        _llm = ChatOpenAI(
            model=settings.GLM_MODEL_NAME,
            api_key=settings.GLM_API_KEY,
            base_url=settings.GLM_API_BASE_URL,
            temperature=0.5,
            top_p=0.9,
            streaming=True,
        )
    return _llm


def get_llm_for_agent() -> ChatOpenAI:
    """获取用于 Agent 的 LLM 实例（需要支持 tool calling）。"""
    global _llm_for_agent
    if _llm_for_agent is None:
        _llm_for_agent = ChatOpenAI(
            model=settings.GLM_MODEL_NAME,
            api_key=settings.GLM_API_KEY,
            base_url=settings.GLM_API_BASE_URL,
            temperature=0.5,
            top_p=0.9,
            streaming=True,
        )
    return _llm_for_agent


def load_system_prompt(filename: str) -> str:
    """从 prompts/ 目录加载系统提示词文件。"""
    from pathlib import Path

    prompts_dir = Path(__file__).resolve().parent.parent / "prompts"
    file_path = prompts_dir / filename
    if not file_path.exists():
        logger.warning("提示词文件不存在: %s", file_path)
        return ""
    return file_path.read_text("utf-8")


async def stream_llm_response(messages: list[dict]) -> AsyncGenerator[str, None]:
    """调用 LLM 并流式 yield 文本块。

    Args:
        messages: OpenAI 格式的消息列表 [{"role": "...", "content": "..."}]

    Yields:
        文本块字符串
    """
    from langchain_core.messages import HumanMessage, SystemMessage, AIMessage

    lc_messages = []
    for msg in messages:
        role = msg.get("role", "user")
        content = msg.get("content", "")
        if role == "system":
            lc_messages.append(SystemMessage(content=content))
        elif role == "assistant":
            lc_messages.append(AIMessage(content=content))
        else:
            lc_messages.append(HumanMessage(content=content))

    llm = get_llm()
    async for chunk in llm.astream(lc_messages):
        if chunk.content:
            yield chunk.content
