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
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgress
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepository
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import javax.inject.Inject

fun JoinedWordProgress.toWordUiState() = WordUiState(
    word = word,
    translation = translation,
    examples = examples
)

@HiltViewModel
class WordsInProgressViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val joinedWordProgressRepository: JoinedWordProgressRepository
) : ViewModel() {
    private val wordProgressRoute: Destination.WordsProgress = savedStateHandle.toRoute()

    private val _uiState = MutableStateFlow<WordsProgressUiState>(
        WordsProgressUiState(
            words = emptyList(),
            pageTitle = if (wordProgressRoute.isLearning) {
                "Слова в обучении"
            } else {
                "Выученные слова"
            },
            searchString = ""
        )
    )
    val uiState = _uiState.asStateFlow()

    private val searchString = MutableStateFlow<String>("")
    private var allSuitableWords: List<WordUiState> = emptyList()
    init {
        viewModelScope.launch {
            allSuitableWords = joinedWordProgressRepository
                .getPerWordOverallProgress()
                .first()
                .filter { it.isLearning == wordProgressRoute.isLearning }
                .map { it.toWordUiState() }

            _uiState.update { it.copy(words = allSuitableWords) }
        }

        viewModelScope.launch {
            searchString.collect { value ->
                _uiState.update { oldState ->
                    oldState.copy(searchString = value)
                }
            }
        }
    }

    fun updateSearchString(value: String) {
        searchString.value = value
    }

    fun onSearch() {
        _uiState.update {
            it.copy(
                words = allSuitableWords.filter { it.word.contains(searchString.value) }
            )
        }
    }
}