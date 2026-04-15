import { getModel } from "../chat_model.js";
import { HumanMessage, AIMessage, SystemMessage } from "@langchain/core/messages";
import type { Request, Response } from "express";
import { ragTool } from "../tools/rag.js";
import { readFileSync } from "fs";
import { join } from "path";

interface Message {
  id: number; // 雪花算法生成的唯一ID
  role: "user" | "model"; // 消息角色
  content: string; // 消息内容
  timestamp: number; // 消息发送时间
}

interface ChatRequestBody {
  messages: Message[];
}

// 从项目根目录加载系统提示词
const SYSTEM_PROMPT = readFileSync(
  join(process.cwd(), "system_prompt.txt"),
  "utf-8",
);

export const ChatController = {
  ai_chat: async (req: Request<ChatRequestBody>, res: Response) => {
    console.log("Received request:", req.body);
    const { messages } = req.body;

    if (!messages || !Array.isArray(messages)) {
      res.status(400).json({ error: "messages is required" });
      return;
    }

    // 设置 SSE 响应头
    res.writeHead(200, {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
    });

    try {
      // 从最后一条用户消息中提取查询内容用于 RAG 检索
      const lastUserMsg = [...messages]
        .reverse()
        .find((m) => m.role === "user");
      let ragContext = "";

      if (lastUserMsg) {
        // 先发送检索状态
        res.write(`data: ${JSON.stringify({ status: "searching" })}\n\n`);

        const ragResult = await ragTool.invoke({ query: lastUserMsg.content });
        ragContext =
          typeof ragResult === "string" ? ragResult : String(ragResult);
      }

      // 构建 system prompt，注入 RAG 上下文
      const systemContent = ragContext
        ? `${SYSTEM_PROMPT}\n\n【知识库参考内容】\n${ragContext}\n\n请优先基于以上参考内容回答，如果参考内容中没有相关信息，可以结合自身知识回答。`
        : SYSTEM_PROMPT;

      // 转换为 LangChain 消息格式
      const lcMessages = [
        new SystemMessage(systemContent),
        ...messages.map((m) =>
          m.role === "model"
            ? new AIMessage(m.content)
            : new HumanMessage(m.content),
        ),
      ];

      const stream = await getModel().stream(lcMessages);
      for await (const chunk of stream) {
        const text =
          typeof chunk.content === "string"
            ? chunk.content
            : JSON.stringify(chunk.content);
        res.write(`data: ${JSON.stringify({ text })}\n\n`);
      }
      res.write(`data: ${JSON.stringify({ done: true })}\n\n`);
    } catch (err: any) {
      res.write(`data: ${JSON.stringify({ error: err.message })}\n\n`);
    }

    res.end();
  },
};
