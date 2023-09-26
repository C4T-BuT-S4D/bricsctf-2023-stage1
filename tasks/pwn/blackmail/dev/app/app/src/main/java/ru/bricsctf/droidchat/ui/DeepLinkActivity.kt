package ru.bricsctf.droidchat.ui

import android.net.Uri
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import ru.bricsctf.droidchat.ui.chat.ChatActivity
import ru.bricsctf.droidchat.ui.list.ListActivity

class DeepLinkActivity: AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val intent = intent

        val _data = intent.data ?: finish()
        val data: Uri = _data as Uri

        if (data.scheme != "droidchat") {
            finish()
        }

        when (data.host) {
            "list" -> {
                intent.setClass(this, ListActivity::class.java)
            }
            "chat" -> {
                intent.setClass(this, ChatActivity::class.java)
                if (data.pathSegments.size != 1) finish()
                intent.putExtra("dl.chats.uid", data.pathSegments.first().toInt())
            }
            else -> {
                finish()
            }
        }

        startActivity(intent)
        finish()
    }
}