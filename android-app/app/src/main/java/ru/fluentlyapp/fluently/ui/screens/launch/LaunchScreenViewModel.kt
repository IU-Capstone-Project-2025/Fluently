package ru.fluentlyapp.fluently.ui.screens.launch

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.AuthRepository
import javax.inject.Inject

@HiltViewModel
class LaunchScreenViewModel @Inject constructor(
    val authRepository: AuthRepository
) : ViewModel() {
    val isUserLogged = MutableStateFlow<Boolean?>(null)

    init {
        viewModelScope.launch(Dispatchers.IO) {
            isUserLogged.update { authRepository.isUserLogged() }
        }
    }
}