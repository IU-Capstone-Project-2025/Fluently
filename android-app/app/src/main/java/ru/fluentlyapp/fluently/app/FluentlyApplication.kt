package ru.fluentlyapp.fluently.app

import android.app.Application
import coil3.ImageLoader
import coil3.PlatformContext
import coil3.SingletonImageLoader
import coil3.util.DebugLogger
import dagger.hilt.android.HiltAndroidApp
import ru.fluentlyapp.fluently.BuildConfig
import timber.log.Timber

@HiltAndroidApp
class FluentlyApplication : Application(), SingletonImageLoader.Factory {
    // Set the logger for Coil library
    override fun newImageLoader(context: PlatformContext): ImageLoader {
        val imageLoader = ImageLoader.Builder(context)
        if (BuildConfig.DEBUG) {
            imageLoader.logger(DebugLogger())
        }
        return imageLoader.build()
    }

    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            // Enable logging only in debug builds
            Timber.plant(Timber.DebugTree())
        }
    }
}