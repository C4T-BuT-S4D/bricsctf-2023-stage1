package ru.bricsctf.droidchat.data.api

import okhttp3.Interceptor
import okhttp3.Response
import ru.bricsctf.droidchat.data.TokenStore
import javax.inject.Inject

class AuthInterceptor
@Inject constructor(
    private val tokenStore: TokenStore
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        var req = chain.request()

        // vuln!
        if (req.url.toString().startsWith(API_HOST)) {
            val newUrl = req.url.newBuilder()
                // or idk use Authorization header
                .addQueryParameter("token", tokenStore.token ?: "")
                .build()
            req = req.newBuilder().url(newUrl).build()
        }

        return chain.proceed(req)
    }
}