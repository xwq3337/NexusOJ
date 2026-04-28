"""LangChain Tool: 代码分析（复杂度、风格、质量）。"""

from langchain_core.tools import tool

from app.services.llm_service import get_llm
from langchain_core.messages import HumanMessage, SystemMessage


@tool
async def analyze_code(code: str, language: str, analysis_type: str = "all") -> str:
    """分析用户提交的代码，提供复杂度评估、代码风格和质量建议。

    Args:
        code: 待分析的源代码
        language: 编程语言（如 cpp, python, java）
        analysis_type: 分析类型 - "complexity"（复杂度）, "style"（风格）, "quality"（质量）, "all"（全部）
    """
    type_descriptions = {
        "complexity": "时间复杂度和空间复杂度分析",
        "style": "代码风格评估（命名规范、缩进、注释）",
        "quality": "代码质量评估（潜在 bug、边界处理、可读性）",
        "all": "全面的代码分析，包括复杂度、风格和质量",
    }

    focus = type_descriptions.get(analysis_type, type_descriptions["all"])

    system_prompt = f"""你是一位资深代码审查专家。请对以下 {language} 代码进行{focus}。

请按以下结构输出：
1. **复杂度分析**：时间复杂度 O(?)、空间复杂度 O(?)，并解释原因
2. **代码风格**：命名规范、缩进格式、注释完整度（评分 1-10）
3. **质量评估**：潜在问题、边界条件处理、优化建议
4. **改进建议**：具体的改进方案和重构建议

请用中文回答，简洁专业。"""

    llm = get_llm()
    response = await llm.ainvoke([
        SystemMessage(content=system_prompt),
        HumanMessage(content=f"```{language}\n{code}\n```"),
    ])
    return response.content
