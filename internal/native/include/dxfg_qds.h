// Copyright (c) 2025 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief dxFeed Graal Native SDK QDS functions declarations
 */

#ifndef DXFG_QDS_H
#define DXFG_QDS_H

#ifdef __cplusplus
extern "C" {
#    include <cstdint>
#else
#    include <stdint.h>
#endif

#include "dxfg_common.h"

#include "graal_isolate.h"
#include "dxfg_javac.h"

/** @defgroup Qds Qds
 *  @brief The functions of the QDS Tools of the Java library: the parser of symbol lists and the command line of the
 *  tools.
 *
 *  A list that a function returns is owned by the caller: release it with dxfg_CList_String_release() (see the group
 *  Javac for the ownership rules). The exception of a failed function is available from
 *  dxfg_get_and_clear_thread_exception_t().
 *  @{
 */

/**
 * @brief Parses a list of symbols (a symbol set filter specification) into the symbols that it contains.
 *
 * The specification is parsed as a QD symbol set filter of the default data scheme (see the QDS documentation for the
 * syntax of the filters); a pattern that contains the ".ipf" substring is recognized as the URL of an IPF file that
 * gives the set of symbols. The symbols of the set are decoded to strings and returned in the order in which the set
 * gives them, which is not specified. The call fails if the specification is invalid.
 *
 * Java: `Tools.parseSymbols(String, DataScheme)`, with the default scheme (`QDFactory.getDefaultScheme()`).
 *
 * @deprecated Use dxfg_v4_Tools_parseSymbols().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbolList The specification of the symbol list: a UTF-8 null-terminated string; the caller keeps its
 * ownership.
 * @return The list of the symbols (release it with dxfg_CList_String_release()), or NULL on error. Use
 * dxfg_get_and_clear_thread_exception_t() to get the exception.
 */
DXFG_DEPRECATED("use dxfg_v4_Tools_parseSymbols()")
dxfg_string_list*   dxfg_Tools_parseSymbols(graal_isolatethread_t *thread, const char* symbolList);

/**
 * @brief Runs the QDS tools with the command line arguments.
 *
 * The first argument is the name of the tool, the others are its arguments; several tools can be run together,
 * separated by the argument "+". The function blocks the calling thread until the tools finish. With no arguments
 * (an empty list or NULL) it prints the usage to the standard error stream and returns #DXFG_EXECUTE_SUCCESSFULLY.
 *
 * Warning: with at least one argument the Java `Tools.main()` ends with `System.exit()`: the whole process, not only
 * the isolate, is terminated with the exit code 0 if the tools succeeded and 1 if they failed, so the function does
 * not return to the caller. Use it only from a process that is a launcher of a tool.
 *
 * Java: `Tools.main(String[])`.
 *
 * @deprecated Use dxfg_v4_Tools_main().
 *
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] args The list of the arguments; the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if no arguments are passed (the process is terminated otherwise, see above),
 * or #DXFG_EXECUTE_FAIL (-1) on error. Use dxfg_get_and_clear_thread_exception_t() to determine if an exception was
 * thrown.
 */
DXFG_DEPRECATED("use dxfg_v4_Tools_main()")
int32_t             dxfg_Tools_main(graal_isolatethread_t *thread, dxfg_string_list* args);

/** @} */ // end of Qds

#ifdef __cplusplus
}
#endif

#endif // DXFG_QDS_H
