package ru.bricsctf.droidchat.ui.chat

import android.graphics.Bitmap
import android.os.Bundle
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import ru.bricsctf.droidchat.data.ChatsRepository
import ru.bricsctf.droidchat.data.NewMessage
import ru.bricsctf.droidchat.data.StickersRepository
import ru.bricsctf.droidchat.data.UsersRepository
import ru.bricsctf.droidchat.data.domain.GetChatUseCase
import ru.bricsctf.droidchat.data.model.MessageType
import ru.bricsctf.droidchat.data.model.Uid
import ru.bricsctf.droidchat.data.model.User
import javax.inject.Inject

@HiltViewModel
class ChatViewModel
@Inject constructor(
    private val getChatUseCase: GetChatUseCase,
    private val chatsRepository: ChatsRepository,
    private val usersRepository: UsersRepository,
    private val stickersRepository: StickersRepository
) : ViewModel() {
    private val _state = MutableLiveData<ChatState>()
    val state: LiveData<ChatState> = _state

    private val _stickers = MutableLiveData<List<PickableSticker>>()
    val stickers: LiveData<List<PickableSticker>> = _stickers

    private var with: Uid = -1
    private lateinit var userWith: User

    fun tryOpenChat(intentArgs: Bundle) {
        val dlUid = intentArgs.getInt("dl.chats.uid", -1)

        viewModelScope.launch {
            if (getChatUseCase(dlUid).isSuccess) {
                with = dlUid
                userWith = usersRepository.getUser(dlUid).getOrThrow()
            } else {
                val argsExtra = ChatActivityArguments.createFromBundle(intentArgs.getBundle("args")!!)
                with = argsExtra.user.id
                userWith = argsExtra.user
                sendPendingSticker(argsExtra)
            }
            showStickerPicker()
            refreshChat()
        }
    }

    private suspend fun showStickerPicker() {
        val stickers = stickersRepository.getStickers().getOrThrow()
        _stickers.value = stickers.map {
            PickableSticker(id = it.id, stickerBitmap = stickersRepository.getStickerBitmap(it))
        }
    }

    fun sendTextMessage(text: String) {
        viewModelScope.launch {
            chatsRepository.sendMessage(with,
                NewMessage(
                    type = MessageType.text,
                    text = text,
                    stickerId = null
                )
            ).getOrThrow()
        }
    }

    fun sendSticker(id: String) {
        viewModelScope.launch {
            chatsRepository.sendMessage(with,
                NewMessage(
                    type = MessageType.sticker,
                    text = null,
                    stickerId = id
                )
            ).getOrThrow()
        }
    }

    private suspend fun refreshChat() {
        while (true) {
            delay(1000)
            _state.value = ChatState(
                with = userWith,
                messages = getChatUseCase(with).getOrThrow().messages.map {
                    DisplayedMessage(
                        fromMe = it.fromMe,
                        type = it.type,
                        text = it.text,
                        stickerBitmap = it.sticker?.let { it1 ->
                            stickersRepository.getStickerBitmap(it1)
                        }
                    )
                }
            )
        }
    }

    private suspend fun sendPendingSticker(args: ChatActivityArguments) {
        val sticker = args.pendingSticker ?: return

        _state.value = ChatState(
            with = args.user,
            messages = listOf(DisplayedMessage(
                fromMe = true,
                type = MessageType.sticker,
                stickerBitmap = stickersRepository.getStickerBitmap(sticker),
                text = null
            ))
        )

        chatsRepository.sendMessage(
            args.user.id,
            NewMessage(type=MessageType.sticker, stickerId = sticker.id, text = null)
        ).getOrThrow()
    }
}

data class ChatState(
    val with: User,
    val messages: List<DisplayedMessage>
)

data class DisplayedMessage(
    val fromMe: Boolean,
    val type: MessageType,
    val text: String?,
    val stickerBitmap: Bitmap?
)

data class PickableSticker(
    val id: String,
    val stickerBitmap: Bitmap
)
