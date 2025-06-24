package ru.fluentlyapp.fluently.ui.screens.home

import android.net.Uri
import androidx.core.net.toUri

// TODO: remove hardcoded values from there
data class HomeScreenUiState(
    val goal: String = "Travelling",
    val avatarPicture: Uri = "https://previews.dropbox.com/p/thumb/ACqRahHLohfalD3qXFydxawY0bYoXYhM__WdNemYAjP8NlB_sUDVeD5wXkmzJ9rS4WEeRFLu13yfDC4YzSvk-mUCUDC60PPDJ9vjQys2-9J871mwKJRZBtMzgteC1O3cDOiDIVz1uQY4kGZr67ts2DPPR79VuTLNwFZhSxMBnYbDonh7LV4lYwCw8jLQPImHBn97YRW_xkZVLDRefKtgUf2EdqMSG1FAv7MU93cX7Bsg3HjrDhnQAKIhY5wwHP9J1xRPcfziJOfTruMO4OOLvHa43H6MDyjnwVvJDeS3MB6YHyX_kryCaw2rBFirFlCjKk3Y3JwQlAU_z8-f8FM_b0I0/p.jpeg?is_prewarmed=true".toUri(),
    val wordOfTheDay: String = "Car",
    val wordOfTheDayTranslation: String = "Машина",
    val notesNumber: Int = 0,
    val learnedWordsNumber: Int = 0,
    val inProgressWordsNumber: Int = 0,
    val hasOngoingLesson: Boolean = false,
)