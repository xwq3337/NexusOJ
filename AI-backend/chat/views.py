import json
from pathlib import Path

from django.http import JsonResponse, StreamingHttpResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_POST

from chat.services import format_rag_context, stream_chat

_SYSTEM_PROMPT = (Path(__file__).resolve().parent.parent / "system_prompt.txt").read_text("utf-8")


def _sse_event(data):
    return f"data: {json.dumps(data, ensure_ascii=False)}\n\n"


@csrf_exempt
@require_POST
def chat_stream(request):
    try:
        body = json.loads(request.body)
    except (json.JSONDecodeError, ValueError):
        return JsonResponse({"error": "invalid JSON"}, status=400)

    messages = body.get("messages")
    if not messages or not isinstance(messages, list):
        return JsonResponse({"error": "messages is required"}, status=400)

    def event_generator():
        try:
            last_user_msg = None
            for m in reversed(messages):
                if m.get("role") == "user":
                    last_user_msg = m
                    break

            rag_context = ""
            if last_user_msg:
                yield _sse_event({"status": "searching"})
                rag_context = format_rag_context(last_user_msg.get("content", ""))

            if rag_context:
                system_content = (
                    f"{_SYSTEM_PROMPT}\n\n"
                    f"【知识库参考内容】\n{rag_context}\n\n"
                    f"请优先基于以上参考内容回答，如果参考内容中没有相关信息，可以结合自身知识回答。"
                )
            else:
                system_content = _SYSTEM_PROMPT

            api_messages = [{"role": "system", "content": system_content}]
            for m in messages:
                role = "assistant" if m.get("role") == "model" else m.get("role", "user")
                api_messages.append({"role": role, "content": m.get("content", "")})

            for text_chunk in stream_chat(api_messages):
                yield _sse_event({"text": text_chunk})

            yield _sse_event({"done": True})

        except Exception as e:
            yield _sse_event({"error": str(e)})

    response = StreamingHttpResponse(event_generator(), content_type="text/event-stream")
    response["Cache-Control"] = "no-cache"
    response["X-Accel-Buffering"] = "no"
    return response
