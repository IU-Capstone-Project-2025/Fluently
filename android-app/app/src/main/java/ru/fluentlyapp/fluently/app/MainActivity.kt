package ru.fluentlyapp.fluently.app

import android.os.Build
import android.os.Bundle
import android.view.Window
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.asPaddingValues
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.systemBars
import androidx.compose.ui.Modifier
import dagger.hilt.android.AndroidEntryPoint
import ru.fluentlyapp.fluently.app.navigation.FluentlyNavHost
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            FluentlyTheme {
                FluentlyNavHost(
                    modifier = Modifier
                        .background(color = FluentlyTheme.colors.surface)
                        .fillMaxSize()
                )
            }
        }
    }
}
