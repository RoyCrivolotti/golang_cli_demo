# Notifier library

This simple notifier provides an instantiable class with a configurable URL that allows the consumer to send text
messages in the body of a `POST` request. It uses `channels` to avoid doing blocking operations, meaning that the user
could send multiple requests at the same time, provided that they properly handle the `channels` being returned.

Should an error be encountered in a `POST` request, the error being returned contains both the original error message
and the text message that failed to be notified.