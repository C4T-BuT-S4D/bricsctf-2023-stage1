package ru.bricsctf.droidchat.data.domain

import ru.bricsctf.droidchat.data.ChatsRepository
import ru.bricsctf.droidchat.data.UsersRepository
import ru.bricsctf.droidchat.data.model.ChatPreview
import ru.bricsctf.droidchat.data.model.Uid
import ru.bricsctf.droidchat.data.model.User
import javax.inject.Inject

class GetChatsWithUsersUseCase
@Inject constructor(
    private val usersRepository: UsersRepository,
    private val chatsRepository: ChatsRepository
)
{
    private val usersCache = mutableMapOf<Uid, User>()

    suspend operator fun invoke(): Result<List<ChatPreview>> {
        val previews = chatsRepository.getPreviews()

        if (previews.isFailure)
            return Result.failure(previews.exceptionOrNull()!!)

        return Result.success(previews.getOrThrow().map {
            val msg = it.messages.first()

            ChatPreview(
                user = getUserCached(it.with),
                isMine = msg.from == usersRepository.getMe().getOrThrow().id,
                isSticker = msg.sticker != null,
                message = msg.text
            )
        })
    }

    private suspend fun getUserCached(uid: Uid): User {
        return usersCache[uid] ?:
            usersRepository.getUser(uid).getOrThrow().also { usersCache[uid] = it }
    }
}