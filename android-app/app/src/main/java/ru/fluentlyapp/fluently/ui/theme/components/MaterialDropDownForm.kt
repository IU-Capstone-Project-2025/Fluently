package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.Text
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.*
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview

@Composable
fun MaterialDropDownForm(
    modifier: Modifier = Modifier,
    value: String?,
    onValueChange: (String) -> Unit,
    values: List<String>
) {
    var expanded by remember { mutableStateOf(false) }
    Box(
        modifier = modifier
            .clip(RoundedCornerShape(8.dp))
            .clickable(onClick = { expanded = true })
            .background(FluentlyTheme.colors.primaryVariant)
            .border(
                width = 2.dp,
                color = if (expanded) FluentlyTheme.colors.primary else Color.Unspecified,
                shape = RoundedCornerShape(8.dp)
            )
            .padding(horizontal = 16.dp, vertical = 8.dp)
    ) {
        Text(text = value ?: "")

        DropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false },
            containerColor = FluentlyTheme.colors.primaryVariant,
        ) {
            values.forEach { option ->
                DropdownMenuItem(
                    text = { Text(option) },
                    onClick = {
                        onValueChange(option)
                        expanded = false
                    }
                )
            }
        }
    }
}

@SmallPhonePreview
@Composable
fun DropdownPickerPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface),
            contentAlignment = Alignment.Center
        ) {
            var value = remember { mutableStateOf("") }

            MaterialDropDownForm(
                value = value.value,
                modifier = Modifier,
                onValueChange = { value.value = it },
                values = listOf("bear", "lemon", "abracadabra")
            )
        }
    }
}