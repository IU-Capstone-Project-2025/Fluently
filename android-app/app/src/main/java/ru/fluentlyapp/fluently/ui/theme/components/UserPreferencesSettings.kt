package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.CefrLevel
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun UserPreferencesSettings(
    modifier: Modifier,
    userPreferences: UserPreferences,
    onPreferencesChange: (UserPreferences) -> Unit,
    availableTopics: List<String>
) {
    Column(modifier, horizontalAlignment = Alignment.CenterHorizontally) {
        Text(
            text = stringResource(R.string.approximate_level_of_english),
            fontWeight = FontWeight.Bold
        )
        Spacer(modifier = Modifier.height(4.dp))
        MaterialDropDownForm(
            value = userPreferences.cefrLevel.key,
            onValueChange = { key ->
                val newCefrLevel =
                    CefrLevel.entries.firstOrNull { it.key == key } ?: CefrLevel.A1
                onPreferencesChange(
                    userPreferences.copy(
                        cefrLevel = newCefrLevel
                    )
                )
            },
            values = CefrLevel.entries.map { it.toString() }
        )
        Spacer(modifier = Modifier.height(16.dp))
        Text(
            text = stringResource(R.string.topic_to_learn),
            fontWeight = FontWeight.Bold,
        )
        Spacer(modifier = Modifier.height(4.dp))
        MaterialDropDownForm(
            value = userPreferences.goal,
            onValueChange = { topic ->
                onPreferencesChange(
                    userPreferences.copy(
                        goal = topic
                    )
                )
            },
            values = availableTopics
        )
        Spacer(modifier = Modifier.height(16.dp))
        Text(
            text = stringResource(R.string.words_per_lesson),
            fontWeight = FontWeight.Bold
        )
        Spacer(modifier = Modifier.height(4.dp))
        MaterialDropDownForm(
            value = userPreferences.wordsPerDay.toString(),
            onValueChange = { newWordsPerDay ->
                onPreferencesChange(
                    userPreferences.copy(
                        wordsPerDay = newWordsPerDay.toIntOrNull() ?: 10
                    )
                )
            },
            values = (5..30 step 5).map { it.toString() }
        )
    }
}