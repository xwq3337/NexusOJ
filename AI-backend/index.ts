import express from "express";
import * as dotenv from "dotenv";
import { ChatController } from "./src/controllers/chat.js";

dotenv.config({ path: [".env"] });

const app = express();
app.use(express.json());

app.post("/chat", ChatController.ai_chat);

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`AI Backend running on http://localhost:${PORT}`);
});
