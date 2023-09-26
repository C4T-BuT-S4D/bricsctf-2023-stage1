package ru.bricsctf.droidchat.data.model

import android.os.Bundle
import com.squareup.moshi.JsonClass

typealias Uid = Int

@JsonClass(generateAdapter = true)
data class User(
    val id: Uid,
    val username: String
) {
    companion object {
        fun createFromBundle(b: Bundle) =
            User(
                id = b.getInt("_User__id"),
                username = b.getString("_User__username")!!
            )
    }

    fun toBundle() = Bundle().apply {
        putInt("_User__id", id)
        putString("_User__username", username)
    }
}
