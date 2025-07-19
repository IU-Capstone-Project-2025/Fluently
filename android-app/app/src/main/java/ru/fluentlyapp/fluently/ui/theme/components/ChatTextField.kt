package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.key
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview

@Composable
fun ChatTextField(
    modifier: Modifier = Modifier,
    text: String,
    isEnabled: Boolean,
    onTextChange: (text: String) -> Unit,
    onSendClick: (text: String) -> Unit,
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.Bottom
    ) {
        Box(
            modifier = Modifier
                .weight(1f)
                .clip(RoundedCornerShape(16.dp))
                .background(color = FluentlyTheme.colors.surface)
                .border(
                    width = 2.dp,
                    color = FluentlyTheme.colors.primary,
                    shape = RoundedCornerShape(16.dp)
                )
                .padding(horizontal = 16.dp, vertical = 8.dp)
        ) {
            BasicTextField(
                modifier = Modifier
                    .heightIn(min = 40.dp)
                    .fillMaxWidth(),
                enabled = isEnabled,
                value = text,
                textStyle = TextStyle(
                    fontSize = 16.sp,
                    lineHeight = 20.sp
                ),
                keyboardOptions = KeyboardOptions(
                    imeAction = ImeAction.Send
                ),
                keyboardActions = KeyboardActions(
                    onSend = {
                        if (text.isNotEmpty()) {
                            onSendClick(text)
                            onTextChange("")
                        }
                    }
                ),
                onValueChange = onTextChange,
                decorationBox = { innerTextField ->
                    Box(contentAlignment = Alignment.CenterStart) {
                        if (text.isEmpty()) {
                            Text(
                                text = stringResource(R.string.Message),
                                fontSize = 16.sp,
                                lineHeight = 20.sp,
                                color = FluentlyTheme.colors.onSurfaceVariant
                            )
                        }
                        innerTextField()
                    }
                }
            )
        }
        Spacer(modifier = Modifier.width(8.dp))
        Box(
            modifier = Modifier
                .size(56.dp)
                .clip(CircleShape)
                .alpha(if (text.isEmpty() || !isEnabled) .5f else 1f)
                .background(FluentlyTheme.colors.primary)
                .clickable(
                    enabled = !text.isEmpty(),
                    onClick = {
                        onSendClick(text)
                        onTextChange("")
                    }
                ),
            contentAlignment = Alignment.Center
        ) {
            Icon(
                modifier = Modifier.size(40.dp),
                painter = painterResource(R.drawable.ic_arrow_up),
                tint = FluentlyTheme.colors.onPrimary,
                contentDescription = null
            )
        }
    }
}

@SmallPhonePreview
@Composable
fun ChatTextFieldPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface)
        ) {
            var text by remember { mutableStateOf("") }
            ChatTextField(
                modifier = Modifier
                    .align(Alignment.BottomCenter)
                    .padding(8.dp)
                    .fillMaxWidth(),
                text = text,
                onTextChange = { text = it },
                onSendClick = {},
                isEnabled = true
            )
        }
    }
}