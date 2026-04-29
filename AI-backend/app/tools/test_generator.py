"""LangChain Tool: 根据题目描述自动生成测试用例。"""

import json
import logging
import re

from langchain_core.tools import tool
from langchain_core.messages import HumanMessage, SystemMessage

from app.services.llm_service import get_llm
from app.services.go_backend_client import get_problem_detail

logger = logging.getLogger(__name__)

# 匹配说明性文字的特征词
_EXPLANATION_PATTERNS = re.compile(
    r"(应返回|测试|验证|检查|情况|处理|目的|用于|边界|压力|退化|正常|特殊|空输入)"
)


def normalize_test_cases(test_cases: list) -> list:
    """校验并修正 LLM 输出的测试用例格式。"""
    normalized = []
    for tc in test_cases:
        if not isinstance(tc, dict):
            continue

        raw_input = tc.get("input", "")
        raw_expected = tc.get("expected", "")
        raw_explanation = tc.get("explanation", "")

        # input: 数组 → 空格分隔字符串，其他类型 → str
        if isinstance(raw_input, list):
            raw_input = " ".join(str(x) for x in raw_input)
        else:
            raw_input = str(raw_input)

        # expected: 统一为字符串
        raw_expected = str(raw_expected)
        raw_explanation = str(raw_explanation)

        # 如果 expected 看起来是说明文字，且 explanation 为空，则互换
        if raw_expected and not raw_explanation and _EXPLANATION_PATTERNS.search(raw_expected):
            raw_explanation = raw_expected
            raw_expected = ""

        normalized.append({
            "input": raw_input,
            "expected": raw_expected,
            "explanation": raw_explanation,
        })

    return normalized


SYSTEM_PROMPT_TEST_CASES = """你是一位专业的算法竞赛测试工程师。根据题目描述和约束条件生成测试用例。

要求：
1. 生成涵盖以下类型的测试用例：
   - 边界值（最小/最大输入）
   - 特殊情况（空输入、全相同元素、退化情况）
   - 随机/正常情况
   - 压力测试（接近数据范围上限）
2. 输出严格的 JSON 数组格式

字段格式约束（必须严格遵守）：
- input：string 类型，为直接输入 stdin 的完整内容，多行用 \\n 分隔，同行用空格分隔。绝不能是数组或列表。
- expected：string 类型，为程序直接输出的内容（按输出格式）。必须是具体的数值或格式化结果，不能是说明文字。
- explanation：string 类型，为该测试用例的目的说明。

正确示例：
```json
[
  {"input": "5\\n1 2 3 4 5", "expected": "15", "explanation": "正常情况，验证基本功能"},
  {"input": "0", "expected": "0", "explanation": "空输入边界测试"}
]
```

错误示例（绝对禁止）：
```json
[
  {"input": [1, 2, 3], "expected": "应返回正确结果", "explanation": "正常测试"},
  {"input": "1 2 3", "expected": "验证排序功能", "explanation": ""}
]
```

只输出 JSON 数组，不要有其他文字。"""


@tool
async def generate_test_cases(
    problem_id: int,
    user_code: str = "",
    count: int = 5,
    token: str = "",
) -> str:
    """根据题目描述生成测试用例。当用户需要测试用例或想验证代码鲁棒性时使用。

    Args:
        problem_id: 题目 ID
        user_code: 用户提交的代码（可选，用于生成针对性测试）
        count: 生成测试用例数量，默认 5
        token: JWT token（由路由层注入，不由 LLM 生成）
    """
    # 从 Go 后端获取题目详情
    resp = await get_problem_detail(problem_id, token)
    if not resp:
        return f"无法获取题目 {problem_id} 的信息，请确认题目 ID 是否正确。"

    problem = resp.get("problem", resp)

    title = problem.get("title", "")
    context = problem.get("context", "")
    input_desc = problem.get("input_description", "")
    output_desc = problem.get("output_description", "")
    samples = problem.get("judge_sample", [])
    tags = problem.get("tags", [])

    # 构建样本参考
    sample_text = ""
    if samples:
        sample_text = "\n已有的示例测试用例：\n"
        for i, s in enumerate(samples[:3], 1):
            sample_text += f"  示例{i}: 输入={s.get('input', '')}, 输出={s.get('output', '')}\n"

    code_text = ""
    if user_code:
        code_text = f"\n用户的代码：\n```\n{user_code}\n```\n请特别关注代码中可能存在的边界问题，生成能暴露潜在 bug 的测试用例。"

    user_msg = f"""题目：{title}
题目描述：{context}
输入说明：{input_desc}
输出说明：{output_desc}
标签：{', '.join(tags) if isinstance(tags, list) else tags}
{sample_text}
{code_text}
请生成 {count} 个测试用例。注意：input 必须是 string，expected 必须是具体的程序输出。"""

    llm = get_llm()
    response = await llm.ainvoke([
        SystemMessage(content=SYSTEM_PROMPT_TEST_CASES),
        HumanMessage(content=user_msg),
    ])

    # 验证输出是有效 JSON
    content = response.content.strip()
    # 提取 JSON 块（LLM 可能包裹在 ```json ... ``` 中）
    if "```json" in content:
        content = content.split("```json")[1].split("```")[0].strip()
    elif "```" in content:
        content = content.split("```")[1].split("```")[0].strip()

    try:
        test_cases = json.loads(content)
        if isinstance(test_cases, list):
            test_cases = normalize_test_cases(test_cases)
            return json.dumps({"test_cases": test_cases}, ensure_ascii=False, indent=2)
    except json.JSONDecodeError:
        pass

    # JSON 解析失败，直接返回原始文本
    return content
