package ru.bricsctf.droidchat.data.api

import ru.bricsctf.droidchat.data.model.Uid

//const val API_HOST = "https://c5ae-178-71-77-120.ngrok.io"
const val API_HOST = "https://droidchat-ab2f2aaa594034df.brics-ctf.ru"


class ApiRoutes {
    companion object {
        fun users(): String {
            return "$API_HOST/api/users/"
        }

        fun usersMe(): String {
            return "$API_HOST/api/users/me"
        }

        fun user(uid: Uid): String {
            return "$API_HOST/api/users/$uid"
        }

        fun usersToken(): String {
            return "$API_HOST/api/users/token"
        }

        fun myChats(): String {
            return "$API_HOST/api/chats/"
        }

        fun chatWith(uid: Uid): String {
            return "$API_HOST/api/chats/$uid"
        }

        fun stickers(): String {
            return "$API_HOST/api/stickers/"
        }
    }
}