import { GoogleGenAI } from '@google/genai'

// NOTE: In production you should proxy requests to keep API keys secret.
const apiKey = (import.meta.env.VITE_GEMINI_API_KEY as string) || ''
const ai = new GoogleGenAI({ apiKey })

export const getGeminiResponse = async (prompt: string, context?: string): Promise<string> => {
  if (!apiKey) {
    return 'Configuration Error: API Key is missing. Please set process.env.API_KEY.'
  }

  try {
    const modelId = 'gemini-2.5-flash'
    const systemInstruction = `You are NexusAI, an expert algorithm instructor and coding assistant for NexusOJ. 
    Your goal is to help users understand algorithms, debug code, and improve their competitive programming skills.
    
    Context:
    ${context || 'The user is browsing the NexusOJ platform.'}
    
    Guidelines:
    - Be concise and encouraging.
    - If the user asks for a solution, guide them with hints first before providing code.
    - Support languages: Python, C++, Java, Rust, Go, JavaScript, C.
    - You can explain time and space complexity.
    `

    const response = await ai.models.generateContent({
      model: modelId,
      contents: prompt,
      config: {
        systemInstruction
      }
    })

    return response.text || "I couldn't generate a response. Please try again."
  } catch (error) {
    console.error('Gemini API Error:', error)
    return 'Sorry, I encountered an error while processing your request.'
  }
}
