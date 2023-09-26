package ru.bricsctf.droidchat.ui.list

import android.content.Context
import android.util.AttributeSet
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.constraintlayout.widget.ConstraintLayout
import ru.bricsctf.droidchat.data.model.ChatPreview
import ru.bricsctf.droidchat.databinding.ViewChatlistitemBinding

class ChatListItemView: ConstraintLayout {
    constructor(context: Context): super(context) {}
    constructor(context: Context, attributeSet: AttributeSet): super(context, attributeSet) {}

    private val _activity = context as ListActivity

    private val binding = ViewChatlistitemBinding.inflate(LayoutInflater.from(context), rootView as ViewGroup )
    private val username = binding.chatName
    private val content = binding.chatContent
    private val button = binding.stickerButton
    private val clickRoot = binding.clickRoot

    private var _chatPreview: ChatPreview? = null
    var chatPreview: ChatPreview
        get() = _chatPreview!!
        set(p) {
            _chatPreview = p

            username.text = p.user.username
            content.text = if (p.isSticker) "[Sticker]" else p.message!!

            clickRoot.setOnClickListener {
                _activity.openChat(p.user)
            }

            if (!p.isMine) {
                button.isEnabled = true
                button.visibility = View.VISIBLE
                button.setOnClickListener {
                    _activity.onStickerButton(chatPreview)
                }
            }
        }
}
