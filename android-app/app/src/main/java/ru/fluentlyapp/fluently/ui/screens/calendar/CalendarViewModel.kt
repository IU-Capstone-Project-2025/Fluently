package ru.fluentlyapp.fluently.ui.screens.calendar

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgress
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepository
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import java.time.LocalDate
import java.time.ZoneId
import java.time.ZoneOffset
import javax.inject.Inject

fun JoinedWordProgress.toWordUiState() = WordUiState(
    word = word,
    translation = translation,
    examples = examples
)

@HiltViewModel
class CalendarViewModel @Inject constructor(
    private val joinedWordProgressRepository: JoinedWordProgressRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow<CalendarScreenUiState>(
        CalendarScreenUiState(
            day = LocalDate.now(),
            learnedWords = emptyList(),
            inProgressWords = emptyList(),
            showIsLearning = true
        )
    )
    val uiState = _uiState.asStateFlow()

    fun setDate(newDate: LocalDate) {
        viewModelScope.launch {
            joinedWordProgressRepository
                .getJoinedWordProgresses(
                    begin = newDate.atStartOfDay(ZoneId.systemDefault()).toInstant(),
                    end = newDate.plusDays(1).atStartOfDay(ZoneId.systemDefault()).toInstant()
                )
                .first()
                .also { progresses ->
                    _uiState.update {
                        it.copy(
                            day = newDate,
                            learnedWords = progresses.filter { !it.isLearning }
                                .map { it.toWordUiState() },
                            inProgressWords = progresses.filter { it.isLearning }
                                .map { it.toWordUiState() },
                        )
                    }
                }
        }
    }

    fun updateShowIsLearning(showIsLearning: Boolean) {
        _uiState.update {
            it.copy(
                showIsLearning = showIsLearning
            )
        }
    }

    init {
        setDate(LocalDate.now())
    }
}
