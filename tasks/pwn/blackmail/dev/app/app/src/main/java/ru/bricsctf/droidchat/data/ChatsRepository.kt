package ru.bricsctf.droidchat.data

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass
import com.squareup.moshi.Moshi
import com.squareup.moshi.adapter
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import ru.bricsctf.droidchat.data.api.ApiHelper
import ru.bricsctf.droidchat.data.api.ApiRoutes
import ru.bricsctf.droidchat.data.model.MessageType
import ru.bricsctf.droidchat.data.model.Sticker
import ru.bricsctf.droidchat.data.model.Uid
import javax.inject.Inject

@OptIn(ExperimentalStdlibApi::class)
class ChatsRepository
@Inject constructor(
    private val apiHelper: ApiHelper,
    private val moshi: Moshi
) {
    private val JSON = "application/json; charset=utf-8".toMediaType()

    suspend fun getChat(with: Uid) =
        apiHelper.doCall(
            Request.Builder().url(ApiRoutes.chatWith(with)).build()
        ) {
            moshi.adapter<ApiResultChat>().fromJson(it)!!.chat
        }

    suspend fun getPreviews() =
        apiHelper.doCall(
            Request.Builder().url(ApiRoutes.myChats()).build()
        ) {
            moshi.adapter<ApiResultPreviews>().fromJson(it)!!.chats
        }

    suspend fun sendMessage(to: Uid, message: NewMessage) =
        apiHelper.doCall(
            Request.Builder()
                .url(ApiRoutes.chatWith(to))
                .post(
                    moshi.adapter<NewMessage>()
                        .toJson(message)!!
                        .toRequestBody(JSON)
                )
                .build()
        ) { }
}

@JsonClass(generateAdapter = true)
data class ApiResultPreviews(
    val chats: List<ApiChat>
)

@JsonClass(generateAdapter = true)
data class ApiResultChat(
    val chat: ApiChat
)

@JsonClass(generateAdapter = true)
data class ApiChat(
    val with: Uid,
    val messages: List<ApiMessage>
)

@JsonClass(generateAdapter = true)
data class ApiMessage(
    val from: Uid,
    val type: MessageType,
    val text: String?,
    val sticker: Sticker?
)

@JsonClass(generateAdapter = true)
data class NewMessage(
    val type: MessageType,
    val text: String?,
    @Json(name = "sticker_id") val stickerId: String?
)

