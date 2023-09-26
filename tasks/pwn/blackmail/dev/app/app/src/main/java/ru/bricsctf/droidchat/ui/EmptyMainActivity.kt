package ru.bricsctf.droidchat.ui

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import dagger.hilt.android.AndroidEntryPoint
import ru.bricsctf.droidchat.R
import ru.bricsctf.droidchat.data.TokenStore
import ru.bricsctf.droidchat.ui.list.ListActivity
import ru.bricsctf.droidchat.ui.login.LoginActivity
import javax.inject.Inject

const val LOGIN_REQUEST = 0x1337C001

@AndroidEntryPoint
class EmptyMainActivity: AppCompatActivity() {
    @Inject lateinit var tokenStore: TokenStore

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_emptymain)

//        tokenStore.token = null
        if (tokenStore.token == null) {
            val intent = Intent(this, LoginActivity::class.java)
            startActivityForResult(intent, LOGIN_REQUEST)
        } else {
            val intent = Intent(this, ListActivity::class.java)
            startActivity(intent)
            finish()
        }
    }

    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)

        if (resultCode == RESULT_OK) {
            val intent = Intent(this, ListActivity::class.java)
            startActivity(intent)
            finish()
        } else if (resultCode == RESULT_CANCELED) {
            finish()
        }
    }
}