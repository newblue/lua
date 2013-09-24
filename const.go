// Package lua provides an interface to the Lua 5.1 interpreter.
package lua

/*
#include <lua.h>
#include <lualib.h>
*/
import "C"
import "errors"

const (
	Version    = C.LUA_VERSION
	Versionnum = C.LUA_VERSION_NUM
	Copyright  = C.LUA_COPYRIGHT
)

const (
	Signature = C.LUA_SIGNATURE // mark for precompiled code (`<esc>Lua')
	Multret   = C.LUA_MULTRET   // option for multiple returns in 'call' functions
	Minstack  = C.LUA_MINSTACK  // minimum Lua stack available to a Go function
)

// Thread status; 0 is OK
const (
	Ok        = 0
	Yield     = C.LUA_YIELD
	Errrun    = C.LUA_ERRRUN
	Errsyntax = C.LUA_ERRSYNTAX
	Errmem    = C.LUA_ERRMEM
	Errerr    = C.LUA_ERRERR
)

var errs map[int]error = map[int]error{
	Errrun:    errors.New("run time error"),
	Errsyntax: errors.New("syntax error"),
	Errmem:    errors.New("out of memory"),
	Errerr:    errors.New("error in error handling"),
}

func numtoerror(errnum int) error {
	if errnum < 1 {
		return nil
	}
	if e, ok := errs[errnum]; ok {
		return e
	}
	return errors.New("unknown error")
}

// Pseudo-indices. Unless otherwise noted, any function that accepts valid
// indices can also be called with these pseudo-indices, which represent
// some Lua values that are accessible to Go code but which are not in
// the stack. Pseudo-indices are used to access the thread environment,
// the function environment, the registry, and the upvalues of a Go function.
//
// The thread environment (where global variables live) is always at
// pseudo-index Globalsindex. The environment of the running Go function
// is always at pseudo-index Environindex.
//
// To access and change the value of global variables, you can use regular
// table operations over an environment table. For instance, to access the
// value of a global variable, do:
//	s.Getfield(luajit.Globalsindex, varname);
const (
	Registryindex = C.LUA_REGISTRYINDEX
	Environindex  = C.LUA_ENVIRONINDEX // env of running Go function
	Globalsindex  = C.LUA_GLOBALSINDEX // thread env, where globals live
)

// Returns the pseudo-index for the nth upvalue of a Go closure.
//
// Whenever a Go closure is called from Lua, its upvalues are located
// at specific pseudo-indices. These pseudo-indices are located using
// Upvalueindex. The first value associated with a function is at position
// Upvalueindex(1), and so on.
func Upvalueindex(n int) int {
	return (Globalsindex - n) + 1 // Upvalueindex(1) is reserved for Go func pointer
}

// Basic types
const (
	Tnone          = C.LUA_TNONE
	Tnil           = C.LUA_TNIL
	Tboolean       = C.LUA_TBOOLEAN
	Tlightuserdata = C.LUA_TLIGHTUSERDATA
	Tnumber        = C.LUA_TNUMBER
	Tstring        = C.LUA_TSTRING
	Ttable         = C.LUA_TTABLE
	Tfunction      = C.LUA_TFUNCTION
	Tuserdata      = C.LUA_TUSERDATA
	Tthread        = C.LUA_TTHREAD
)

// Garbage-collection function and options
const (
	// Stops the garbage collector.
	GCstop = C.LUA_GCSTOP
	// Restarts the garbage collector.
	GCrestart = C.LUA_GCRESTART
	// Performs a full garbage-collection cycle.
	GCcollect = C.LUA_GCCOLLECT
	// Returns the current amount of memory (in Kbytes) in use by Lua.
	GCcount = C.LUA_GCCOUNT
	// Returns the remainder of dividing the current amount of bytes of memory
	// in use by Lua by 1024.
	GCcountb = C.LUA_GCCOUNTB
	// Performs an incremental step of garbage collection. The step "size" is
	// controlled by data (larger values mean more steps) in a non-specified
	// way. If you want to control the step size you must experimentally
	// tune the value of data. The function returns 1 if the step finished a
	// garbage-collection cycle.
	GCstep = C.LUA_GCSTEP
	// Sets data as the new value for the pause of the collector. The function
	// returns the previous value of the pause.
	GCsetpause = C.LUA_GCSETPAUSE
	// Sets data as the new value for the step multiplier of the collector. The
	// function returns the previous value of the step multiplier.
	GCsetstepmul = C.LUA_GCSETSTEPMUL
)

// Debug event codes
const (
	// The call hook is called when the interpreter calls a function. The
	// hook is called just after Lua enters the new function, before
	// the function gets its arguments.
	Hookcall = C.LUA_HOOKCALL
	// The return hook is called when the interpreter returns from
	// a function. The hook is called just before Lua leaves the
	// function. You have no access to the values to be returned by
	// the function.
	Hookret = C.LUA_HOOKRET
	// The line hook is called when the interpreter is about to start
	// the execution of a new line of code, or when it jumps back in
	// the code (even to the same line). (This event only happens while
	// Lua is executing a Lua function.)
	Hookline = C.LUA_HOOKLINE
	// The count hook is called after the interpreter executes every
	// count instructions. (This event only happens while Lua is
	// executing a Lua function.)
	Hookcount   = C.LUA_HOOKCOUNT
	Hooktailret = C.LUA_HOOKTAILRET
)

// Debug event masks
const (
	Maskcall  = 1 << Hookcall
	Maskret   = 1 << Hookret
	Maskline  = 1 << Hookline
	Maskcount = 1 << Hookcount
)

// These are for Sethook and others
const (
	namehooks = "_hooks"
	namecall  = "call"
	nameret   = "ret"
	nameline  = "line"
	namecount = "count"
)

// lualib constants
const (
	Filehandle  = C.LUA_FILEHANDLE
	Colibname   = C.LUA_COLIBNAME   // coroutine
	Tablibname  = C.LUA_TABLIBNAME  // table
	IOlibname   = C.LUA_IOLIBNAME   // io
	OSlibname   = C.LUA_OSLIBNAME   // os
	Strlibname  = C.LUA_STRLIBNAME  // string
	Mathlibname = C.LUA_MATHLIBNAME // math
	Dblibname   = C.LUA_DBLIBNAME   // debug
	Loadlibname = C.LUA_LOADLIBNAME // package
)
