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
    val primaryVariant: Color = Color(0xFFE8CAB3),

    val onSurface: Color = Color(0xFF1B1B1B),
    val onSurfaceVariant: Color = Color(0xFF707070),
    val onPrimary: Color = Color(0xFFFFFFFF),
    val onSecondary: Color = Color(0xFFFFFFFF),

    val error: Color = Color(0xFFb00020),

    val surfaceInverse: Color = Color(0xFF090909),
    val onSurfaceInverse: Color = Color(0xFFe4e4e4),
    val onSurfaceVariantInverse: Color = Color(0xFF8f8f8f),

    // Additional local colors
    val googleBlue: Color = Color(0xFF4285F4),
    val correct: Color = Color(0xFF0CB20C),
    val wrong: Color = Color(0xFFCC1C1C)
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