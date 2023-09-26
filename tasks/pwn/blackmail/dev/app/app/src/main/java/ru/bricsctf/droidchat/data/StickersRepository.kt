package ru.bricsctf.droidchat.data

import android.graphics.Bitmap
import android.graphics.BitmapFactory
import com.squareup.moshi.JsonClass
import com.squareup.moshi.Moshi
import com.squareup.moshi.adapter
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.withContext
import okhttp3.OkHttpClient
import okhttp3.Request
import ru.bricsctf.droidchat.data.api.ApiHelper
import ru.bricsctf.droidchat.data.api.ApiRoutes
import ru.bricsctf.droidchat.data.model.Sticker
import ru.bricsctf.droidchat.di.IODispatcher
import javax.inject.Inject

@OptIn(ExperimentalStdlibApi::class)
class StickersRepository
@Inject constructor(
    @IODispatcher private val ioDispatcher: CoroutineDispatcher,
    private val okHttpClient: OkHttpClient,
    private val apiHelper: ApiHelper,
    private val moshi: Moshi
) {
    private val cache: MutableMap<String, Bitmap> = mutableMapOf()

    suspend fun getStickerBitmap(sticker: Sticker): Bitmap {
        if (!cache.containsKey(sticker.id)) {
            withContext(ioDispatcher) {
                val req = Request.Builder()
                    .url(sticker.url)
                    .build()
                val resp = okHttpClient.newCall(req).execute()
                val bitmap = BitmapFactory.decodeStream(resp.body!!.byteStream())
                cache[sticker.id] = bitmap
            }
        }
        return cache[sticker.id]!!
    }

    suspend fun getStickers() =
        apiHelper.doCall(Request.Builder()
            .url(ApiRoutes.stickers())
            .build()
        ) {
            moshi.adapter<ApiResultStickers>().fromJson(it)!!.stickers
        }
}

@JsonClass(generateAdapter = true)
data class ApiResultStickers(
    val stickers: List<Sticker>
)