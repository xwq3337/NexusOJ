from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
import uvicorn
import asyncio
import json
import os
from dotenv import load_dotenv
from langchain.agents import create_agent
from langchain_community.chat_models import ChatZhipuAI
from tools.rag_tool import rag_tool

app = FastAPI()


class QueryRequest(BaseModel):
    """请求体模型"""
    query: str


# 加载环境变量
load_dotenv(verbose=True)

# 初始化模型
model = ChatZhipuAI(
    model=os.getenv("AGENT_MODEL_NAME"),
    api_key=os.getenv("ZHIPU_API_KEY"),
    temperature=0.5,
)

# 加载系统提示词
with open("system_prompt.txt", 'r', encoding='utf-8') as f:
    system_prompt = f.read()

# 创建 agent
agent = create_agent(
    model,
    tools=[rag_tool],
    system_prompt=system_prompt
)


async def stream_generator(query: str):
    """流式输出生成器 - 使用 agent 并实现逐字流式输出"""
    try:
        # 使用 agent 获取响应
        for chunk in agent.stream({
            "messages": [
                {"role": "user", "content": query}
            ]
        }):
            # 从 chunk 中提取 AIMessage 的内容
            if chunk and 'model' in chunk:
                messages = chunk['model'].get('messages', [])
                for msg in messages:
                    if hasattr(msg, 'content'):
                        content = msg.content
                        # 将内容逐字发送，实现流式效果
                        for char in content:
                            yield f"data: {json.dumps({'content': char, 'done': False})}\n\n"
                            await asyncio.sleep(0.02)  # 控制流式速度

        # 发送结束标记
        yield f"data: {json.dumps({'content': '', 'done': True})}\n\n"

    except Exception as e:
        # 错误处理
        yield f"data: {json.dumps({'content': f'Error: {str(e)}', 'done': True})}\n\n"


@app.post("/test")
async def test_endpoint(request: QueryRequest):
    """
    POST /test 接口
    请求体: {"query": "hello"}
    返回: AI Agent 流式输出
    """
    return StreamingResponse(
        stream_generator(request.query),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no"  # 禁用 nginx 缓冲
        }
    )


@app.get("/")
async def root():
    return {"message": "Agent API Server", "status": "running"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5557)