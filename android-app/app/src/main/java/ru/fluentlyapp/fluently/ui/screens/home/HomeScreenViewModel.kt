package ru.fluentlyapp.fluently.ui.screens.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import javax.inject.Inject
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.supervisorScope
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepository
import ru.fluentlyapp.fluently.feature.wordoftheday.WordOfTheDayRepository
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreenUiState.OngoingLessonState
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import ru.fluentlyapp.fluently.utils.safeLaunch
import timber.log.Timber


@HiltViewModel
class HomeScreenViewModel @Inject constructor(
    private val lessonRepository: LessonRepository,
    private val joinedWordProgressRepository: JoinedWordProgressRepository,
    private val wordOfTheDayRepo: WordOfTheDayRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(HomeScreenUiState())
    val uiState = _uiState.asStateFlow()

    private val _commandsChannel = Channel<HomeCommands>()
    val commandsChannel = _commandsChannel.receiveAsFlow()

    init {
        viewModelScope.safeLaunch {
            supervisorScope {
                safeLaunch {
                    lessonRepository.currentComponent().collect { currentComponent ->
                        if (
                            uiState.value.ongoingLessonState in setOf(
                                OngoingLessonState.LOADING,
                                OngoingLessonState.ERROR
                            )
                        ) {
                            return@collect
                        }

                        if (currentComponent == null) {
                            _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.NOT_STARTED) }
                        } else {
                            _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
                        }
                    }
                }

                safeLaunch {
                    joinedWordProgressRepository.getPerWordOverallProgress().collect { progresses ->
                        _uiState.update {
                            it.copy(
                                learnedWordsNumber = progresses.count { !it.isLearning },
                                inProgressWordsNumber = progresses.count { it.isLearning }
                            )
                        }
                    }
                }

                safeLaunch {
                    wordOfTheDayRepo.updateWordOfTheDay()
                }

                safeLaunch {
                    wordOfTheDayRepo.isWordOfTheDayLearning().collect { isLearning ->
                        _uiState.update {
                            it.copy(
                                hasWordOfTheDaySaved = isLearning
                            )
                        }
                    }
                }

                safeLaunch {
                    wordOfTheDayRepo.getWordOfTheDay().collect { wordOfTheDay ->
                        if (wordOfTheDay == null) {
                            return@collect
                        }

                        _uiState.update {
                            it.copy(
                                wordOfTheDay = WordUiState(
                                    word = wordOfTheDay.word,
                                    translation = wordOfTheDay.translation,
                                    examples = wordOfTheDay.examples
                                )
                            )
                        }
                    }
                }
            }
        }
    }

    fun startLearningWordOfTheDay() {
        viewModelScope.safeLaunch {
            wordOfTheDayRepo.startLearningWordOfTheDay()
        }
    }

    fun ensureOngoingLesson() {
        _uiState.update {
            it.copy(ongoingLessonState = OngoingLessonState.LOADING)
        }

        viewModelScope.safeLaunch {
            try {
                if (lessonRepository.hasSavedLesson()) {
                    _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
                    _commandsChannel.send(HomeCommands.NavigateToLesson)
                    return@safeLaunch
                }

                lessonRepository.fetchAndSaveOngoingLesson()
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
                _commandsChannel.send(HomeCommands.NavigateToLesson)
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.ERROR) }
            }
        }
    }
}
