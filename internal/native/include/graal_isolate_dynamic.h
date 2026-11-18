// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The functions of graal_isolate.h as pointer types (graal_create_isolate_fn_t, etc.), for a program that loads
 * the library of the SDK at run time (`dlopen` and `dlsym`, `LoadLibrary` and `GetProcAddress`) instead of linking it.
 *
 * native-image generates this header for the library of the SDK; the build of the SDK fails if the declarations here
 * differ from the ones that its GraalVM generates (see graal_isolate.h). It has the same include guard as
 * graal_isolate.h and the same types, so only the one included first counts: include this one first to get the
 * pointer types (dxfg_api.h, which includes graal_isolate.h, can still be included after it; its functions are
 * declared for linking, a program that loads them at run time declares their pointer types itself). See
 * graal_isolate.h for the functions and the isolate arguments.
 *
 * @code{.c}
 * #include <dlfcn.h>
 * #include "graal_isolate_dynamic.h"
 *
 * void *library = dlopen("libDxFeedGraalNativeSdk.so", RTLD_NOW); // .dylib on macOS
 * graal_create_isolate_fn_t create_isolate = (graal_create_isolate_fn_t) dlsym(library, "graal_create_isolate");
 * graal_tear_down_isolate_fn_t tear_down_isolate =
 *     (graal_tear_down_isolate_fn_t) dlsym(library, "graal_tear_down_isolate");
 *
 * graal_isolate_t *isolate = NULL;
 * graal_isolatethread_t *thread = NULL;
 * if (create_isolate(NULL, &isolate, &thread) == 0) {
 *     // ... the SDK functions, loaded the same way ...
 *     tear_down_isolate(thread);
 * }
 * @endcode
 *
 * On Windows: `LoadLibraryA("DxFeedGraalNativeSdk.dll")` and `GetProcAddress`.
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

/** @brief The type of graal_isolate.h's graal_create_isolate() (`dlsym(library, "graal_create_isolate")`). */
typedef int (*graal_create_isolate_fn_t)(graal_create_isolate_params_t* params, graal_isolate_t** isolate,
                                         graal_isolatethread_t** thread);

/** @brief The type of graal_isolate.h's graal_attach_thread() (`dlsym(library, "graal_attach_thread")`). */
typedef int (*graal_attach_thread_fn_t)(graal_isolate_t* isolate, graal_isolatethread_t** thread);

/** @brief The type of graal_isolate.h's graal_get_current_thread() (`dlsym(library, "graal_get_current_thread")`). */
typedef graal_isolatethread_t* (*graal_get_current_thread_fn_t)(graal_isolate_t* isolate);

/** @brief The type of graal_isolate.h's graal_get_isolate() (`dlsym(library, "graal_get_isolate")`). */
typedef graal_isolate_t* (*graal_get_isolate_fn_t)(graal_isolatethread_t* thread);

/** @brief The type of graal_isolate.h's graal_detach_thread() (`dlsym(library, "graal_detach_thread")`). */
typedef int (*graal_detach_thread_fn_t)(graal_isolatethread_t* thread);

/** @brief The type of graal_isolate.h's graal_tear_down_isolate() (`dlsym(library, "graal_tear_down_isolate")`). */
typedef int (*graal_tear_down_isolate_fn_t)(graal_isolatethread_t* thread);

/**
 * @brief The type of graal_isolate.h's graal_detach_all_threads_and_tear_down_isolate()
 * (`dlsym(library, "graal_detach_all_threads_and_tear_down_isolate")`).
 */
typedef int (*graal_detach_all_threads_and_tear_down_isolate_fn_t)(graal_isolatethread_t* thread);

#if defined(__cplusplus)
}
#endif
#endif
