package ru.fluentlyapp.fluently.network.model

data class Chat(
    val chat: List<Message>
)

enum class Author(val key: String) {
    LLM("llm"),
    USER("user")
}

data class Message(
    val author: Author,
    val message: String
)