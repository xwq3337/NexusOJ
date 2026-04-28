"""POST /generate-tests — 测试用例生成端点。"""

import json
import logging

from fastapi import APIRouter, Depends, Request
from fastapi.responses import JSONResponse
from langchain_core.messages import HumanMessage, SystemMessage
from pydantic import BaseModel

from app.deps import get_optional_user_id
from app.services.llm_service import get_llm
from app.services.go_backend_client import get_problem_detail

logger = logging.getLogger(__name__)

router = APIRouter()


class GenerateTestsRequest(BaseModel):
    problem_id: int
    user_code: str = ""
    count: int = 5


@router.post("/generate-tests")
async def generate_tests(
    req: GenerateTestsRequest,
    request: Request,
    user_id: int = Depends(get_optional_user_id),
):
    """根据题目描述生成测试用例，返回 JSON 格式结果。"""

    # 从 request 提取 token 用于调用 Go 后端
    auth_header = request.headers.get("authorization", "")
    token = auth_header.replace("Bearer ", "") if auth_header.startswith("Bearer ") else auth_header

    # 获取题目详情
    problem = await get_problem_detail(req.problem_id, token)
    if not problem:
        return JSONResponse({"error": f"题目 {req.problem_id} 不存在或无法访问"}, status_code=404)

    title = problem.get("title", "")
    context = problem.get("context", "")
    input_desc = problem.get("input_description", "")
    output_desc = problem.get("output_description", "")
    samples = problem.get("judge_sample", [])
    tags = problem.get("tags", [])

    # 构建已有样本参考
    sample_text = ""
    if samples:
        sample_text = "\n已有示例：\n"
        for i, s in enumerate(samples[:3], 1):
            sample_text += f"  示例{i}: 输入={s.get('input', '')}, 输出={s.get('output', '')}\n"

    # 用户代码提示
    code_text = ""
    if req.user_code:
        code_text = (
            f"\n用户代码：\n```\n{req.user_code}\n```\n"
            "请特别关注代码中的潜在缺陷，生成能暴露这些问题的测试用例。"
        )

    system_prompt = """你是一位专业的算法竞赛测试工程师。根据题目描述生成全面覆盖的测试用例。

要求：
1. 必须生成涵盖以下类型的测试用例：
   - 边界值（最小/最大输入规模）
   - 特殊/退化情况（空输入、全相同、单元素）
   - 随机正常情况
   - 压力数据（接近约束上限）
2. 输出严格的 JSON 数组，每个元素包含 input、expected、explanation
3. 只输出 JSON，不要有任何其他文字"""

    user_msg = (
        f"题目：{title}\n"
        f"描述：{context}\n"
        f"输入说明：{input_desc}\n"
        f"输出说明：{output_desc}\n"
        f"标签：{', '.join(tags) if isinstance(tags, list) else tags}\n"
        f"{sample_text}{code_text}\n"
        f"生成 {req.count} 个测试用例。"
    )

    llm = get_llm()
    response = await llm.ainvoke([
        SystemMessage(content=system_prompt),
        HumanMessage(content=user_msg),
    ])

    # 解析 LLM 输出
    content = response.content.strip()
    if "```json" in content:
        content = content.split("```json")[1].split("```")[0].strip()
    elif "```" in content:
        content = content.split("```")[1].split("```")[0].strip()

    try:
        test_cases = json.loads(content)
        if isinstance(test_cases, list):
            return {"test_cases": test_cases}
    except json.JSONDecodeError:
        pass

    # JSON 解析失败，返回原始文本
    return {"test_cases": [], "raw_response": content}
