"""LangChain Tool: 从 Go 后端获取用户能力画像。"""

from langchain_core.tools import tool

from app.services.go_backend_client import get_user_ability


@tool
async def get_user_ability_profile(token: str = "") -> str:
    """获取当前用户的能力画像，包括各知识点掌握度、最强/最弱标签。当需要分析用户学习情况或提供个性化建议时使用。

    Args:
        token: JWT token（由路由层注入，不由 LLM 生成）
    """
    if not token:
        return "用户未登录，无法获取能力画像。"

    ability = await get_user_ability(token)
    if not ability:
        return "无法获取用户能力数据，可能是新用户暂无做题记录。"

    # 格式化为 LLM 易读的文本
    overall = ability.get("overall_score", 0)
    tag_scores = ability.get("tag_scores", {})
    strongest = ability.get("strongest_tags", [])
    weakest = ability.get("weakest_tags", [])
    languages = ability.get("languages", {})

    lines = [
        f"用户综合能力评分: {overall:.2f}（满分 1.0）",
        "",
        "各知识点掌握度:",
    ]

    for tag, score in sorted(tag_scores.items(), key=lambda x: x[1]):
        bar = "█" * int(score * 10) + "░" * (10 - int(score * 10))
        lines.append(f"  {tag}: [{bar}] {score:.2f}")

    if strongest:
        lines.append(f"\n最强领域: {', '.join(strongest)}")
    if weakest:
        lines.append(f"最弱领域: {', '.join(weakest)}")
    if languages:
        lang_list = [f"{lang}({cnt}次)" for lang, cnt in languages.items()]
        lines.append(f"常用语言: {', '.join(lang_list)}")

    return "\n".join(lines)
