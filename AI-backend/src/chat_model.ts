import { ChatZhipuAI } from "@langchain/community/chat_models/zhipuai";

let _model: ChatZhipuAI | null = null;

export function getModel() {
  if (!_model) {
    _model = new ChatZhipuAI({
      zhipuAIApiKey: process.env.GLM_API_KEY as string,
      model: process.env.GLM_MODEL_NAME as string,
      temperature: 0.5,
      topP: 0.9,
      streaming: true,
    });
  }
  return _model;
}
