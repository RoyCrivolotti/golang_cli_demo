# Notifier library

This simple notifier provides an instantiable class with a configurable URL that allows the consumer to send text
messages in the body of a `POST` request.

The library implements two version of the `Notify` method: one being async using `channels`, and one being `sync`.

Should an error be encountered in a `POST` request, the error being returned contains both the original error message
and the text message that failed to be notified.