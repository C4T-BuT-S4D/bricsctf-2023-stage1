package ru.bricsctf.droidchat.data.model

import android.os.Bundle
import com.squareup.moshi.JsonClass


data class Chat(
    val with: User,
    val messages: List<Message>
)

@JsonClass(generateAdapter = false)
enum class MessageType {
    text, sticker
}

data class Message(
    val fromMe: Boolean,
    val type: MessageType,
    val text: String?,
    val sticker: Sticker?
)

@JsonClass(generateAdapter = true)
data class Sticker(
    val id: String,
    val url: String
) {
    companion object {
        fun createFromBundle(b: Bundle): Sticker {
            return Sticker(
                id = b.getString("_Sticker__id")!!,
                url = b.getString("_Sticker__url")!!
            )
        }
    }

    fun toBundle(): Bundle {
        return Bundle().apply {
            putString("_Sticker__id", id)
            putString("_Sticker__url", url)
        }
    }
}


