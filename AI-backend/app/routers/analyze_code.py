"""POST /analyze-code — 代码分析端点（复杂度、风格、质量）。"""

import json
import logging

from fastapi import APIRouter, Depends
from fastapi.responses import StreamingResponse
from langchain_core.messages import HumanMessage, SystemMessage
from pydantic import BaseModel

from app.deps import get_optional_user_id
from app.services.llm_service import get_llm, load_system_prompt

logger = logging.getLogger(__name__)

router = APIRouter()


class AnalyzeCodeRequest(BaseModel):
    code: str
    language: str = "cpp"
    analysis_type: str = "all"  # complexity | style | quality | all


def _sse_event(data: dict) -> str:
    return f"data: {json.dumps(data, ensure_ascii=False)}\n\n"


@router.post("/analyze-code")
async def analyze_code_stream(
    req: AnalyzeCodeRequest,
    user_id: int = Depends(get_optional_user_id),
):
    """对用户提交的代码进行全面分析，SSE 流式返回结果。"""

    type_prompts = {
        "complexity": "时间复杂度和空间复杂度分析，给出 Big-O 估计并解释推导过程",
        "style": "代码风格评估，包括命名规范、缩进格式、注释完整度、代码组织",
        "quality": "代码质量评估，包括潜在 bug、边界条件处理、内存安全、可读性",
        "all": "全面的代码审查，包括复杂度分析、代码风格评估和质量评估",
    }

    focus = type_prompts.get(req.analysis_type, type_prompts["all"])

    system_prompt = f"""你是一位资深代码审查专家。请对用户提交的代码进行{focus}。

输出格式要求：
1. **复杂度分析** - 时间复杂度 O(?) 和空间复杂度 O(?)，附推导说明
2. **代码风格** - 命名规范、格式化、注释（评分 1-10）
3. **质量评估** - 潜在问题列表、边界处理评价、改进建议
4. **综合评分** - 总体评价（优秀/良好/需改进）及关键改进点

请用中文回答，简洁专业，重点突出。"""

    llm = get_llm()

    async def event_generator():
        try:
            messages = [
                SystemMessage(content=system_prompt),
                HumanMessage(content=f"请分析以下 {req.language} 代码：\n```{req.language}\n{req.code}\n```"),
            ]
            async for chunk in llm.astream(messages):
                if chunk.content:
                    yield _sse_event({"text": chunk.content})
            yield _sse_event({"done": True})
        except Exception as e:
            logger.error("Analyze code stream error: %s", e, exc_info=True)
            yield _sse_event({"error": str(e)})

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={"Cache-Control": "no-cache", "X-Accel-Buffering": "no"},
    )
