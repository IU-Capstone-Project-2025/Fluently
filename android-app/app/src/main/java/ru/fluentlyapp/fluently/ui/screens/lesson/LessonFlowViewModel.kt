package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.common.model.Dialog
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.feature.dialog.DialogRepository
import ru.fluentlyapp.fluently.network.model.Author
import ru.fluentlyapp.fluently.network.model.Chat
import ru.fluentlyapp.fluently.network.model.Message
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.FinishDecorationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.LearningPartCompleteObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.OnboardingDecorationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.DialogObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordObserver
import ru.fluentlyapp.fluently.utils.safeLaunch
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class LessonFlowViewModel @Inject constructor(
    private val lessonRepository: LessonRepository,
    private val dialogRepository: DialogRepository
) : ViewModel() {
    // Start collecting the lesson state from the repository
    val currentComponent: StateFlow<LessonComponent?> = lessonRepository.currentComponent()
        .stateIn(
            viewModelScope,
            SharingStarted.Eagerly,
            initialValue = null
        )

    private val _commandsChannel = Channel<LessonFlowCommand>()
    val commandsChannel = _commandsChannel.receiveAsFlow()

    /**
     * Check if the current component is an instance of type `T`, and if it is, then try
     * to save the result of the block lambda into the `lessonRepository`.
     */
    private inline fun <reified T : LessonComponent> safeApplyAndUpdate(
        produceNewComponent: (oldComponent: T) -> T
    ) {
        (currentComponent.value as? T)?.let { oldComponent ->
            produceNewComponent(oldComponent).also { newComponent ->
                try {
                    viewModelScope.launch {
                        lessonRepository.updateCurrentComponent(newComponent)
                    }
                } catch (ex: Exception) {
                    Timber.e(ex)
                }
            }
        }
    }

    // Observer that is responsible for handling the "choose translation" exercises
    val chooseTranslationObserver = object : ChooseTranslationObserver() {
        override fun onVariantClick(variantIndex: Int) {
            safeApplyAndUpdate<Exercise.ChooseTranslation> {
                it.copy(selectedVariant = variantIndex)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for the "learn new word" exercises
    val newWordObserver = object : NewWordObserver() {
        override fun onUserKnowsWord(doesUserKnowWord: Boolean) {
            safeApplyAndUpdate<Exercise.NewWord> {
                it.copy(doesUserKnow = doesUserKnowWord)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for fill the gap exercises
    val fillGapsObserver = object : FillGapsObserver() {

        override fun onVariantClick(variantIndex: Int) {
            safeApplyAndUpdate<Exercise.FillTheGap> {
                it.copy(selectedVariant = variantIndex)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for input word exercise
    val inputWordObserver = object : InputWordObserver {
        override fun onConfirmInput(inputtedWord: String) {
            safeApplyAndUpdate<Exercise.InputWord> {
                it.copy(inputtedWord = inputtedWord.trim())
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    val onboardingDecorationObserver = object : OnboardingDecorationObserver() {
        override fun onContinue() {
            viewModelScope.safeLaunch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    val finishDecorationObserver = object : FinishDecorationObserver() {
        override fun onFinish() {
            viewModelScope.launch {
                try {
                    lessonRepository.dropLesson()
                    _commandsChannel.send(LessonFlowCommand.UserFinishesLesson)
                } catch (ex: Exception) {
                    Timber.e(ex)
                }
            }
        }
    }

    val learningPartCompleteObserver = object : LearningPartCompleteObserver() {
        override fun onMoveNext() {
            viewModelScope.safeLaunch {
                lessonRepository.sendCurrentProgress()
                lessonRepository.moveToNextComponent()
            }
        }
    }

    val dialogObserver = object : DialogObserver() {
        override fun onMoveNext() {
            viewModelScope.safeLaunch {
                lessonRepository.moveToNextComponent()
            }
        }

        override fun onCompleteDialog() {
            viewModelScope.launch {
                try {
                    dialogRepository.sendFinish()
                    safeApplyAndUpdate<Dialog> {
                        it.copy(isFinished = true)
                    }
                } catch (ex: Exception) {
                    Timber.e(ex)
                }
            }
        }

        private fun Dialog.toChat() = Chat(
            chat = messages.map {
                Message(
                    author = if (it.fromUser) Author.USER else Author.LLM,
                    message = it.text
                )
            }
        )

        private fun Chat.toDialog(oldDialog: Dialog?) = Dialog(
            messages = chat.withIndex().map { (index, message) ->
                Dialog.Message(
                    messageId = index.toLong(),
                    text = message.message,
                    fromUser = message.author == Author.USER
                )
            },
            isFinished = oldDialog?.isFinished == true,
            id = oldDialog?.id ?: -1
        )

        override fun onSendMessage(message: String) {
            val dialog: Dialog? = currentComponent.value as? Dialog
            if (dialog == null || dialog.isFinished) {
                return
            }
            viewModelScope.launch {
                try {
                    val currentChat = dialog.toChat().copy(
                        chat = dialog.toChat().chat + Message(
                            author = Author.USER,
                            message = message
                        )
                    )
                    val oldDialog = currentChat.toDialog(dialog)
                    safeApplyAndUpdate<Dialog> { oldDialog }
                    Timber.v(currentChat.toString())
                    val updatedChat = dialogRepository.sendChat(currentChat)
                    val updatedDialog = updatedChat.toDialog(oldDialog)
                    Timber.d("Updated dialog: $updatedDialog")
                    safeApplyAndUpdate<Dialog> { updatedDialog }
                } catch (ex: Exception) {
                    Timber.e(ex)
                    return@launch
                }
            }
        }
    }
}