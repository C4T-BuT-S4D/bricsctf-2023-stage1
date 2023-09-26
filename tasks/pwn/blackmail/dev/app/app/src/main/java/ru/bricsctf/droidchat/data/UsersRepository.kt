package ru.bricsctf.droidchat.data

import com.squareup.moshi.JsonClass
import com.squareup.moshi.Moshi
import com.squareup.moshi.adapter
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.Request
import okhttp3.RequestBody.Companion.toRequestBody
import ru.bricsctf.droidchat.data.api.ApiHelper
import ru.bricsctf.droidchat.data.api.ApiRoutes
import ru.bricsctf.droidchat.data.model.Uid
import ru.bricsctf.droidchat.data.model.User
import javax.inject.Inject

@OptIn(ExperimentalStdlibApi::class)
class UsersRepository
@Inject constructor(
    private val apiHelper: ApiHelper,
    private val moshi: Moshi,
    private val tokenStore: TokenStore
) {
    private val JSON = "application/json; charset=utf-8".toMediaType()

    // Cache getMe() call!
    private lateinit var me: User

    suspend fun login(username: String, password: String) =
        apiHelper.doCall(
            Request.Builder()
                .url(ApiRoutes.usersToken())
                .post(
                    moshi.adapter<ApiCreds>()
                        .toJson(ApiCreds(username, password))!!
                        .toRequestBody(JSON)
                )
                .build()
        ) {
            val token = moshi.adapter<ApiToken>().fromJson(it)!!
            tokenStore.token = token.token
        }

    suspend fun register(username: String, password: String) =
        apiHelper.doCall(
            Request.Builder()
                .url(ApiRoutes.users())
                .post(
                    moshi.adapter<ApiCreds>()
                        .toJson(ApiCreds(username, password))!!
                        .toRequestBody(JSON)
                )
                .build()
        ) {
            moshi.adapter<ApiResultUser>().fromJson(it)!!.user
        }

    suspend fun getMe() =
        if (this::me.isInitialized)
            Result.success(me)
        else
            apiHelper.doCall(
                Request.Builder().url(ApiRoutes.usersMe()).build()
            ) {
                me = moshi.adapter<ApiResultUser>().fromJson(it)!!.user
                me
            }

    suspend fun getUser(id: Uid) =
        apiHelper.doCall(
            Request.Builder().url(ApiRoutes.user(id)).build()
        ) {
            moshi.adapter<ApiResultUser>().fromJson(it)!!.user
        }

    suspend fun getUsers() =
        apiHelper.doCall(
            Request.Builder().url(ApiRoutes.users()).build()
        ) {
            moshi.adapter<ApiResultUsers>().fromJson(it)!!.users
        }

    fun logout() {
        tokenStore.token = null
    }
}

@JsonClass(generateAdapter = true)
data class ApiResultUser(
    val user: User
)

@JsonClass(generateAdapter = true)
data class ApiResultUsers(
    val users: List<User>
)

@JsonClass(generateAdapter = true)
data class ApiToken(
    val token: String
)

@JsonClass(generateAdapter = true)
data class ApiCreds(
    val username: String,
    val password: String
)
