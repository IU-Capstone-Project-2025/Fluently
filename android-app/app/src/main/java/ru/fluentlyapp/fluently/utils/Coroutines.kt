package ru.fluentlyapp.fluently.utils

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import timber.log.Timber

fun CoroutineScope.safeLaunch(block: suspend CoroutineScope.() -> Unit) = launch {
    try {
        block()
    } catch (ex: Exception) {
        Timber.e(ex)
    }
}