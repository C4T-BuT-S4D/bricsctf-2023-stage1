package ru.bricsctf.droidchat.data.domain

import ru.bricsctf.droidchat.data.ChatsRepository
import ru.bricsctf.droidchat.data.UsersRepository
import ru.bricsctf.droidchat.data.model.Chat
import ru.bricsctf.droidchat.data.model.Message
import ru.bricsctf.droidchat.data.model.Uid
import ru.bricsctf.droidchat.data.model.User
import javax.inject.Inject

class GetChatUseCase
@Inject constructor(
    private val usersRepository: UsersRepository,
    private val chatsRepository: ChatsRepository
){

    private var userCache: User? = null

    suspend operator fun invoke(with: Uid): Result<Chat> {
        val me = usersRepository.getMe().getOrThrow()

        val apiChatRes = chatsRepository.getChat(with)
        if (apiChatRes.isFailure) {
            return Result.failure(apiChatRes.exceptionOrNull()!!)
        }
        val apiChat = apiChatRes.getOrThrow()

        if (userCache == null) {
            userCache = usersRepository.getUser(apiChat.with).getOrThrow()
        }

        val chat = Chat(
            with = userCache!!,
            messages = apiChat.messages.map {
                Message(
                    fromMe = it.from == me.id,
                    type = it.type,
                    text = it.text,
                    sticker = it.sticker
                )
            }
        )

        return Result.success(chat)
    }
}