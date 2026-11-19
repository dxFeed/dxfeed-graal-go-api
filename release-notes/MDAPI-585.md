* **\[MDAPI-585]\[GO]** Fixed a hang of many concurrent calls on Linux
  * An OS thread that calls the API is attached to the isolate on its first call and stays attached, instead of
    attaching and detaching on every call: with thousands of goroutines calling at once, attaching and detaching made
    the calls hang on Linux for minutes. A call on an attached thread no longer attaches or detaches it.
  * A goroutine locked to its OS thread with `runtime.LockOSThread` that calls the API should unlock before it exits
    (Threads in docs/troubleshooting.md).
