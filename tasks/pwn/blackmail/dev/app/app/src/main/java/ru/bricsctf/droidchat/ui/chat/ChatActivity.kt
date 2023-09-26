package ru.bricsctf.droidchat.ui.chat

import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.Gravity
import android.view.View
import android.widget.EditText
import android.widget.ImageView
import android.widget.LinearLayout
import android.widget.TextView
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.Observer
import dagger.hilt.android.AndroidEntryPoint
import ru.bricsctf.droidchat.R
import ru.bricsctf.droidchat.data.model.MessageType
import ru.bricsctf.droidchat.databinding.ActivityChatBinding

@AndroidEntryPoint
class ChatActivity : AppCompatActivity() {
    private lateinit var binding: ActivityChatBinding
    private val viewModel: ChatViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityChatBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val messagesList = binding.messagesList
        val message = binding.message
        val send = binding.sendButton

        var cachedLastMessage: DisplayedMessage? = null

        viewModel.state.observe(this@ChatActivity, Observer {
            val chat = it ?: return@Observer

            binding.toolbar.title = chat.with.username

            if (chat.messages.isEmpty()) return@Observer

            messagesList.removeAllViews()
            val last = chat.messages.map(this::messageToView).map {v ->
                messagesList.addView(v)
                v
            }.lastOrNull()
            if (chat.messages.last() != cachedLastMessage) {
                last?.focusable = View.FOCUSABLE
                last?.isFocusableInTouchMode = true
                last?.requestFocus()
                cachedLastMessage = chat.messages.last()
            }
        })

        viewModel.stickers.observe(this@ChatActivity, Observer {
            val stickers = it ?: return@Observer

            binding.stickers.removeAllViews()
            stickers.map {s ->
                val imageView = ImageView(this@ChatActivity).apply {
                    setImageBitmap(s.stickerBitmap)
                    setOnClickListener {
                        viewModel.sendSticker(s.id)
                    }
                }
                binding.stickers.addView(imageView)
            }
        })

        message.editText!!.afterTextChanged {
            send.isEnabled = it.isNotBlank()
        }

        send.setOnClickListener {
            viewModel.sendTextMessage(message.editText!!.text.toString())
            message.editText!!.text.clear()
        }

        binding.toolbar.setNavigationOnClickListener { finish() }

        // init viewModel
        viewModel.tryOpenChat(intent.extras!!)
    }

    private fun messageToView(message: DisplayedMessage): View {
        val view: View

        if (message.type == MessageType.sticker) {
            val imageView = ImageView(this@ChatActivity).apply {
                setImageBitmap(message.stickerBitmap) }
            view = imageView
        } else {
            val textView = TextView(this@ChatActivity).apply {
                text = message.text
            }
            view = textView
        }

        view.setPadding(8, 8, 8, 8)

        if (message.fromMe) {
            view.setBackgroundColor(resources.getColor(R.color.purple_200))
        }

        view.layoutParams = LinearLayout.LayoutParams(
            LinearLayout.LayoutParams.WRAP_CONTENT,
            LinearLayout.LayoutParams.WRAP_CONTENT
        ).apply {
            weight = 1.0f
            gravity = if (message.fromMe) Gravity.RIGHT else Gravity.LEFT
        }

        return view
    }
}

/**
 * Extension function to simplify setting an afterTextChanged action to EditText components.
 */
fun EditText.afterTextChanged(afterTextChanged: (String) -> Unit) {
    this.addTextChangedListener(object : TextWatcher {
        override fun afterTextChanged(editable: Editable?) {
            afterTextChanged.invoke(editable.toString())
        }

        override fun beforeTextChanged(s: CharSequence, start: Int, count: Int, after: Int) {}

        override fun onTextChanged(s: CharSequence, start: Int, before: Int, count: Int) {}
    })
}
