package ru.bricsctf.droidchat.ui.chat

import android.os.Bundle
import ru.bricsctf.droidchat.data.model.Sticker
import ru.bricsctf.droidchat.data.model.User

data class ChatActivityArguments(
    val user: User,
    val pendingSticker: Sticker?
) {
    companion object {
        fun createFromBundle(b: Bundle): ChatActivityArguments =
            ChatActivityArguments(
                user = b.getBundle("_ChatActivityArguments__user")
                    !!.let{ User.createFromBundle(it) },
                pendingSticker = b.getBundle("_ChatActivityArguments__pendingSticker")
                    ?.let { Sticker.createFromBundle(it) }
            )
    }

    fun toBundle() = Bundle().apply {
        putBundle("_ChatActivityArguments__user", user.toBundle())
        pendingSticker?.let { putBundle("_ChatActivityArguments__pendingSticker", it.toBundle()) }
    }
}
