package ru.bricsctf.droidchat.data.model

data class ChatPreview(
    val user: User,
    val isMine: Boolean,
    val isSticker: Boolean,
    val message: String?  // null if isSticker
)
