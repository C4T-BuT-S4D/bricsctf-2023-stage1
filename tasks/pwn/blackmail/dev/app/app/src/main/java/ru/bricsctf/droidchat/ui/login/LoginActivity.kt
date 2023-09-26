package ru.bricsctf.droidchat.ui.login

import android.content.Intent
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.View
import android.view.inputmethod.EditorInfo
import android.widget.EditText
import androidx.activity.viewModels
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.Observer
import com.google.android.material.dialog.MaterialAlertDialogBuilder
import dagger.hilt.android.AndroidEntryPoint
import ru.bricsctf.droidchat.databinding.ActivityLoginBinding

@AndroidEntryPoint
class LoginActivity : AppCompatActivity() {
    private lateinit var binding: ActivityLoginBinding
    private val loginViewModel: LoginViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityLoginBinding.inflate(layoutInflater)
        setContentView(binding.root)

        val username = binding.textUsername
        val usernameText = username.editText!!
        val password = binding.textPassword
        val passwordText = password.editText!!
        val login = binding.buttonLogin
        val loading = binding.loading


        setResult(RESULT_CANCELED)

        loginViewModel.loginFormState.observe(this@LoginActivity, Observer {
            val loginState = it ?: return@Observer

            // disable login button unless both username / password is valid
            login.isEnabled = loginState.isDataValid
            username.error = loginState.usernameError
            password.error = loginState.passwordError
        } )

        loginViewModel.loginResult.observe(this@LoginActivity, Observer {
            val result = it ?: return@Observer

            if (result.error == null) {
                loading.visibility = View.GONE

                setResult(RESULT_OK, Intent())
                finish()
            } else {
                MaterialAlertDialogBuilder(this)
                    .setTitle("Can't sign in")
                    .setMessage("Wrong username or password. You can try to sign up with these credentials.\n${result.error}")
                    .setPositiveButton("Try again") { _, _ ->
                        loading.visibility = View.GONE
                    }
                    .setNeutralButton("Sign up") { _, _ ->
                        loginViewModel.register(
                            usernameText.text.toString(),
                            passwordText.text.toString()
                        )
                    }
                    .show()
            }
        })

        usernameText.afterTextChanged {
            loginViewModel.loginDataChanged(
                usernameText.text.toString(),
                passwordText.text.toString()
            )
        }

        passwordText.apply {
            passwordText.afterTextChanged {
                loginViewModel.loginDataChanged(
                    usernameText.text.toString(),
                    passwordText.text.toString()
                )
            }


            setOnEditorActionListener { _, actionId, _ ->
                when (actionId) {
                    EditorInfo.IME_ACTION_DONE ->
                        loginViewModel.login(
                            usernameText.text.toString(),
                            passwordText.text.toString()
                        )
                }
                false
            }

            login.setOnClickListener {
                loading.visibility = View.VISIBLE
                loginViewModel.login(
                    usernameText.text.toString(),
                    passwordText.text.toString()
                )
            }
        }
    }
}

/**
 * Extension function to simplify setting an afterTextChanged action to EditText components.
 */
fun EditText.afterTextChanged(afterTextChanged: (String) -> Unit) {
    this.addTextChangedListener(object : TextWatcher {
        override fun afterTextChanged(editable: Editable?) {
            afterTextChanged.invoke(editable.toString())
        }

        override fun beforeTextChanged(s: CharSequence, start: Int, count: Int, after: Int) {}

        override fun onTextChanged(s: CharSequence, start: Int, before: Int, count: Int) {}
    })
}