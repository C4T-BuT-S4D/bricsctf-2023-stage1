package ru.bricsctf.droidchat.data.api

import com.squareup.moshi.JsonClass
import com.squareup.moshi.Moshi
import com.squareup.moshi.adapter
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.withContext
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import ru.bricsctf.droidchat.di.IODispatcher
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ApiHelper
@Inject constructor(
    @IODispatcher private val ioDispatcher: CoroutineDispatcher,
    private val okHttpClient: OkHttpClient,
    private val moshi: Moshi
) {

    private suspend fun makeRequest(req: Request) =
        withContext(ioDispatcher) { okHttpClient.newCall(req).execute() }

    @OptIn(ExperimentalStdlibApi::class)
    private fun <T> errorOr(res: Response, block: (String) -> T) =
        if (res.isSuccessful) {
            Result.success(block.invoke(res.body!!.string()))
        } else {

            val jsonError = res.body!!.string()
            val error = moshi.adapter<ApiError>().fromJson(jsonError)!!
            Result.failure(Exception(error.error))
        }

    suspend fun <T> doCall(req: Request, resulter: (String) -> T) =
        errorOr(makeRequest(req), resulter)
}

@JsonClass(generateAdapter = true)
data class ApiError(
    val error: String
)

