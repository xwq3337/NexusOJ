export interface ChatMessage {
  id?: string
  sender_id : string
  receiver_id : string
  status: boolean
  message_type: 'text' | 'image' | 'voice' | 'video' | 'file'
  content: string
  created_at: string // "2025-06-13T04:31:21.728Z"
}
