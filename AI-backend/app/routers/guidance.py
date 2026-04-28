"""POST /personalized-guidance — 个性化学习指导端点。"""

import json
import logging

from fastapi import APIRouter, Depends, Request
from fastapi.responses import StreamingResponse
from langchain_core.messages import HumanMessage, SystemMessage, AIMessage, ToolMessage
from langchain_core.tools import tool
from pydantic import BaseModel

from app.deps import get_optional_user_id
from app.services.llm_service import get_llm, get_llm_for_agent
from app.services.go_backend_client import get_user_ability, get_user_recommendations
from app.tools.rag_search import knowledge_search
from app.tools.problem_lookup import lookup_problem

logger = logging.getLogger(__name__)

router = APIRouter()


class GuidanceRequest(BaseModel):
    question: str = ""
    messages: list = []


def _sse_event(data: dict) -> str:
    return f"data: {json.dumps(data, ensure_ascii=False)}\n\n"


@router.post("/personalized-guidance")
async def personalized_guidance(
    request: Request,
    user_id: int = Depends(get_optional_user_id),
):
    """基于用户能力画像的个性化学习指导，SSE 流式返回。"""
    try:
        body = await request.json()
    except Exception:
        from fastapi.responses import JSONResponse
        return JSONResponse({"error": "invalid JSON"}, status_code=400)

    question = body.get("question", "")
    messages = body.get("messages", [])

    # 提取 token
    auth_header = request.headers.get("authorization", "")
    token = auth_header.replace("Bearer ", "") if auth_header.startswith("Bearer ") else auth_header

    # 先获取用户能力数据
    yield_direct = False
    ability_text = "用户能力数据暂无"

    if token:
        ability = await get_user_ability(token)
        if ability:
            tag_scores = ability.get("tag_scores", {})
            strongest = ability.get("strongest_tags", [])
            weakest = ability.get("weakest_tags", [])
            overall = ability.get("overall_score", 0)

            lines = [
                f"综合能力评分: {overall:.2f}",
                "各知识点掌握度:",
            ]
            for tag, score in sorted(tag_scores.items(), key=lambda x: x[1]):
                bar = "█" * int(score * 10) + "░" * (10 - int(score * 10))
                lines.append(f"  {tag}: [{bar}] {score:.2f}")
            if strongest:
                lines.append(f"最强领域: {', '.join(strongest)}")
            if weakest:
                lines.append(f"薄弱领域: {', '.join(weakest)}")

            ability_text = "\n".join(lines)

            # 同时获取推荐题目
            recommendations = await get_user_recommendations(token, page=1, page_size=5)
            if recommendations and recommendations.get("problems"):
                lines.append("\n系统推荐的练习题目:")
                for p in recommendations["problems"][:5]:
                    lines.append(f"  - {p.get('title', '')}（难度 {p.get('difficulty', 0):.1f}，原因: {p.get('reason', '')}）")
                ability_text = "\n".join(lines)

    system_prompt = f"""你是一位个性化的编程学习导师。你拥有用户的详细能力画像数据。

以下是当前用户的能力分析：
{ability_text}

基于以上数据，请为用户提供：
1. 学习状况分析：当前水平和能力分布
2. 薄弱环节诊断：哪些知识点需要重点加强
3. 学习路径建议：推荐的具体练习方向和题目类型
4. 目标设定：短期可达成的学习目标

请用友好、鼓励的语气，用中文回答。给出具体可操作的建议，而非笼统的方向。"""

    user_question = question or "请分析我的学习状况并给出练习建议。"

    llm = get_llm()

    async def event_generator():
        try:
            lc_messages = [SystemMessage(content=system_prompt)]

            # 添加对话历史
            for m in messages[-10:]:  # 保留最近 10 条
                role = m.get("role", "user")
                content = m.get("content", "")
                if role in ("assistant", "model"):
                    lc_messages.append(AIMessage(content=content))
                else:
                    lc_messages.append(HumanMessage(content=content))

            lc_messages.append(HumanMessage(content=user_question))

            async for chunk in llm.astream(lc_messages):
                if chunk.content:
                    yield _sse_event({"text": chunk.content})
            yield _sse_event({"done": True})
        except Exception as e:
            logger.error("Guidance stream error: %s", e, exc_info=True)
            yield _sse_event({"error": str(e)})

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={"Cache-Control": "no-cache", "X-Accel-Buffering": "no"},
    )
