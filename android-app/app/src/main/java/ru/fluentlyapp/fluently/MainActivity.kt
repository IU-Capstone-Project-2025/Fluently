package ru.fluentlyapp.fluently

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.ui.Modifier
import androidx.navigation.compose.NavHost
import dagger.hilt.android.AndroidEntryPoint
import ru.fluentlyapp.fluently.navigation.FluentlyNavHost

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            FluentlyNavHost(modifier = Modifier.fillMaxSize())
        }
    }
}
