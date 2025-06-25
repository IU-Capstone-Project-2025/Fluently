package ru.fluentlyapp.fluently.ui.screens.home

import android.net.Uri
import androidx.core.net.toUri

// TODO: remove hardcoded values from there
data class HomeScreenUiState(
    val goal: String = "Travelling",
    val avatarPicture: Uri = "https://uc2bc8e6296b05c14f59323964db.previews.dropboxusercontent.com/p/thumb/ACoAPJoylOt475_bn7qs-IsD7p6fL2ZpS-_MNmDNbT7cCZQDjy7bY2bD6V6iPjkkljD4GOQK1zW8BSWS24y9BxOfuPS1VXSxhtCjd_gfx2-RnFAUWNuxt73LXfQs3VkZ_GVKPKqTuSinsGMvJbPa9_uQ9SOZ39SzmK21O7OKxZEF_CgCfwKk5ZeKvsVlzbzyACTGjSvU8kBWj89MSwRzATHsxAxT7vTVQhXt7oqaY7TGyc0REizYU0EjssXCUMqiBh44ffL1wnWLMnVUS46-FcL-Yfj35PoC-P_vR-hgJDFp2iFPvr2VHu5MXUCPsT2aGE9Haq-acYyAWyWoC61Bun1z6aUfUAQj4xFllFTQcJOVfMerU-zxfgBtrn8YQ9WCb7E/p.jpeg?is_prewarmed=true".toUri(),
    val wordOfTheDay: String = "Car",
    val wordOfTheDayTranslation: String = "Машина",
    val notesNumber: Int = 0,
    val learnedWordsNumber: Int = 0,
    val inProgressWordsNumber: Int = 0,
    val hasOngoingLesson: Boolean = false,
)