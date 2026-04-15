from django.urls import path
from chat.views import chat_stream

urlpatterns = [
    path("chat", chat_stream, name="chat"),
]
