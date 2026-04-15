import { MilvusClient } from "@zilliz/milvus2-sdk-node";
import { DynamicStructuredTool } from "@langchain/core/tools";
import { z } from "zod";

const COLLECTION_NAME = "algorithm_knowledge";

// 懒加载 Milvus 客户端
let milvusClient: MilvusClient | null = null;

function getMilvusClient(): MilvusClient {
  if (!milvusClient) {
    milvusClient = new MilvusClient({
      address: `${process.env.MILVUS_HOST}:${process.env.MILVUS_PORT}`,
      username: process.env.MILVUS_USER || "",
      password: process.env.MILVUS_PASSWORD || "",
    });
  }
  return milvusClient;
}

async function getEmbedding(text: string): Promise<number[] | null> {
  try {
    const res = await fetch("https://open.bigmodel.cn/api/paas/v4/embeddings", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${process.env.GLM_API_KEY}`,
      },
      body: JSON.stringify({ model: "embedding-2", input: text }),
    });
    const data = await res.json();
    return data.data?.[0]?.embedding ?? null;
  } catch (e) {
    console.error("获取embedding失败:", e);
    return null;
  }
}

async function search(query: string, topK = 3) {
  const client = getMilvusClient();
  const queryEmbedding = await getEmbedding(query);

  if (!queryEmbedding) {
    console.error("无法获取查询文本的向量");
    return [];
  }

  const results = await client.search({
    collection_name: COLLECTION_NAME,
    data: [queryEmbedding],
    anns_field: "embedding",
    params: { nprobe: 10 },
    limit: topK,
    output_fields: ["file_name", "content"],
  });

  return results.results.map((hit: any) => ({
    file_name: hit.file_name,
    content: hit.content,
    score: hit.score,
  }));
}

export const ragTool = new DynamicStructuredTool({
  name: "rag_search",
  description: "用于检索与用户查询相关的知识库文档，适用于回答关于算法、编程等技术问题",
  schema: z.object({
    query: z.string().describe("用户查询内容"),
  }),
  func: async ({ query }) => {
    try {
      const results = await search(query);

      if (!results.length) {
        return `未在知识库中找到与 '${query}' 相关的信息`;
      }

      let context = `根据知识库中找到的相关信息，关于 '${query}' 的内容如下：\n\n`;
      for (const result of results) {
        context += `来源文档: ${result.file_name}\n`;
        context += `内容摘要: ${result.content.slice(0, 500)}...\n`;
        context += `相似度得分: ${result.score.toFixed(4)}\n`;
        context += "-".repeat(80) + "\n";
      }

      return context;
    } catch (e: any) {
      return `检索过程中发生错误: ${e.message}`;
    }
  },
});
