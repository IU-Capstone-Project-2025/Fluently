package ru.fluentlyapp.fluently.ui.utils

import androidx.compose.ui.tooling.preview.Preview

@Preview(
    name = "Large Phone (iPhone 16-ish)",
    device = "spec:width=430dp,height=932dp,dpi=460"
)
@Preview(
    name = "Medium Phone",
    device = "spec:width=393dp,height=851dp,dpi=440"
)
@Preview(
    name = "Small Phone",
    device = "spec:width=360dp,height=740dp,dpi=400"
)
@Preview(
    name = "Tablet Vertical",
    device = "spec:width=800dp,height=1280dp,dpi=320"
)
annotation class DevicePreviews