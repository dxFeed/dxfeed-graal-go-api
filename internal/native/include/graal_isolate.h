// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The functions of GraalVM that create, attach to and tear down an isolate, the instance of the SDK in the
 * process: graal_create_isolate(), graal_attach_thread(), graal_tear_down_isolate(), etc.
 *
 * native-image generates this header for the library of the SDK; the build of the SDK fails if the declarations
 * here differ from the ones that its GraalVM generates, so only the comments are the SDK's. The x64 libraries (macOS
 * x64 and the x86_64 slice of the iOS Simulator) are built with an older GraalVM, which declares version 4 of
 * graal_create_isolate_params_t: the same fields, named _reserved_1, _reserved_2 and _reserved_3 instead of argc,
 * argv and ignore_unrecognized_args; these libraries take the parameters of this header too.
 * graal_isolate_dynamic.h declares the same functions as pointer types, for loading the library at run time.
 *
 * An isolate has its own heap and its own Java state (the system properties, the endpoints, etc.): the handles of
 * one isolate cannot be used in another. Every thread that calls the SDK must be attached to the isolate, and every
 * function of the SDK takes the graal_isolatethread_t of the calling thread:
 *
 * @code{.c}
 * graal_isolate_t *isolate = NULL;
 * graal_isolatethread_t *thread = NULL;
 * if (graal_create_isolate(NULL, &isolate, &thread) != 0) { // the current thread is attached
 *     return; // the reason is printed to stderr
 * }
 *
 * // ... the SDK functions with the thread ...
 *
 * // Another thread:
 * graal_isolatethread_t *other = NULL;
 * if (graal_attach_thread(isolate, &other) == 0) {
 *     // ... the SDK functions with other ...
 *     graal_detach_thread(other); // before the thread exits
 * }
 *
 * // At the end (see dxfg_v4_system_shutdown()):
 * dxfg_v4_system_shutdown(thread);
 * // ... release the remaining handles ...
 * graal_tear_down_isolate(thread);
 * @endcode
 */

#ifndef __GRAAL_ISOLATE_H
#define __GRAAL_ISOLATE_H

struct __graal_isolate_t;
/**
 * @brief An isolate: the instance of the SDK (its heap and Java state) in the process.
 *
 * The pointer is valid until the isolate is torn down; graal_attach_thread() and graal_get_current_thread() take it.
 */
typedef struct __graal_isolate_t graal_isolate_t;

struct __graal_isolatethread_t;
/**
 * @brief The state of one thread attached to an isolate: the first parameter (`thread`) of the SDK functions.
 *
 * It belongs to the thread that attached itself (graal_create_isolate(), graal_attach_thread()) and must not be used
 * by other threads; it is valid until the thread detaches (graal_detach_thread()) or the isolate is torn down.
 */
typedef struct __graal_isolatethread_t graal_isolatethread_t;

#ifdef _WIN64
/** @brief An unsigned integer of the size of a pointer: the type of the sizes in graal_create_isolate_params_t. */
typedef unsigned long long __graal_uword;
#else
/** @brief An unsigned integer of the size of a pointer: the type of the sizes in graal_create_isolate_params_t. */
typedef unsigned long __graal_uword;
#endif

/**
 * @brief The values of graal_create_isolate_params_t.pkey (protection domains, internal: leave it 0).
 */
#define NO_PROTECTION_DOMAIN 0
/** @brief See #NO_PROTECTION_DOMAIN. */
#define NEW_PROTECTION_DOMAIN -1

/**
 * @brief The version of graal_create_isolate_params_t that this header declares (its `version` field).
 */
enum { __graal_create_isolate_params_version = 5 };

/**
 * @brief The parameters of graal_create_isolate().
 *
 * Zero the structure, then set `version` to #__graal_create_isolate_params_version and the fields to change; the
 * fields that are 0 keep their defaults. The isolate reads only the fields of the given version, so a program built
 * with an older header works with a newer library.
 *
 * The isolate arguments (`argc`, `argv`) are parsed like the command line of a native executable, and `argv[0]` (the
 * program name) is skipped:
 * - `-D<name>=<value>` sets a system property, as dxfg_system_set_property() does after the creation (for example the
 *   properties that must be set before the first endpoint is created);
 * - `-Xmx<size>`, `-Xms<size>`, `-Xmn<size>` set the maximum, the initial and the young generation sizes of the heap
 *   (`64m`, `1g`);
 * - `-XX:<option>=<value>`, `-XX:+<option>`, `-XX:-<option>` set the runtime options of the isolate, for example
 *   `-XX:TearDownWarningSeconds=<seconds>` (prints the threads that block graal_tear_down_isolate(); the debug builds
 *   of the SDK use 10 seconds) or `-XX:-EnableSignalHandling` (on Linux and macOS the isolate does not install its
 *   handlers of `SIGPIPE` and `SIGXFSZ`).
 *
 * An unknown `-XX:` option or a wrong value (`-Xmx64q`) fails graal_create_isolate() and prints the reason to stderr,
 * unless `ignore_unrecognized_args` is 1. The other arguments are not checked.
 *
 * @code{.c}
 * char *argv[] = {"app", "-Dlog.level=OFF", "-XX:TearDownWarningSeconds=10"};
 * graal_create_isolate_params_t params;
 * memset(&params, 0, sizeof(params));
 * params.version = __graal_create_isolate_params_version;
 * params.argc = 3;
 * params.argv = argv;
 * if (graal_create_isolate(&params, &isolate, &thread) != 0) { ... }
 * @endcode
 */
struct __graal_create_isolate_params_t {
    /** @brief The version of this structure: #__graal_create_isolate_params_version, set after zeroing it. */
    int version;

    /* Fields introduced in version 1 */
    /** @brief The size of the virtual address space to reserve for the heap (0: the default). */
    __graal_uword  reserved_address_space_size;

    /* Fields introduced in version 2. Internal usage, do not use. */
    /** @brief Internal: the path of an auxiliary image to load (leave it NULL). */
    const char    *auxiliary_image_path;
    /** @brief Internal: the bytes reserved for loading an auxiliary image (leave it 0). */
    __graal_uword  auxiliary_image_reserved_space_size;

    /* Fields introduced in version 3 */
    /** @brief The number of the strings in `argv`, including `argv[0]`. */
    int            argc;
    /** @brief The isolate arguments, parsed like a command line (see graal_create_isolate_params_t). */
    char         **argv;
    /** @brief Internal: the protection key or domain of the isolate (leave it 0, #NO_PROTECTION_DOMAIN). */
    int            pkey;

    /* Fields introduced in version 4 */
    /** @brief 1: unknown options in `argv` are ignored instead of failing graal_create_isolate(). */
    char           ignore_unrecognized_args;
    /** @brief Internal, leave it 0. */
    char           _reserved_4;

    /* Fields introduced in version 5 */
    /** @brief Internal, leave it 0. */
    char           _reserved_5;
};
/** @brief The parameters of graal_create_isolate(): see struct __graal_create_isolate_params_t. */
typedef struct __graal_create_isolate_params_t graal_create_isolate_params_t;

#if defined(__cplusplus)
extern "C" {
#endif

/**
 * @brief Creates a new isolate and attaches the current thread to it.
 *
 * @param[in] params The parameters (see graal_create_isolate_params_t), or NULL for the defaults.
 * @param[out] isolate Receives the isolate, if not NULL.
 * @param[out] thread Receives the isolate thread of the current thread, if not NULL.
 * @return 0 on success, a non-zero value on failure (an invalid isolate argument prints the reason to stderr).
 */
int graal_create_isolate(graal_create_isolate_params_t* params, graal_isolate_t** isolate,
                         graal_isolatethread_t** thread);

/**
 * @brief Attaches the current thread to the isolate.
 *
 * A thread must be attached before it calls the SDK, and it should detach (graal_detach_thread()) before it exits.
 * If the thread is already attached, the call succeeds and returns its isolate thread.
 *
 * @param[in] isolate The isolate.
 * @param[out] thread Receives the isolate thread of the current thread.
 * @return 0 on success, a non-zero value on failure.
 */
int graal_attach_thread(graal_isolate_t* isolate, graal_isolatethread_t** thread);

/**
 * @brief Returns the isolate thread of the current thread in the isolate.
 *
 * @param[in] isolate The isolate.
 * @return The isolate thread, or NULL if the current thread is not attached to the isolate or on another error.
 */
graal_isolatethread_t* graal_get_current_thread(graal_isolate_t* isolate);

/**
 * @brief Returns the isolate of an isolate thread.
 *
 * @param[in] thread The isolate thread.
 * @return The isolate, or NULL on error.
 */
graal_isolate_t* graal_get_isolate(graal_isolatethread_t* thread);

/**
 * @brief Detaches the isolate thread from its isolate and discards its state.
 *
 * No SDK function may be running on the isolate thread, and it must not be used afterwards.
 *
 * @param[in] thread The isolate thread of the current thread.
 * @return 0 on success, a non-zero value on failure.
 */
int graal_detach_thread(graal_isolatethread_t* thread);

/**
 * @brief Tears down the isolate of the isolate thread (which must still be attached): waits for the other attached
 * threads to detach, then discards the objects, the threads and the other state of the isolate.
 *
 * The call blocks while Java threads of the isolate run that do not stop on `Thread.interrupt()`. The threads of an
 * open endpoint do not, so close the endpoints first with dxfg_v4_system_shutdown() (see it for the whole shutdown
 * sequence). The isolate argument `-XX:TearDownWarningSeconds=<seconds>` prints the stack traces
 * of the threads that block the tear-down.
 *
 * @param[in] thread The isolate thread of the current thread.
 * @return 0 on success, a non-zero value on failure.
 */
int graal_tear_down_isolate(graal_isolatethread_t* thread);

/**
 * @brief Detaches all the threads that were started outside Java and attached to the isolate afterwards (including
 * the current one), shuts down the threads started in Java and tears down the isolate (see graal_tear_down_isolate()).
 *
 * None of the detached threads may be running Java code (an SDK function) when it is called or afterwards: that is
 * undefined and likely fatal behavior. It blocks as graal_tear_down_isolate() does.
 *
 * @param[in] thread The isolate thread of the current thread.
 * @return 0 on success, a non-zero value on a (non-fatal) failure.
 */
int graal_detach_all_threads_and_tear_down_isolate(graal_isolatethread_t* thread);

#if defined(__cplusplus)
}
#endif
#endif
