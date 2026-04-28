"""POST /chat — 增强型对话端点，使用 LangChain ReAct Agent + Tool Calling。"""

import json
import logging
from typing import Any

from fastapi import APIRouter, Depends, Request
from fastapi.responses import StreamingResponse
from langchain_core.messages import HumanMessage, SystemMessage, AIMessage
from langgraph.prebuilt import create_react_agent

from app.config import settings
from app.deps import get_optional_user_id
from app.services.llm_service import get_llm_for_agent, load_system_prompt
from app.tools.rag_search import knowledge_search
from app.tools.code_analyzer import analyze_code
from app.tools.test_generator import generate_test_cases
from app.tools.user_profile import get_user_ability_profile
from app.tools.problem_lookup import lookup_problem

logger = logging.getLogger(__name__)

router = APIRouter()

# Agent 工具列表
ALL_TOOLS = [
    knowledge_search,
    analyze_code,
    generate_test_cases,
    get_user_ability_profile,
    lookup_problem,
]


def _inject_token_into_tool_args(tool_call: dict, token: str) -> dict:
    """向需要认证的 tool 参数中注入 JWT token。"""
    args = tool_call.get("args", {})
    if "token" in args.__fields__ if hasattr(args, "__fields__") else False:
        args["token"] = token
    return args


def _sse_event(data: dict) -> str:
    """将字典序列化为 SSE data 帧格式。"""
    return f"data: {json.dumps(data, ensure_ascii=False)}\n\n"


@router.post("/chat")
async def chat_stream(
    request: Request,
    user_id: int = Depends(get_optional_user_id),
):
    """SSE 流式对话端点。LLM 自主决定调用哪些 tool。"""
    try:
        body = await request.json()
    except Exception:
        from fastapi.responses import JSONResponse
        return JSONResponse({"error": "invalid JSON"}, status_code=400)

    messages = body.get("messages", [])
    if not messages or not isinstance(messages, list):
        from fastapi.responses import JSONResponse
        return JSONResponse({"error": "messages is required"}, status_code=400)

    # 从 request 中提取原始 Authorization header，转发到 Go 后端
    auth_header = request.headers.get("authorization", "")
    token = auth_header.replace("Bearer ", "") if auth_header.startswith("Bearer ") else auth_header

    # 加载系统提示词
    system_prompt = load_system_prompt("system_prompt.txt")

    async def event_generator():
        try:
            # 构建 LangChain 消息列表
            lc_messages = [SystemMessage(content=system_prompt)]
            for m in messages:
                role = m.get("role", "user")
                content = m.get("content", "")
                # 前端用 "model" 表示 assistant
                if role in ("assistant", "model"):
                    lc_messages.append(AIMessage(content=content))
                elif role == "system":
                    lc_messages.append(SystemMessage(content=content))
                else:
                    lc_messages.append(HumanMessage(content=content))

            # 创建 ReAct Agent
            llm = get_llm_for_agent()

            # 使用简单的 LLM 直接流式调用（带 tool 绑定）
            llm_with_tools = llm.bind_tools(
                [
                    {
                        "type": "function",
                        "function": {
                            "name": tool.name,
                            "description": tool.description,
                            "parameters": {
                                "type": "object",
                                "properties": tool.args_schema.schema().get("properties", {}),
                                "required": tool.args_schema.schema().get("required", []),
                            },
                        },
                    }
                    for tool in ALL_TOOLS
                ]
            )

            # 第一次调用：让 LLM 决定是否调用 tool
            yield _sse_event({"status": "thinking"})

            first_response = await llm_with_tools.ainvoke(lc_messages)

            # 处理 tool calls
            if hasattr(first_response, "tool_calls") and first_response.tool_calls:
                lc_messages.append(first_response)

                for tool_call in first_response.tool_calls:
                    tool_name = tool_call["name"]
                    tool_args = dict(tool_call["args"])

                    # 注入 token 到需要认证的 tool
                    if "token" not in tool_args:
                        tool_args["token"] = token

                    yield _sse_event({"status": "tool_call", "tool": tool_name})

                    # 执行 tool
                    tool_map = {t.name: t for t in ALL_TOOLS}
                    selected_tool = tool_map.get(tool_name)
                    if selected_tool:
                        try:
                            tool_result = await selected_tool.ainvoke(tool_args)
                            yield _sse_event({"status": "tool_result", "tool": tool_name})
                        except Exception as e:
                            logger.error("Tool %s 执行失败: %s", tool_name, e)
                            tool_result = f"工具 {tool_name} 执行出错: {e}"
                    else:
                        tool_result = f"未知工具: {tool_name}"

                    # 将 tool 结果添加到消息中
                    from langchain_core.messages import ToolMessage
                    lc_messages.append(ToolMessage(content=str(tool_result), tool_call_id=tool_call["id"]))

                # 带 tool 结果再次调用 LLM 生成最终回答
                async for chunk in llm.astream(lc_messages):
                    if chunk.content:
                        yield _sse_event({"text": chunk.content})
            else:
                # 无 tool call，直接流式输出（对首次响应已完整返回的情况）
                if first_response.content:
                    # 已经有了完整响应，分块发送模拟流式效果
                    content = first_response.content
                    chunk_size = max(1, len(content) // 20)
                    for i in range(0, len(content), chunk_size):
                        yield _sse_event({"text": content[i:i + chunk_size]})
                else:
                    # 无内容则重新流式调用
                    async for chunk in llm.astream(lc_messages):
                        if chunk.content:
                            yield _sse_event({"text": chunk.content})

            yield _sse_event({"done": True})

        except Exception as e:
            logger.error("Chat stream error: %s", e, exc_info=True)
            yield _sse_event({"error": str(e)})

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={"Cache-Control": "no-cache", "X-Accel-Buffering": "no"},
    )
