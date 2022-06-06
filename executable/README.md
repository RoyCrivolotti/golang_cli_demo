# Executable

This executable reads the input streamed by the user to `stdin` and, consuming the simple notifier library, sends the
messages (with each line being a message, regardless of its content or lack thereof) to the configured URL and prints
the messages that encountered an error to the user's screen.

The errors for the messages that fail to be processed are printed as the error occurs instead of after the program
finished, which gives the user some real-time feedback.
