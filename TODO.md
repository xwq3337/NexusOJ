# 需求文档

## 提交记录 AI 分析

AI 分析结果进行 localstorage 存储，并记录分析时间
object {
    "id": "string", // 记录 id
    "analysis": "string", // 分析结果md
    "analysis_time": "string" // 分析时间
}

## AI 生成测试用例，仅发送题目 id 即可, 不需要发送用户代码 ,返回测试结果

object {
    "test_cases": [ // 测试用例列表
        {
            "input": "string", // 输入
            "expected": "string", // 期望输出
            "explanation": "string" // 实际输出
        }
    ],
    "analysis_time": "string" // 分析时间
}

## 发送当前页面的信息给 AI Assistant

例如:
    1. 题目代码编辑页: 发送题目id和用户代码，后端代理增加扩展信息(题目描述、测试用例、输入输出描述、用户代码等) {当用户询问题目时，可能会用到}
    2. 博客详情页面：发送博客 id，后端代理增加扩展信息(博客内容、博客标题、博客作者等)。{当用户让AI总结博客时，可能会用到}
    3. 等等
