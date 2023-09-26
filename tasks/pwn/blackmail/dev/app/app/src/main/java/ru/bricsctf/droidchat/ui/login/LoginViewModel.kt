package ru.bricsctf.droidchat.ui.login

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.launch
import ru.bricsctf.droidchat.data.UsersRepository
import javax.inject.Inject

@HiltViewModel
class LoginViewModel
@Inject constructor(
    private val usersRepository: UsersRepository
)
: ViewModel() {
    private val rgUsername = Regex("[a-zA-Z0-9]{4,32}")

    private val _loginForm = MutableLiveData<LoginFormState>()
    val loginFormState: LiveData<LoginFormState> = _loginForm

    private val _loginResult = MutableLiveData<LoginResult>()
    val loginResult: LiveData<LoginResult> = _loginResult

    fun login(username: String, password: String) {
        viewModelScope.launch {
            val result: Result<Unit> = usersRepository.login(username, password)

            if (result.isSuccess) {
                _loginResult.value = LoginResult(error = null)
            } else {
                _loginResult.value = LoginResult(error = result.exceptionOrNull()!!.toString())
            }
        }
//        val result = usersRepository.login(username, password)
    }

    fun register(username: String, password: String) {
        viewModelScope.launch {
            val result = usersRepository.register(username, password)

            if (result.isSuccess) {
                login(username, password)
            } else {
                _loginResult.value = LoginResult(error = result.exceptionOrNull()!!.toString())
            }
        }
    }

    fun loginDataChanged(username: String, password: String) {
        if (username.isNotEmpty() && !isUsernameValid(username)) {
            _loginForm.value = LoginFormState(usernameError = "Not a valid username")
        } else if (password.isNotEmpty() && !isPasswordValid(password)) {
            _loginForm.value = LoginFormState(passwordError = "Uhh")
        } else if (username.isEmpty() || password.isEmpty()) {
            _loginForm.value = LoginFormState(isDataValid = false)
        } else {
            _loginForm.value = LoginFormState(isDataValid = true)
        }
    }

    private fun isUsernameValid(username: String) = rgUsername.matches(username)

    private fun isPasswordValid(password: String) = password.length in 8..72

}

data class LoginFormState(
    val usernameError: String? = null,
    val passwordError: String? = null,
    val isDataValid: Boolean = false
)

data class LoginResult(
    val error: String? = null
)
