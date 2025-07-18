package ru.fluentlyapp.fluently.ui.previewdata

import ru.fluentlyapp.fluently.ui.theme.components.WordUiState

val words = listOf(
    WordUiState(
        word = "summer",
        translation = "лето",
        examples = listOf(
            "I left Bobby Jr. at an amusement park. Told his father he was at summer camp." to
                    "Я Бобби как-то забыла в луна-парке, а отцу сказала, что он в лагере, выгадала неделю.",
            "In summer camp the kids used to pay me to make their beds." to
                    "В лагере дети платили мне, чтобы я убирала кровати.",
        )
    ),
    WordUiState(
        word = "car",
        translation = "машина",
        examples = listOf(
            "Mary can drive a car" to
                    "Мэри умеет управлять машиной",
        )
    )
)