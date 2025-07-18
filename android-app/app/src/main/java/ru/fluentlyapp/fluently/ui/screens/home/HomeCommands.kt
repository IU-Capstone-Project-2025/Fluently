package ru.fluentlyapp.fluently.ui.screens.home

sealed interface HomeCommands {
    object NavigateToLesson : HomeCommands
}