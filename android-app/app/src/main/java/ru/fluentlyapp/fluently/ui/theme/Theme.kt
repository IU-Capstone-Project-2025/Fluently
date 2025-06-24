package ru.fluentlyapp.fluently.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.CompositionLocalProvider
import androidx.compose.runtime.staticCompositionLocalOf
import androidx.compose.ui.graphics.Color

data class FluentlyColors(
    val surface: Color = Color(0xFFF6F6F6),
    val surfaceContainerLow: Color = Color(0xFFEEF2F5),
    val surfaceContainerHigh: Color = Color(0xD0D0D3E7),

    val tertiary: Color = Color(0xFFA043DB),
    val tertiaryVariant1: Color = Color(0xFFCC8EF1),
    val tertiaryVariant2: Color = Color(0xFFF5E9FD),

    val secondary: Color = Color(0xFF119CE2),
    val secondaryVariant: Color = Color(0xFFDFF1FD),

    val primary: Color = Color(0xFFF08616),
    val primaryVariant: Color = Color(0xFFFFEFE2),

    val onSurface: Color = Color(0xFF1B1B1B),
    val onSurfaceVariant: Color = Color(0xFF707070),
    val onPrimary: Color = Color(0xFFFFFFFF),
    val onSecondary: Color = Color(0xFFFFFFFF),

    val error: Color = Color(0xFFb00020),

    // Additional local colors
    val googleBlue: Color = Color(0xFF4285F4)
)

val DefaultPalette = FluentlyColors()

val LocalFluentlyColors = staticCompositionLocalOf<FluentlyColors> { error("Colors are not provided") }

@Composable
fun ProvideFluentlyColors(fluentlyColors: FluentlyColors, content: @Composable () -> Unit) {
    CompositionLocalProvider(LocalFluentlyColors provides fluentlyColors) {
        content()
    }
}

object FluentlyTheme {
    val colors: FluentlyColors
        @Composable get() = LocalFluentlyColors.current
}

@Composable
fun FluentlyTheme(content: @Composable () -> Unit) {
    ProvideFluentlyColors(DefaultPalette) {
        MaterialTheme(
            content = content
        )
    }
}