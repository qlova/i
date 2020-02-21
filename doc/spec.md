# Island

## Reserved words:
if, for, in, go, return, break, main, to, i, usm

## Builtin constants
true := 1
false := 0

## Builtin functions
print, out, freeze, embed, open, type, subtype

## Standard library
os

## Prototypes
number, collection, instruction, data, connection, metatype, nothing, thing, magic

## Magic types
native, tag, unit, sequencer, undefined, dynamic

### Number types
real, rational, natural, integer, duplex, complex, quaternion, octonion, sedenion

### Data types
string, symbol, bit, byte, color, image, sound, video, time, logical

### Collection types
frozen, array, list, table, tensor, vector, matrix, set, tree, stack, queue, option, database, sequence

### Connection types
pipe, file, http, tcp, udp, 

### Instruction types
script, function, concept, builtin

## Ownership

$= means purchase & store inside.
= means alias.

$ means to trade/own.
% means to share. (reference count)
@ means to borrow.

copy by default.

## Symbols
` = literal name
! = not
 # = length/break
% = share
& = and
_ = the empty set/type
| = or/else
" = string
' = symbol
: = line-level block
} = close block
. = static_index/privacy
[ = dynamic_index
, = seperator/sequencer
{ = fake-block
; = tag statement / error handling
^ = aliaser

## Operators
-+/*\^><=~
