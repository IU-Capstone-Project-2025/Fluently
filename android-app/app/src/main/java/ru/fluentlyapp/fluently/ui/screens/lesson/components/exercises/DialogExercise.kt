package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.animation.core.tween
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.zIndex
import ru.fluentlyapp.fluently.common.model.Dialog
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.theme.components.ChatTextField
import ru.fluentlyapp.fluently.ui.theme.components.DialogTopFloatingButton
import ru.fluentlyapp.fluently.ui.theme.components.MessageBubble
import ru.fluentlyapp.fluently.ui.utils.MediumPhonePreview
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview
import kotlin.random.Random

abstract class DialogObserver {
    abstract fun onSendMessage(message: String)
    abstract fun onCompleteDialog()
    abstract fun onMoveNext()
}

@Composable
fun DialogExercise(
    modifier: Modifier = Modifier,
    exerciseState: Dialog,
    dialogObserver: DialogObserver,
    isCompleted: Boolean
) {
    Box(modifier = modifier.background(FluentlyTheme.colors.surface)) {
        DialogTopFloatingButton(
            modifier = Modifier
                .padding(top = 16.dp)
                .align(Alignment.TopCenter),
            text = if (isCompleted) {
                "Дальше"
            } else {
                "Закончить диалог"
            },
            onClick = {
                if (isCompleted) {
                    dialogObserver.onMoveNext()
                } else {
                    dialogObserver.onCompleteDialog()
                }
            }
        )

        val listState = rememberLazyListState()

        LaunchedEffect(exerciseState.messages.size) {
            if (exerciseState.messages.isNotEmpty()) {
                listState.animateScrollToItem(exerciseState.messages.size - 1)
            }
        }

        LazyColumn(
            state = listState,
            modifier = Modifier
                .fillMaxSize()
                .zIndex(-1f),
            verticalArrangement = Arrangement.Bottom,
            contentPadding = PaddingValues(bottom = 60.dp, start = 10.dp, end = 10.dp, top = 0.dp),
        ) {
            items(
                items = exerciseState.messages,
                key = { it.messageId }
            ) { message ->
                val alignment = if (message.fromUser) {
                    Alignment.TopEnd
                } else {
                    Alignment.TopStart
                }

                Column(
                    modifier = Modifier.animateItem(
                        fadeInSpec = tween()
                    )
                ) {
                    Spacer(modifier = Modifier.height(8.dp))
                    Box(
                        modifier = Modifier.fillMaxWidth(),
                        contentAlignment = alignment
                    ) {
                        // I don't give a shit why BoxWithConstraints complain on smth here
                        Box(
                            modifier = Modifier.fillMaxWidth(0.7f),
                            contentAlignment = alignment
                        ) {
                            MessageBubble(
                                text = message.text,
                                fromUser = message.fromUser
                            )
                        }
                    }
                }
            }
        }

        var currentText by remember { mutableStateOf("") }
        ChatTextField(
            modifier = Modifier
                .fillMaxWidth()
                .padding(8.dp)
                .heightIn(min = 60.dp)
                .align(Alignment.BottomCenter)
                .zIndex(0f),
            text = currentText,
            onTextChange = { currentText = it },
            onSendClick = { dialogObserver.onSendMessage(it) }
        )
    }
}

@SmallPhonePreview
@MediumPhonePreview
@Composable
fun DialogExercisePreview() {
    FluentlyTheme {
        val dialog = remember {
            mutableStateListOf<Dialog.Message>(
                Dialog.Message(messageId = -1, text = "fuck", fromUser = true)
            )
        }
        DialogExercise(
            modifier = Modifier
                .fillMaxSize(),
            exerciseState = Dialog(
                messages = dialog,
                isFinished = false
            ),
            dialogObserver = object : DialogObserver() {
                override fun onCompleteDialog() {}
                override fun onMoveNext() {}
                override fun onSendMessage(message: String) {
                    dialog.add(
                        Dialog.Message(
                            messageId = Random.nextLong(),
                            text = message,
                            fromUser = Random.nextBoolean()
                        )
                    )
                }
            },
            isCompleted = false
        )
    }
}