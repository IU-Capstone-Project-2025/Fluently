package ru.fluentlyapp.fluently.ui.screens.wordsprogress

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Clear
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.Icon
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.platform.LocalFocusManager
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

@Composable
fun SearchBar(
    modifier: Modifier,
    onQueryChange: (String) -> Unit,
    onSearch: () -> Unit,
    query: String
) {
    BasicTextField(
        value = query,
        singleLine = true,
        onValueChange = onQueryChange,
        textStyle = TextStyle(fontSize = 16.sp),
        keyboardOptions = KeyboardOptions.Default.copy(
            imeAction = ImeAction.Search
        ),
        keyboardActions = KeyboardActions(
            onSearch = {
                onSearch()
            }
        ),
        decorationBox = { innerTextField ->
            Row(
                modifier = modifier
                    .height(40.dp)
                    .clip(shape = RoundedCornerShape(8.dp))
                    .background(FluentlyTheme.colors.secondaryVariant)
                    .padding(horizontal = 8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(
                    Icons.Default.Search,
                    contentDescription = null,
                    tint = FluentlyTheme.colors.secondary,
                    modifier = Modifier.clip(CircleShape).clickable(onClick = onSearch)
                )
                Spacer(modifier = Modifier.width(4.dp))
                Box(
                    modifier = Modifier
                        .fillMaxHeight()
                        .weight(1f),
                    contentAlignment = Alignment.CenterStart
                ) {
                    innerTextField()
                }
                Icon(
                    Icons.Default.Clear,
                    contentDescription = null,
                    tint = FluentlyTheme.colors.secondary,
                    modifier = Modifier
                        .clip(CircleShape)
                        .clickable(
                            onClick = {
                                onQueryChange("")
                                onSearch()
                            }
                        )
                )
            }
        }
    )
}

@DevicePreviews
@Composable
fun SearchBarPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface),
            contentAlignment = Alignment.Center
        ) {
            var query by remember { mutableStateOf("") }
            SearchBar(
                modifier = Modifier
                    .fillMaxWidth(.8f)
                    .height(40.dp),
                onQueryChange = { query = it },
                onSearch = {
                    query = "SEARCHED!!!"
                },
                query = query
            )
        }
    }
}
