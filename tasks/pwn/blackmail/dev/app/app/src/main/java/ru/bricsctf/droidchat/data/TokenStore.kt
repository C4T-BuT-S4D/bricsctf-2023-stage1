package ru.bricsctf.droidchat.data

import android.content.Context
import android.content.Context.MODE_PRIVATE
import androidx.core.content.edit
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject
import javax.inject.Singleton

const val SHAREDPREFS_KEY = "token"
const val SHAREDPREFS_NAME = "droidchat"

@Singleton
class TokenStore
@Inject constructor(
    @ApplicationContext context: Context
) {
    private val prefs = context.getSharedPreferences(SHAREDPREFS_NAME, MODE_PRIVATE)

    var token: String?
        get() = prefs.getString(SHAREDPREFS_KEY, null)
        set(value) = prefs.edit(commit = true) { putString(SHAREDPREFS_KEY, value) }
}