// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

/**
 * @file
 * @brief The v4 functions of dxfg_qds.h and the v4 names of its types (see the group @ref V4_Qds).
 */

#ifndef DXFG_V4_QDS_H
#define DXFG_V4_QDS_H

#include "dxfg_qds.h"
#include "dxfg_v4_common.h"
#include "dxfg_v4_javac.h"

#ifdef __cplusplus
extern "C" {
#endif

/** @defgroup V4_Qds Qds (v4)
 *  @brief The v4 functions and type names of the group @ref Qds. Every function returns
 *  #DXFG_EXECUTE_SUCCESSFULLY or #DXFG_EXECUTE_FAIL and writes its result to a #DXFG_OUT parameter,
 *  see @ref V4_Conventions.
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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] symbolList The specification of the symbol list: a UTF-8 null-terminated string; the caller keeps its
 * ownership.
 * @param[out] symbols The pointer to the variable that receives the list of the symbols (release it with
 * dxfg_v4_string_list_release()); must not be NULL.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) on successful function execution or #DXFG_EXECUTE_FAIL (-1) on error. Use
 * dxfg_v4_get_and_clear_thread_exception() to determine if an exception was thrown.
 */
int32_t dxfg_v4_Tools_parseSymbols(graal_isolatethread_t *thread, const char* symbolList,
                                   DXFG_OUT dxfg_v4_string_list **symbols);

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
 * @param[in] thread The current GraalVM Isolate's thread.
 * @param[in] args The list of the arguments; the caller keeps its ownership.
 * @return #DXFG_EXECUTE_SUCCESSFULLY (0) if no arguments are passed (the process is terminated otherwise, see above),
 * or #DXFG_EXECUTE_FAIL (-1) on error. Use dxfg_v4_get_and_clear_thread_exception() to determine if an exception was
 * thrown.
 */
int32_t dxfg_v4_Tools_main(graal_isolatethread_t *thread, dxfg_v4_string_list* args);

/** @} */ // end of V4_Qds

#ifdef __cplusplus
}
#endif

#endif // DXFG_V4_QDS_H
