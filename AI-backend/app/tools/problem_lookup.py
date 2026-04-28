"""LangChain Tool: 从 Go 后端获取题目详情。"""

import logging

from langchain_core.tools import tool

from app.services.go_backend_client import get_problem_detail

logger = logging.getLogger(__name__)


@tool
async def lookup_problem(problem_id: int, token: str = "") -> str:
    """获取题目详情，包括描述、输入输出说明、标签、难度。当需要了解某道题的具体内容时使用。

    Args:
        problem_id: 题目 ID
        token: JWT token（由路由层注入，不由 LLM 生成）
    """
    problem = await get_problem_detail(problem_id, token)
    if not problem:
        return f"未找到题目 {problem_id}，请确认题目 ID 是否正确。"

    title = problem.get("title", "")
    context = problem.get("context", "")
    input_desc = problem.get("input_description", "")
    output_desc = problem.get("output_description", "")
    difficulty = problem.get("difficulty", 0)
    tags = problem.get("tags", [])
    acceptance = problem.get("accept", 0)
    submissions = problem.get("submission", 0)

    rate = f"{acceptance / submissions * 100:.1f}%" if submissions > 0 else "N/A"

    lines = [
        f"题目 {problem_id}: {title}",
        f"难度: {difficulty:.1f} | 通过率: {rate} ({acceptance}/{submissions})",
        f"标签: {', '.join(tags) if isinstance(tags, list) else tags}",
        "",
        f"【题目描述】\n{context}",
        "",
        f"【输入说明】\n{input_desc}",
        "",
        f"【输出说明】\n{output_desc}",
    ]

    return "\n".join(lines)
