package ru.fluentlyapp.fluently.app

import android.app.Application
import coil3.ImageLoader
import coil3.PlatformContext
import coil3.SingletonImageLoader
import coil3.util.DebugLogger
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class FluentlyApplication : Application(), SingletonImageLoader.Factory {
    override fun newImageLoader(context: PlatformContext): ImageLoader {
        return ImageLoader.Builder(context)
            .logger(DebugLogger())
            .build()
    }
}