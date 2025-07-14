package ru.fluentlyapp.fluently.ui.screens.wordsprogress

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.toRoute
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.app.navigation.Destination
import ru.fluentlyapp.fluently.app.navigation.Destination.WordsProgress.WordsProgressType
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgress
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepository
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import javax.inject.Inject

@HiltViewModel
class WordsInProgressViewModel @Inject constructor(
    private val savedStateHandle: SavedStateHandle,
    private val joinedWordProgressRepository: JoinedWordProgressRepository
) : ViewModel() {
    private val wordsProgressType: Destination.WordsProgress = savedStateHandle.toRoute()

    private val _uiState = MutableStateFlow<WordsProgressUiState>(
        WordsProgressUiState(
            words = emptyList(),
            pageTitle = when (wordsProgressType.progressType) {
                WordsProgressType.LEARNED -> "Выученные слова"
                WordsProgressType.IN_PROGRESS -> "Слова в обучении"
            },
            searchString = ""
        )
    )
    val uiState = _uiState.asStateFlow()

    private val searchString = MutableStateFlow<String>("")

    private var allWords: List<JoinedWordProgress> = emptyList()

    private fun matchesWordProgressType(word: JoinedWordProgress): Boolean {
        return (wordsProgressType.progressType == WordsProgressType.LEARNED && !word.isLearning) ||
                (wordsProgressType.progressType == WordsProgressType.IN_PROGRESS && word.isLearning)
    }

    init {
        viewModelScope.launch {
            allWords = joinedWordProgressRepository.getPerWordOverallProgress().first()
        }

        viewModelScope.launch {
            launch {
                searchString.collect { value ->
                    _uiState.update { oldState ->
                        oldState.copy(
                            searchString = value,
                            words = allWords.filter {
                                it.word.contains(value) && matchesWordProgressType(it)
                            }.map {
                                WordUiState(
                                    word = it.word,
                                    translation = it.translation,
                                    examples = it.examples
                                )
                            }
                        )
                    }
                }
            }
        }
    }
}