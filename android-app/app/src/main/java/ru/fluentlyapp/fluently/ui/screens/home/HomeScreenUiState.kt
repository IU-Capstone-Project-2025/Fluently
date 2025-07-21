package ru.fluentlyapp.fluently.ui.screens.home

import android.net.Uri
import androidx.core.net.toUri
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState


data class HomeScreenUiState(
    val avatarPicture: Uri? = null,
    val wordOfTheDay: WordUiState? = null,
    val hasWordOfTheDaySaved: Boolean = false,
    val learnedWordsNumber: Int = 0,
    val inProgressWordsNumber: Int = 0,
    val preferredTopic: String? = null,
    val ongoingLessonState: OngoingLessonState = OngoingLessonState.NOT_STARTED
) {
    enum class OngoingLessonState {
        HAS_PAUSED,
        NOT_STARTED,
        ERROR,
        LOADING
    }
}