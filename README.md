# Salmon - Golabal Hub

Global Back-end sever for Salmon streaming platform providing account management and library synchronization.

This application manages user accounts and stores/synchronizes media library meta data from all instances of [salmon media servers](https://github.com/steams/salmon-media-server).
The application manages login sessions, utilizing http-only cookies and a custom auth header to protect against CSRF and XSS attacks. Sensitive data is encrypted using AES and passwords are salted and hashed using Bycrypt.

This application also serves an instance of the [salmon web player](https://github.com/Steams/salmon-web-client) for media playback.
