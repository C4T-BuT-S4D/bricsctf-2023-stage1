package ru.bricsctf.droidchat.ui.list

import android.content.Intent
import android.os.Bundle
import android.widget.ImageView
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.Observer
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import dagger.hilt.android.AndroidEntryPoint
import ru.bricsctf.droidchat.R
import ru.bricsctf.droidchat.data.model.ChatPreview
import ru.bricsctf.droidchat.data.model.Sticker
import ru.bricsctf.droidchat.data.model.User
import ru.bricsctf.droidchat.databinding.ActivityListBinding
import ru.bricsctf.droidchat.databinding.DialogStickerBinding
import ru.bricsctf.droidchat.ui.chat.ChatActivity
import ru.bricsctf.droidchat.ui.chat.ChatActivityArguments

@AndroidEntryPoint
class ListActivity : AppCompatActivity() {
    private lateinit var binding: ActivityListBinding
    private val viewModel: ListViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityListBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val title = binding.toolbar
        val list = binding.chatsList
        val logout = binding.toolbar.menu.findItem(R.id.logout)

        viewModel.chats.observe(this@ListActivity, Observer {
            val data = it ?: return@Observer

            // how slow is this to do each second?
//            Log.d("List", "updating")
            title.subtitle = data.me.username

            list.removeAllViews()
            data.previews.map {p ->
                list.addView(ChatListItemView(this).apply {
                    chatPreview = p
                })
            }
        })

        viewModel.users.observe(this@ListActivity, Observer {
            val data = it ?: return@Observer

            val usernames = data.map { u -> u.username }.toTypedArray()

            MaterialAlertDialogBuilder(this@ListActivity)
                .setTitle("Write to...")
                .setItems(usernames) { _, idx ->
                    openChat(data[idx])
                }
                .show()
        })

        viewModel.stickers.observe(this@ListActivity, Observer {
            val state = it ?: return@Observer

            val binding = DialogStickerBinding.inflate(layoutInflater)
            
            binding.stickers.removeAllViews()
            state.stickers.map { s ->
                binding.stickers.addView(ImageView(this@ListActivity).apply { 
                    setImageBitmap(s.stickerBitmap)
                    setOnClickListener {
                        openChatWithSticker(state.user, Sticker(id=s.id, url=s.url))
                    }
                })
            }

            MaterialAlertDialogBuilder(this@ListActivity)
                .setTitle("React...")
                .setView(binding.root)
                .show()
        })

        binding.newChat.setOnClickListener {
            viewModel.getUsers()
        }

        logout.setOnMenuItemClickListener {
            viewModel.logout()
            finish()
            return@setOnMenuItemClickListener false
        }
    }

    fun onStickerButton(chatPreview: ChatPreview) {
        viewModel.showStickerPicker(chatPreview.user)
    }

    private fun openChatWithSticker(user: User, pendingSticker: Sticker) {
        startActivity(Intent(this, ChatActivity::class.java).apply {
            putExtra("args", ChatActivityArguments(user = user, pendingSticker = pendingSticker).toBundle())
        })
    }

    fun openChat(user: User) {
        startActivity(Intent(this, ChatActivity::class.java).apply {
            putExtra("args", ChatActivityArguments(user = user, pendingSticker = null).toBundle())
        })
    }
}