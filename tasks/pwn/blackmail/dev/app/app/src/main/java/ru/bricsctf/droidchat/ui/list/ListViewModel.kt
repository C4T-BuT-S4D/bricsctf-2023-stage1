package ru.bricsctf.droidchat.ui.list

import android.graphics.Bitmap
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch
import ru.bricsctf.droidchat.data.StickersRepository
import ru.bricsctf.droidchat.data.UsersRepository
import ru.bricsctf.droidchat.data.domain.GetChatsWithUsersUseCase
import ru.bricsctf.droidchat.data.model.ChatPreview
import ru.bricsctf.droidchat.data.model.User
import javax.inject.Inject

@HiltViewModel
class ListViewModel
@Inject constructor(
    private val usersRepository: UsersRepository,
    private val getChatsWithUsersUseCase: GetChatsWithUsersUseCase,
    private val stickersRepository: StickersRepository
)
: ViewModel() {
    private lateinit var userMe: User

    private val _chats = MutableLiveData<ListData>()
    val chats: LiveData<ListData> = _chats

    private val _users = MutableLiveData<List<User>>()
    val users: LiveData<List<User>> = _users

    private val _stickers = MutableLiveData<StickerPickerState>()
    val stickers: LiveData<StickerPickerState> = _stickers

    init {
        viewModelScope.launch {
            val result = usersRepository.getMe()
            userMe = result.getOrThrow()

            while (true) {
                delay(5000)
                refreshChats()
            }
        }
    }

    fun logout() {
        usersRepository.logout()
    }

    fun getUsers() {
        viewModelScope.launch {
            val usersResult = usersRepository.getUsers()
            _users.value = usersResult.getOrThrow()
        }
    }

    fun showStickerPicker(user: User) {
        viewModelScope.launch {
            val stickers = stickersRepository.getStickers().getOrThrow()
            _stickers.value = StickerPickerState(
                user = user,
                stickers = stickers.map {
                    PickableSticker(id = it.id, url = it.url,
                        stickerBitmap = stickersRepository.getStickerBitmap(it))
                }
            )
        }
    }

    private suspend fun refreshChats() {
        _chats.value = ListData(
            me = userMe,
            previews = getChatsWithUsersUseCase().getOrThrow()
        )
    }
}

data class ListData(
    val me: User,
    val previews: List<ChatPreview>
)

data class PickableSticker(
    val id: String,
    val url: String,
    val stickerBitmap: Bitmap
)

data class StickerPickerState(
    val user: User,
    val stickers: List<PickableSticker>
)
