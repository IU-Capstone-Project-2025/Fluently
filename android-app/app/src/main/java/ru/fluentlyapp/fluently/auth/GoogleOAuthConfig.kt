package ru.fluentlyapp.fluently.auth

import net.openid.appauth.ResponseTypeValues

object GoogleOAuthConfig {
    const val AUTH_URI = "https://accounts.google.com/o/oauth2/auth"
    const val TOKEN_URI = "https://oauth2.googleapis.com/token"
    const val REDIRECT_URI = "ru.fluentlyapp.fluently:/"

    const val CLIENT_ID = "543284924233-vvqh7nov7srbubgalq5gpe0ng4cb5i8r.apps.googleusercontent.com"
    const val SCOPE = "https://www.googleapis.com/auth/userinfo.email"
    const val RESPONSE_TYPE = ResponseTypeValues.CODE
}